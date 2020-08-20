package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"

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
		ctx := context.Background()

		if viper.GetBool("global-csv") {
			requests := models.Requests(ctx, getAPI(), viper.GetBool("treasury-get-requests-cmd-all"))
			requestsTable := views.RequestTable(requests)
			csvData := models.TableToData(requestsTable)

			file, err := os.Create(viper.GetString("global-output-file"))
			if err != nil {
				log.Fatalln("error writing csv:", err)
			}

			defer file.Close()

			w := csv.NewWriter(file)
			w.WriteAll(csvData) // calls Flush internally

			if err := w.Error(); err != nil {
				log.Fatalln("error writing csv:", err)
			}
		} else {
			printRequestsTable(ctx, getAPI(), "HUSD Redemption Requests")
		}
	},
}

func printRequestsTable(ctx context.Context, api *eos.API, title string) {
	fmt.Println("\n", title)
	requests := models.Requests(ctx, api, viper.GetBool("treasury-get-requests-cmd-all"))
	requestsTable := views.RequestTable(requests)
	requestsTable.SetStyle(simpletable.StyleCompactLite)
	fmt.Println("\n" + requestsTable.String() + "\n\n")
}

func init() {
	treasuryGetCmd.AddCommand(treasuryGetRequestsCmd)
	treasuryGetRequestsCmd.Flags().BoolP("all", "", false, "include all requests or only requests with additional amounts due")
}
