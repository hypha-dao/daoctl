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

var treasuryGetRequestsCmd = &cobra.Command{
	Use:   "requests",
	Short: "retrieve list of redemption requests",
	// Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		printRequestsTable(ctx, api, "HUSD Redemption Requests")
	},
}

func printRequestsTable(ctx context.Context, api *eos.API, title string) {
	fmt.Println("\n", title)
	requests := models.Requests(ctx, api)
	requestsTable := views.RequestTable(requests)
	requestsTable.SetStyle(simpletable.StyleCompactLite)
	fmt.Println("\n" + requestsTable.String() + "\n\n")
}

func init() {
	treasuryGetCmd.AddCommand(treasuryGetRequestsCmd)
}
