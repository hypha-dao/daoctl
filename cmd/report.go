package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/eoscanada/eos-go"
	"github.com/go-echarts/go-echarts/charts"
	"github.com/hypha-dao/daoctl/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var reportCmd = &cobra.Command{
	Use:   "report [filter]",
	Short: "view reports",
	Run: func(cmd *cobra.Command, args []string) {
		// api := getAPI()
		// ctx := context.Background()
		// accountName := viper.GetString("query-cmd-account")
		// contractName := viper.GetString("query-cmd-contract")
		// actionName := viper.GetString("query-cmd-action")
		// trxID := viper.GetString("query-cmd-trx")
		// limit := viper.GetInt("query-cmd-limit")

		// zlog.Debug("Query parameters",
		// 	zap.String("account", accountName),
		// 	zap.String("contract", contractName),
		// 	zap.String("action", actionName),
		// 	zap.Int("limit", limit))

		// request := viper.GetString("HyperionEndpoint") + "/history/"

		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		periods := models.LoadPeriods(api, true, true)
		// includeExpired := viper.GetBool("get-assignments-cmd-expired")
		roles := models.Roles(ctx, api, periods, "role")

		assignments, err := models.Assignments(ctx, api, roles, periods, "proposal", false)
		if err != nil {
			fmt.Println("Cannot get list of assignments: " + err.Error())
			os.Exit(-1)
		}

		currentPeriod, _ := models.CurrentPeriod(&periods)

		var maxPeriods int
		maxPeriods = 26
		numPeriods := Min(maxPeriods, len(periods)-int(currentPeriod))

		var periodEndingDates []string
		var husdPerPeriod, hvoicePerPeriod, hyphaPerPeriod []float64
		periodEndingDates = make([]string, numPeriods)
		husdPerPeriod = make([]float64, numPeriods)
		hvoicePerPeriod = make([]float64, numPeriods)
		hyphaPerPeriod = make([]float64, numPeriods)

		for i := 0; i < numPeriods; i++ {
			for _, assignment := range assignments {
				if periods[i+int(currentPeriod)].PeriodID >= assignment.StartPeriod.PeriodID && periods[i+int(currentPeriod)].PeriodID <= assignment.EndPeriod.PeriodID {
					husdPerPeriod[i] += float64(assignment.HusdPerPhase.Amount) / 100
					hvoicePerPeriod[i] += float64(assignment.HvoicePerPhase.Amount) / 100
					hyphaPerPeriod[i] += float64(assignment.HyphaPerPhase.Amount) / 100
				}
			}
			periodEndingDates[i] = periods[i+int(currentPeriod)].EndTime.Time.Format("2006 Jan 02")
		}

		line := charts.NewLine()
		line.SetGlobalOptions(charts.TitleOpts{Title: "Assignment Expenses by Period"})
		line.AddXAxis(periodEndingDates).
			AddYAxis("husd", husdPerPeriod)
			// AddYAxis("hvoice", hvoicePerPeriod).
			// AddYAxis("hypha", hyphaPerPeriod)

		f, err := os.Create("line.html")
		if err != nil {
			log.Println(err)
		}
		line.Render(f)

	},
}

// For 2 values
func Min(value_0, value_1 int) int {
	if value_0 < value_1 {
		return value_0
	}
	return value_1
}

// For 1+ values
func Mins(value int, values ...int) int {
	for _, v := range values {
		if v < value {
			value = v
		}
	}
	return value
}

func init() {
	RootCmd.AddCommand(reportCmd)
	// reportCmd.Flags().StringP("account", "", "", "member's account to query")
	// reportCmd.Flags().StringP("contract", "", "dao.hypha", "query actions called on this smart contract (defaults to DAO)")
	// reportCmd.Flags().StringP("action", "", "", "action name to query on the contract (defaults to create)")
	// reportCmd.Flags().StringP("trx", "", "", "transaction ID to query for the full content of that transaction")
	// reportCmd.Flags().IntP("limit", "", 10, "maximum number of records to retrieve")
	// reportCmd.Flags().BoolP("json", "j", false, "print the results as JSON")

}
