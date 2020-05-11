package cmd

import (
	"context"
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/daoctl/views"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getPayoutCmd = &cobra.Command{
	Use:   "payouts [account name]",
	Short: "retrieve payouts",
	Long:  "retrieve all payouts",
	// Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		periods := models.LoadPeriods(api)
		payouts := models.Payouts(ctx, api, periods)
		payoutsTable := views.PayoutTable(payouts)
		payoutsTable.SetStyle(simpletable.StyleCompactLite)

		fmt.Println("\n\n" + payoutsTable.String() + "\n\n")

		if viper.GetBool("get-payouts-cmd-include-proposals") == true {
			propPayout := models.ProposedPayouts(ctx, api, periods)
			propPayoutTable := views.PayoutTable(propPayout)
			propPayoutTable.SetStyle(simpletable.StyleCompactLite)
			fmt.Println("\n\n" + propPayoutTable.String() + "\n\n")
			return
		}
	},
}

func init() {
	getCmd.AddCommand(getPayoutCmd)
	getPayoutCmd.Flags().BoolP("include-proposals", "i", false, "include proposals in the output")
	// getPayoutCmd.Flags().BoolP("", "i", false, "include proposals in the output")
}
