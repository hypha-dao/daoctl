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

// Treasury ...
type Treasury struct {
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
	treasury.loadHolders(api, tokenContract, symbol)
	treasury.loadEthUSDT(viper.GetString("Treasury.EthUSDTContract"), viper.GetString("Treasury.EthUSDTAddress"))
	treasury.loadBtcBalance("")
	return treasury
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
	t.EthUSDTBalance, err = eos.NewAssetFromString(string(body) + " HUSD")
	if err != nil {
		fmt.Println("Unable to format ETH USDT balance as an asset type: " + string(body) + " HUSD")
	}
}

func (t *Treasury) loadBtcBalance(btcAddress string) {
	fmt.Println("Note: Bitcoin Treasury balance not yet supported.")
}

// func GetHusdBankBalance(api *eos.API) {

// 	// tokenContract := eos.AN(viper.GetString("TreasuryTokenContract")) //toAccount(viper.GetString("TreasuryTokenContract"), "TreasuryTokenContract account")
// 	tokenHoldings := GetTokenHoldings(api, viper.GetString("TreasuryTokenContract"), viper.GetString("TreasurySymbol"))

// 	// treasuryAccount := toAccount(viper.GetString("TreasuryContract"), "treasury contract")
// 	// balances, err := api.GetCurrencyBalance(context.Background(), treasuryAccount, treasuryTokenContract)
// }
