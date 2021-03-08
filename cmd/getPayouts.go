package cmd

import (
	"context"
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/daoctl/views"
	"github.com/spf13/cobra"
)

var getPayoutCmd = &cobra.Command{
	Use:   "payouts [account name]",
	Short: "retrieve payouts",
	Long:  "retrieve all payouts",
	// Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		// api := eos.New(viper.GetString("EosioEndpoint"))
		// ctx := context.Background()

		// periods := models.LoadPeriods(api, true, true)

		// if viper.GetBool("global-active") == true {
		// 	printPayoutTable(ctx, api, periods, "Completed Payouts", "payout")
		// }

		// if viper.GetBool("global-include-proposals") == true {
		// 	printPayoutTable(ctx, api, periods, "Open Payout Proposals", "proposal")
		// }

		// if viper.GetBool("global-failed-proposals") == true {
		// 	printPayoutTable(ctx, api, periods, "Failed Payout Proposals", "failedprops")
		// }

		// if viper.GetBool("global-include-archive") == true {
		// 	printPayoutTable(ctx, api, periods, "Archived Payout Proposals", "proparchive")
		// }
	},
}

func printPayoutTable(ctx context.Context, api *eos.API, periods []models.Period, title, scope string) {
	fmt.Println("\n", title)
	payouts := models.Payouts(ctx, api, periods, scope)
	payoutsTable := views.PayoutTable(payouts)
	payoutsTable.SetStyle(simpletable.StyleCompactLite)
	fmt.Println("\n" + payoutsTable.String() + "\n\n")
}

func init() {
	getCmd.AddCommand(getPayoutCmd)
	//getPayoutCmd.Flags().BoolP("include-proposals", "i", false, "include proposals in the output")
	// getPayoutCmd.Flags().BoolP("", "i", false, "include proposals in the output")
}
