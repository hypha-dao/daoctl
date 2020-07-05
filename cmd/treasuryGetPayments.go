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

var treasuryGetPaymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "view a table of payments",
	//Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		printPaymentsTable(ctx, getAPI(), "HUSD Payments")
	},
}

func printPaymentsTable(ctx context.Context, api *eos.API, title string) {
	fmt.Println("\n", title)
	payments := models.Payments(ctx, api)
	paymentsTable := views.PaymentTable(payments)
	paymentsTable.SetStyle(simpletable.StyleCompactLite)
	fmt.Println("\n" + paymentsTable.String() + "\n\n")
}

func init() {
	treasuryGetCmd.AddCommand(treasuryGetPaymentsCmd)
}
