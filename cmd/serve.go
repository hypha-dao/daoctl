package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/hyperion"
	"github.com/hypha-dao/daoctl/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func assetToFloat(a *eos.Asset) float64 {
	f, _ := big.NewFloat(float64(a.Amount) / math.Pow10(int(a.Precision))).Float64()
	return f
}

func getVotingEvents(ctx context.Context) {

	const Format = "2006-01-02T15:04:05"
}

func getTokenSupply(ctx context.Context, api *eos.API, tokenContract, symbol string) (float64, error) {
	type Supply struct {
		TokenSupply eos.Asset `json:"supply"`
	}

	var supply []Supply
	// telosDecide := eos.MustStringToName(viper.GetString("TelosDecideContract"))

	var request eos.GetTableRowsRequest
	request.Code = tokenContract
	request.Scope = symbol
	request.Table = "stat"
	request.Limit = 1
	request.JSON = true
	response, err := api.GetTableRows(ctx, request)
	if err != nil {
		return 0, err
	}
	response.JSONToStructs(&supply)
	return assetToFloat(&supply[0].TokenSupply), nil
}

// testCmd represents the test command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Hypha prometheus metrics server",

	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		api := getAPI()

		hyphaSupply := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "hypha_supply",
			Help: "Total amount of HYPHA tokens in circulation",
		})

		hvoiceSupply := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "hvoice_supply",
			Help: "Total amount of HVOICE tokens rewarded",
		})

		husdSupply := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "husd_supply",
			Help: "Total amount of HUSD tokens circulating",
		})

		// memberCount := prometheus.NewGauge(prometheus.GaugeOpts{
		// 	Name: "member_count",
		// 	Help: "Total number of members",
		// })

		// applicantCount := prometheus.NewGauge(prometheus.GaugeOpts{
		// 	Name: "applicant_count",
		// 	Help: "Total number of open applications",
		// })

		voteEventCount := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "vote_events",
			Help: "Number of times users have voted in last 24 hours",
		})

		daoEventCount := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "dao_events",
			Help: "Number of actions called on the contract in last 24 hours",
		})

		// Periodically update the metrics
		go func() {
			for {
				sup, err := models.GetHvoiceSupply(ctx, api)
				if err == nil {
					// fmt.Println("Updated HVOICE supply to: ", sup.String())
					hvoiceSupply.Set(assetToFloat(sup))
				} else {
					fmt.Println("Retrieved an error retrieving Hypha supply: ", err)
				}

				hypha, err := getTokenSupply(ctx, api, viper.GetString("RewardToken.Contract"), viper.GetString("RewardToken.Symbol"))
				if err == nil {
					// fmt.Println("Updated HYPHA supply to: ", hypha)
					hyphaSupply.Set(hypha)
				} else {
					fmt.Println("Retrieved an error retrieving Hypha supply: ", err)
				}

				husd, err := getTokenSupply(ctx, api, viper.GetString("Treasury.TokenContract"), viper.GetString("Treasury.Symbol"))
				if err == nil {
					// fmt.Println("Updated HUSD supply to: ", husd)
					husdSupply.Set(husd)
				} else {
					fmt.Println("Retrieved an error retrieving HUSD supply: ", err)
				}

				// memberCount.Set(float64(len(models.Members(ctx, api))))
				// applicantCount.Set(float64(len(models.Applicants(ctx, api))))

				query := hyperion.NewQuery("castvote", viper.GetString("TelosDecideContract"), "")
				query.After = time.Now().AddDate(0, 0, -1)
				query.Limit = 1000
				results, err := query.Results()
				if err == nil {
					// fmt.Println("Updated vote count ", len(results))
					voteEventCount.Set(float64(len(results)))
				} else {
					fmt.Println("Error updating the vote count: ", err)
				}

				query = hyperion.NewQuery("", viper.GetString("DAOContract"), "")
				query.After = time.Now().AddDate(0, 0, -1)
				query.Limit = 1000
				results, err = query.Results()
				if err == nil {
					// fmt.Println("Updated dao event count ", len(results))
					daoEventCount.Set(float64(len(results)))
				} else {
					fmt.Println("Error updating the dao event count: ", err)
				}

				time.Sleep(time.Duration(time.Minute))
			}
		}()

		bind := ""
		flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		flagset.StringVar(&bind, "bind", ":"+viper.GetString("ServePort"), "The socket to bind to.")
		fmt.Println("Binding daoctl serve port to: ", viper.GetString("ServePort"))
		flagset.Parse(os.Args[1:])

		r := prometheus.NewRegistry()
		r.MustRegister(hyphaSupply)
		r.MustRegister(hvoiceSupply)
		r.MustRegister(husdSupply)
		// r.MustRegister(memberCount)
		// r.MustRegister(applicantCount)
		r.MustRegister(voteEventCount)
		r.MustRegister(daoEventCount)
		http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
		log.Fatal(http.ListenAndServe(bind, nil))
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
