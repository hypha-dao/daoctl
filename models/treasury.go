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

// TreasuryHolder ...
type TreasuryHolder struct {
	TokenHolder eos.Name
	Balance     eos.Asset
}

// TreasuryConfig struct
type TreasuryConfig struct {
	RedemptionSymbol        *string
	DAOContract             *eos.Name
	RedemptionTokenContract *eos.Name
	Paused                  *uint64
	Threshold               *uint64
	RawTreasuryConfig
}

type RawTreasuryConfig struct {
	RedemptionSymbol string             `json:"redemption_symbol"`
	Names            []NameKV           `json:"names"`
	Strings          []StringKV         `json:"strings"`
	Assets           []AssetKV          `json:"assets"`
	Ints             []IntKV            `json:"ints"`
	UpdatedDate      eos.BlockTimestamp `json:"updated_date"`
}

// Treasury ...
type Treasury struct {
	Config          TreasuryConfig
	TreasuryHolders []TreasuryHolder
	BankBalance     eos.Asset
	EthUSDTBalance  eos.Asset
	BtcBalance      eos.Asset
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

// LoadTreasury ...
func LoadTreasury(api *eos.API, tokenContract, symbol string) Treasury {
	var treasury Treasury
	treasury.loadConfig(api)
	treasury.loadHolders(api, tokenContract, symbol)
	treasury.loadEthUSDT(viper.GetString("Treasury.EthUSDTContract"), viper.GetString("Treasury.EthUSDTAddress"))
	treasury.loadBtcBalance("")
	return treasury
}

func (t *Treasury) loadConfig(api *eos.API) {

	// LoadTreasConfig loads the treasury configuration from the smart contract
	var rto []RawTreasuryConfig
	var request eos.GetTableRowsRequest
	request.Code = viper.GetString("Treasury.Contract")
	request.Scope = viper.GetString("Treasury.Contract")
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

	t.Config.RawTreasuryConfig = rto[0] // keep the raw object around
	t.Config.RedemptionSymbol = &rto[0].RedemptionSymbol
}

func (t *Treasury) loadHolders(api *eos.API, tokenContract, symbol string) {
	var request eos.GetTableByScopeRequest
	request.Code = tokenContract
	request.Table = "accounts"
	request.Limit = 500 // TODO: move to a MaxMembers parameter or check the "more" return value
	response, err := api.GetTableByScope(context.Background(), request)
	errorCheck("get table by scope", err)

	var scopes []Scope
	json.Unmarshal(response.Rows, &scopes)

	var treasuryHolder []TreasuryHolder
	treasuryHolder = make([]TreasuryHolder, len(scopes))

	for index, scope := range scopes {

		tokenHolder := eos.AccountName(scope.Scope)
		balances, err := api.GetCurrencyBalance(context.Background(), tokenHolder, symbol, eos.AN(tokenContract))
		errorCheck("treasury", err)
		if len(balances) > 0 {
			if string(scope.Scope) == viper.GetString("Treasury.Contract") {
				t.BankBalance = balances[0]
			} else {
				treasuryHolder[index].TokenHolder = scope.Scope
				treasuryHolder[index].Balance = balances[0]
			}
		}
	}
	t.TreasuryHolders = treasuryHolder
}

func (t *Treasury) loadEthUSDT(usdtTokenAddress, treasuryWallet string) {
	requestString := "https://api.tokenbalance.com/balance/" + viper.GetString("Treasury.EthUSDTContract") + "/" + viper.GetString("Treasury.EthUSDTAddress")
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
		fmt.Println("Unable to format ETH USDT balance as an asset type: " + string(body) + " HUSD")
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
