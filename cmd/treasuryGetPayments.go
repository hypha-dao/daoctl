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

var treasuryGetPaymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "view a table of payments",
	//Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		if viper.GetBool("global-csv") {
			payments := models.Payments(ctx, getAPI())
			paymentsTable := views.PaymentTable(payments)
			csvData := models.TableToData(paymentsTable)

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
			printPaymentsTable(ctx, getAPI(), "HUSD Payments")
		}
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
