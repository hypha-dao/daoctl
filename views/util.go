package views

import (
	"math"
	"math/big"

	"github.com/eoscanada/eos-go"
	"github.com/leekchan/accounting"
	"github.com/spf13/viper"
)

// FormatAsset returns a string for an eos.Asset, taking into account the AssetsAsFloat configuration parameter
func FormatAsset(a *eos.Asset) string {
	ac := accounting.NewAccounting("", 0, ",", ".", "%s %v", "%s (%v)", "%s --") // TODO: make this configurable
	if viper.GetBool("AssetsAsFloat") {
		return ac.FormatMoneyBigFloat(big.NewFloat(float64(a.Amount) / math.Pow10(int(a.Precision))))
	}
	return a.String()
}
