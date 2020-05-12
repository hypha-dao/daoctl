package cmd

import (
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/daoctl/views"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getTreasuryCmd = &cobra.Command{
	Use:   "treasury",
	Short: "Retrieve multi-chain balance information for the treasury",
	// Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Get Treasury command")
		api := eos.New(viper.GetString("EosioEndpoint"))

		treasuries := models.GetTokenHoldings(api, viper.GetString("TreasuryTokenContract"), viper.GetString("TreasurySymbol"))
		treasuryTable := views.TreasuryTable(treasuries)
		fmt.Println("\n\n" + treasuryTable.String() + "\n\n")
	},
}

func init() {
	getCmd.AddCommand(getTreasuryCmd)

	getTreasuryCmd.Flags().StringP("contract", "", "eosio.token", "Account managing the token")
	getTreasuryCmd.Flags().StringP("symbol", "", "", "Only query this symbol. Try EOS")
}
