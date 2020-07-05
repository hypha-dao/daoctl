package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eosc/cli"
	"github.com/spf13/viper"
)

// Balance ...
type Balance struct {
	Balance              eos.Asset
	RequestedRedemptions eos.Asset
}

// Config struct
type Config struct {
	RedemptionSymbol        *string
	DAOContract             *eos.Name
	RedemptionTokenContract *eos.Name
	Paused                  *uint64
	Threshold               *uint64
	RawConfig               rawConfig
}

type rawConfig struct {
	RedemptionSymbol string             `json:"redemption_symbol"`
	Names            []NameKV           `json:"names"`
	Strings          []StringKV         `json:"strings"`
	Assets           []AssetKV          `json:"assets"`
	Ints             []IntKV            `json:"ints"`
	UpdatedDate      eos.BlockTimestamp `json:"updated_date"`
}

// Treasury ...
type Treasury struct {
	Config              Config
	Members             map[eos.Name]Balance
	RedemptionRequests  []RedemptionRequest
	TotalReqRedemptions eos.Asset
	Circulating         eos.Asset
	BankBalance         eos.Asset
	EthUSDTBalance      eos.Asset
	BtcBalance          eos.Asset
}

func errorCheck(prefix string, err error) {
	if err != nil {
		fmt.Printf("ERROR: %s: %s\n", prefix, err)
		os.Exit(1)
	}
}

func toAccount(in, field string) eos.AccountName {
	acct, err := cli.ToAccountName(in)
	errorCheck(fmt.Sprintf("invalid account format for %q", field), err)

	return acct
}

// Load ...
func Load(api *eos.API, treasuryContract, tokenContract, symbol string) Treasury {
	var treasury Treasury
	treasury.loadConfig(api, treasuryContract)
	// treasury.Members = make(map[eos.Name]TreasuryBalance)
	treasury.loadMembers(api, treasuryContract, tokenContract, symbol)
	treasury.loadEthUSDT(viper.GetString("Treasury.EthUSDTContract"), viper.GetString("Treasury.EthUSDTAddress"))
	treasury.loadBtcBalance("")
	return treasury
}

func (t *Treasury) loadConfig(api *eos.API, treasuryContract string) {

	// LoadTreasConfig loads the treasury configuration from the smart contract
	var rto []rawConfig
	var request eos.GetTableRowsRequest
	request.Code = treasuryContract
	request.Scope = treasuryContract
	request.Table = "config"
	request.Limit = 1
	request.JSON = true
	response, _ := api.GetTableRows(context.Background(), request)
	response.JSONToStructs(&rto)

	// bookmark known values
	for index := range rto[0].Names {
		if rto[0].Names[index].Key == "dao_contract" {
			t.Config.DAOContract = &rto[0].Names[index].Value
		} else if rto[0].Names[index].Key == "token_redemption_contract" {
			t.Config.RedemptionTokenContract = &rto[0].Names[index].Value
		}
	}

	for index := range rto[0].Ints {
		if rto[0].Ints[index].Key == "paused" {
			t.Config.Paused = &rto[0].Ints[index].Value
		} else if rto[0].Ints[index].Key == "threshold" {
			t.Config.Threshold = &rto[0].Ints[index].Value
		}
	}

	t.Config.RawConfig = rto[0] // keep the raw object around
	t.Config.RedemptionSymbol = &rto[0].RedemptionSymbol
}

func (t *Treasury) getRedemptionRequests(api *eos.API, treasuryContract string) map[eos.Name]eos.Asset {
	var request eos.GetTableRowsRequest
	request.Code = treasuryContract
	request.Scope = treasuryContract
	request.Table = "redemptions"
	request.Limit = 500 // TODO: max redemptions?
	request.JSON = true

	var rr []RedemptionRequest
	rrResponse, err := api.GetTableRows(context.Background(), request)
	if err != nil {
		panic(err)
	}
	rrResponse.JSONToStructs(&rr)
	redemptionMap := make(map[eos.Name]eos.Asset)

	t.TotalReqRedemptions, _ = eos.NewAssetFromString("0.00 HUSD")
	for _, element := range rr {
		_, exists := redemptionMap[element.Requestor]
		if exists {
			redemptionMap[element.Requestor] = redemptionMap[element.Requestor].Add(element.Requested)
		} else {
			redemptionMap[element.Requestor] = element.Requested
		}

		t.TotalReqRedemptions = t.TotalReqRedemptions.Add(element.Requested)
	}
	return redemptionMap
}

