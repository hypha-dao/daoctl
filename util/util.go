package util

import (
	"fmt"
	"math"
	"math/big"
	"os"

	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eosc/cli"
	"github.com/leekchan/accounting"
	"github.com/spf13/viper"
)

// ErrorCheck - need to improve this
// TODO: fix
func ErrorCheck(prefix string, err error) {
	if err != nil {
		fmt.Printf("ERROR: %s: %s\n", prefix, err)
		os.Exit(1)
	}
}

// ToAccount ...
func ToAccount(in, field string) eos.AccountName {
	acct, err := cli.ToAccountName(in)
	ErrorCheck(fmt.Sprintf("invalid account format for %q", field), err)

	return acct
}

// FormatAsset returns a string for an eos.Asset, taking into account the AssetsAsFloat configuration parameter
func FormatAsset(a *eos.Asset, precision int) string {
	ac := accounting.NewAccounting("", precision, ",", ".", "%s %v", "%s (%v)", "%s --") // TODO: make this configurable
	if viper.GetBool("AssetsAsFloat") {
		return ac.FormatMoneyBigFloat(big.NewFloat(float64(a.Amount) / math.Pow10(int(a.Precision))))
	}
	return a.String()
}

// AssetMult multiples the asset by the coeffecient
func AssetMult(a eos.Asset, coeff *big.Float) eos.Asset {
	//var amount big.Int
	var f big.Float
	f.SetInt(big.NewInt(int64(a.Amount)))
	amount, _ := f.Mul(&f, coeff).Int64() // big.NewInt(int64(a.Amount)).Mul(coeff)
	newAmount := eos.Int64(amount)
	// intObject := big.NewInt(amount)
	return eos.Asset{Amount: newAmount, Symbol: a.Symbol}
}

// Snip returns a shortened version of the notes map
func Snip(notes *map[string]string) string {

	charsToShow := 45
	suffix := "... <snip>"
	longValue := ""
	for k, v := range *notes {
		longValue = longValue + "{key:" + k + ",value:" + v + "}"
	}

	if len(longValue) < charsToShow {
		charsToShow = len(longValue)
		suffix = ""
	}
	return longValue[:charsToShow] + suffix
}
