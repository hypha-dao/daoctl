package cmd

import (
	"fmt"
	"strconv"

	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/daoctl/views"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getTreasuryCmd = &cobra.Command{
	Use:   "treasury",
	Short: "retrieve multi-chain balance information for the treasury",
	// Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))

		addlBalance, err := eos.NewAssetFromString(viper.GetString("get-treasury-cmd-addl-balance"))
		if err != nil {
			fmt.Println("Unable to read addl-balance parameter, using 0.00 HUSD")
			addlBalance, _ = eos.NewAssetFromString("0.00 HUSD")
		}

		treasury := models.LoadTreasury(api, viper.GetString("Treasury.TokenContract"), viper.GetString("Treasury.Symbol"))
		treasuryTable, circulatingBalance := views.TreasuryTable(treasury.TreasuryHolders)
		fmt.Println("\n" + treasuryTable.String() + "\n\n")

		totalAssets := treasury.EthUSDTBalance.Add(addlBalance)
		usd := float64(totalAssets.Amount)
		circulating := float64(circulatingBalance.Amount)
		coverage := float64(usd / circulating * 100)
		netTreasury := totalAssets.Sub(circulatingBalance)
		output := []string{
			fmt.Sprintf("Awaiting burning|%v", views.FormatAsset(&treasury.BankBalance)),
			fmt.Sprintf("Circulating|%v", views.FormatAsset(&circulatingBalance)),
			fmt.Sprintf("Eth USDT balance|%v", views.FormatAsset(&treasury.EthUSDTBalance)),
			fmt.Sprintf("Additional balances|%v", views.FormatAsset(&addlBalance)),
			fmt.Sprintf("Total balances|%v", views.FormatAsset(&totalAssets)),
			fmt.Sprintf("Net Treasury balance|%v", views.FormatAsset(&netTreasury)),
			string(fmt.Sprintf("Coverage|%v", strconv.FormatFloat(coverage, 'f', 3, 64)) + " %"),
		}

		fmt.Println(columnize.SimpleFormat(output))
	},
}

func init() {
	getCmd.AddCommand(getTreasuryCmd)
	getTreasuryCmd.Flags().StringP("addl-balance", "", "0.00 HUSD", "Value of other accounts (e.g. banking, BTC) to optionally add manually")
}