func (t *Treasury) getHolders(api *eos.API, treasuryContract, tokenContract, symbol string) map[eos.Name]eos.Asset {
	var request eos.GetTableByScopeRequest
	request.Code = tokenContract
	request.Table = "accounts"
	request.Limit = 500 // TODO: move to a MaxMembers parameter or check the "more" return value
	response, err := api.GetTableByScope(context.Background(), request)
	errorCheck("get table by scope", err)

	var scopes []Scope
	json.Unmarshal(response.Rows, &scopes)

	holders := make(map[eos.Name]eos.Asset)
	t.Circulating, _ = eos.NewAssetFromString("0.00 HUSD")

	for _, scope := range scopes {

		tokenHolder := eos.AccountName(scope.Scope)
		balances, err := api.GetCurrencyBalance(context.Background(), tokenHolder, symbol, eos.AN(tokenContract))
		errorCheck("getting currency balance", err)

		if len(balances) > 0 {
			if string(scope.Scope) == treasuryContract {
				t.BankBalance = balances[0]
			} else {
				holders[scope.Scope] = balances[0]
			}
		}
		t.Circulating = t.Circulating.Add(balances[0])
	}
	return holders
}

func (t *Treasury) loadMembers(api *eos.API, treasuryContract, tokenContract, symbol string) {
	holderBalances := t.getHolders(api, treasuryContract, tokenContract, symbol)
	rrMap := t.getRedemptionRequests(api, treasuryContract)
	zeroHusd, _ := eos.NewAssetFromString("0.00 HUSD")

	t.Members = make(map[eos.Name]Balance)
	for holder, tokenBalance := range holderBalances {
		t.Members[holder] = Balance{Balance: tokenBalance, RequestedRedemptions: zeroHusd}
	}

	for Requestor, requestedAmount := range rrMap {
		_, exists := t.Members[Requestor]
		if !exists {
			t.Members[Requestor] = Balance{Balance: zeroHusd, RequestedRedemptions: requestedAmount}
		} else {
			t.Members[Requestor] = Balance{Balance: t.Members[Requestor].Balance, RequestedRedemptions: requestedAmount}
		}
	}
}

func (t *Treasury) loadEthUSDT(usdtTokenAddress, treasuryWallet string) {
	requestString := "https://api.tokenbalance.com/balance/" + viper.GetString("Treasury.EthUSDTContract") + "/" + viper.GetString("Treasury.EthUSDTAddress")
	fmt.Println(requestString)
	resp, err := http.Get(requestString)
	if err != nil {
		fmt.Println("Unable to retrieve ETH USDT balance.")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unable to retrieve ETH USDT balance.")
		return
	}
	// fmt.Println("USDT balance: ", string(body))
	if string(body) == "0.0" {
		fmt.Println("WARNING: Balance of USDT in multisig wallet is 0.00")
		t.EthUSDTBalance, err = eos.NewAssetFromString("0.00 HUSD")
		return
	}
	t.EthUSDTBalance, err = eos.NewAssetFromString(string(body) + " HUSD")
	if err != nil {
		fmt.Println("Unable to format ETH USDT balance as an asset type: " + string(body) + " HUSD.  Assuming 0.00 HUSD")
		t.EthUSDTBalance, _ = eos.NewAssetFromString("0.00 HUSD")
		return
	}
}

func (t *Treasury) loadBtcBalance(btcAddress string) {
	fmt.Println("Note: Bitcoin Treasury balance not yet supported. Use the --addl-balance parameter to add the BTC balance.")
}

// func GetHusdBankBalance(api *eos.API) {

// 	// tokenContract := eos.AN(viper.GetString("TreasuryTokenContract")) //toAccount(viper.GetString("TreasuryTokenContract"), "TreasuryTokenContract account")
// 	tokenHoldings := GetTokenHoldings(api, viper.GetString("TreasuryTokenContract"), viper.GetString("TreasurySymbol"))

// 	// treasuryAccount := toAccount(viper.GetString("TreasuryContract"), "treasury contract")
// 	// balances, err := api.GetCurrencyBalance(context.Background(), treasuryAccount, treasuryTokenContract)
// }
