package models

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eosc/cli"
)

type Treasury struct {
	TokenHolder eos.Name
	Balance     eos.Asset
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

func GetTokenHoldings(api *eos.API, tokenContract, symbol string) []Treasury {

	// func GetTokenHoldings(api *eos.API, tokenContract, symbol string) map[*eos.AccountName]*eos.Asset {
	var request eos.GetTableByScopeRequest
	request.Code = tokenContract
	request.Table = "accounts"
	request.Limit = 500 // TODO: move to a MaxMembers parameter or check the "more" return value
	response, err := api.GetTableByScope(context.Background(), request)
	errorCheck("get table by scope", err)

	var scopes []Scope
	json.Unmarshal(response.Rows, &scopes)

	var treasuries []Treasury
	treasuries = make([]Treasury, len(scopes))

	// var holdingsMap map[*eos.AccountName]*eos.Asset
	// holdingsMap = make(map[*eos.AccountName]*eos.Asset, len(scopes))
	for index, scope := range scopes {

		tokenHolder := eos.AccountName(scope.Scope)
		balances, err := api.GetCurrencyBalance(context.Background(), tokenHolder, symbol, eos.AN(tokenContract))
		errorCheck("treasury", err)
		if len(balances) > 0 {
      treasuries[index].TokenHolder = scope.Scope
      treasuries[index].Balance = balances[0]
      fmt.Println("Token holder: ", scope.Scope, " -- Balance: ", balances[0])
    }
		// holdingsMap[&tokenHolder] = &balances[0]
	}
	return treasuries
}

// func GetHusdBankBalance(api *eos.API) {

// 	// tokenContract := eos.AN(viper.GetString("TreasuryTokenContract")) //toAccount(viper.GetString("TreasuryTokenContract"), "TreasuryTokenContract account")
// 	tokenHoldings := GetTokenHoldings(api, viper.GetString("TreasuryTokenContract"), viper.GetString("TreasurySymbol"))

// 	// treasuryAccount := toAccount(viper.GetString("TreasuryContract"), "treasury contract")
// 	// balances, err := api.GetCurrencyBalance(context.Background(), treasuryAccount, treasuryTokenContract)
// }
