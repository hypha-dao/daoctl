package cmd

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/dao-contracts/dao-go"
	"github.com/hypha-dao/daoctl/hyperion"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/document-graph/docgraph"
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

func getTokenBalance(ctx context.Context, api *eos.API, tokenContract, accountname, symbol string) (eos.Asset, error) {
	type Balance struct {
		Balance eos.Asset `json:"balance"`
	}

	var balances []Balance
	var request eos.GetTableRowsRequest
	request.Code = tokenContract
	request.Scope = accountname
	request.Table = "accounts"
	request.Limit = 100
	request.JSON = true
	response, err := api.GetTableRows(ctx, request)
	if err != nil {
		return eos.Asset{}, fmt.Errorf("could not get balance: GetTableRows: token contract: "+tokenContract+" account: "+accountname+" symbol: "+symbol+": %v", err)
	}

	err = response.JSONToStructs(&balances)
	if err != nil {
		return eos.Asset{}, fmt.Errorf("could not get balance: JSONToStructs: token contract: "+tokenContract+" account: "+accountname+" symbol: "+symbol+": %v", err)
	}

	if len(balances) == 0 {
		return eos.Asset{}, fmt.Errorf("could not get balance: query returned zero rows: token contract: " + tokenContract + " account: " + accountname + " symbol: " + symbol)
	}

	for _, balance := range balances {
		if balance.Balance.Symbol.Symbol == symbol {
			return balance.Balance, nil
		}
	}

	return eos.Asset{}, fmt.Errorf("could not get balance: no rows match the symbol provided: token contract: " + tokenContract + " account: " + accountname + " symbol: " + symbol)
}

func getTokenSupply(ctx context.Context, api *eos.API, tokenContract, symbol string) (eos.Asset, error) {
	type Supply struct {
		TokenSupply eos.Asset `json:"supply"`
	}

	var supply []Supply
	var request eos.GetTableRowsRequest
	request.Code = tokenContract
	request.Scope = symbol
	request.Table = "stat"
	request.Limit = 1
	request.JSON = true
	response, err := api.GetTableRows(ctx, request)

	if err != nil {
		return eos.Asset{}, fmt.Errorf("could not get supply: GetTableRows: token contract: "+tokenContract+"  symbol: "+symbol+": %v", err)
	}

	err = response.JSONToStructs(&supply)
	if err != nil {
		return eos.Asset{}, fmt.Errorf("could not get supply: JSONToStructs: token contract: "+tokenContract+"  symbol: "+symbol+": %v", err)
	}

	if len(supply) == 0 {
		return eos.Asset{}, fmt.Errorf("could not get supply: query returned zero rows: token contract: " + tokenContract + "  symbol: " + symbol)
	}
	return supply[0].TokenSupply, nil
}

func getLegacyObjectsRange(ctx context.Context, api *eos.API, contract eos.AccountName, scope eos.Name, id, count int) ([]dao.Object, bool, error) {

	var objects []dao.Object
	var request eos.GetTableRowsRequest
	request.LowerBound = strconv.Itoa(id)
	request.Code = string(contract)
	request.Scope = string(scope)
	request.Table = "objects"
	request.Limit = uint32(count)
	request.JSON = true
	response, err := api.GetTableRows(ctx, request)
	if err != nil {
		return []dao.Object{}, false, fmt.Errorf("get table rows %v", err)
	}

	err = response.JSONToStructs(&objects)
	if err != nil {
		return []dao.Object{}, false, fmt.Errorf("json to structs %v", err)
	}
	return objects, response.More, nil
}

func getLegacyObjects(ctx context.Context, api *eos.API, contract eos.AccountName, scope eos.Name) ([]dao.Object, error) {

	var allObjects []dao.Object

	cursor := 0
	batchSize := 45

	batch, more, err := getLegacyObjectsRange(ctx, api, contract, scope, cursor, batchSize)
	if err != nil {
		return []dao.Object{}, fmt.Errorf("json to structs %v", err)
	}
	allObjects = append(allObjects, batch...)

	for more {
		cursor += batchSize
		batch, more, err = getLegacyObjectsRange(ctx, api, contract, scope, cursor, batchSize)
		if err != nil {
			return []dao.Object{}, fmt.Errorf("json to structs %v", err)
		}
		allObjects = append(allObjects, batch...)
	}

	return allObjects, nil
}

func getSeedsUsdPrice(ctx context.Context, api *eos.API) (eos.Asset, error) {
	var priceHistory []dao.SeedsPriceHistory
	var request eos.GetTableRowsRequest
	request.Code = "tlosto.seeds"
	request.Scope = "tlosto.seeds"
	request.Table = "pricehistory"
	request.Reverse = true
	request.Limit = 1
	request.JSON = true
	response, err := api.GetTableRows(ctx, request)
	if err != nil {
		return eos.Asset{}, fmt.Errorf("could not get Seeds/USD price: %v", err)
	}

	err = response.JSONToStructs(&priceHistory)
	if err != nil {
		return eos.Asset{}, fmt.Errorf("could not get Seeds/USD price: JSONToStructs: %v", err)
	}

	if len(priceHistory) == 0 {
		return eos.Asset{}, fmt.Errorf("Seeds/USD price queyr returned zero rows")
	}
	return priceHistory[0].SeedsUSD, nil
}

func yamlStringSettings() string {
	c := viper.AllSettings()
	bs, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("unable to marshal config to YAML: %v", err)
	}
	return string(bs)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Hypha prometheus metrics server",

	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		api := getAPI()

		log.Println(yamlStringSettings())

		errorCount := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "hypha_prometheus_errors",
			Help: "Number of errors that this instance has encountered",
		})

		hyphaSupply := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "hypha_balance_supply_hypha",
			Help: "Total amount of HYPHA tokens in circulation",
		})

		hvoiceSupply := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "hypha_balance_supply_hvoice",
			Help: "Total amount of HVOICE tokens rewarded",
		})

		husdSupply := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "hypha_balance_supply_husd",
			Help: "Total amount of HUSD tokens in circulation",
		})

		memberCount := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "hypha_dao_membership_members",
			Help: "Number of members enrolled in the DAO",
		})

		applicantCount := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "hypha_dao_membership_applicants",
			Help: "Number of applicants for the DHO",
		})

		openProposals := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "hypha_dao_proposals_total",
			Help: "Number of open proposals to be voted on",
		})

		seedsBalance := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "hypha_balance_balance_seeds",
			Help: "The SEEDS token balance of " + viper.GetString("DAOContract"),
		})

		escrowedSeedsBalance := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "hypha_balance_balance_seeds_escrow",
			Help: "The SEEDS token balance of " + viper.GetString("EscrowContract"),
		})

		hyphaSeedsAccountBalance := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "hypha_balance_balance_seeds_hyphaseedsaccount",
			Help: "The SEEDS token balance of " + viper.GetString("HyphaSeedsAccount"),
		})

		documentCount := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "hypha_graph_document_total",
			Help: "Total number of documents",
		})

		voteEventCount := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "hypha_dao_voteevents",
			Help: "Number of times users have voted in last 24 hours",
		})

		seedsPriceUsd := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "hypha_price_seedsusd",
			Help: "SEEDS/USD price",
		})

		// Periodically update the metrics
		go func() {
			for {
				sup, err := models.GetHvoiceSupply(ctx, api)
				if err == nil {
					hvoiceSupply.Set(assetToFloat(sup))
				} else {
					errorCount.Add(1)
					log.Println("Retrieval error: telos decide supply: telosdecide: "+viper.GetString("TelosDecideContract")+" symbol: "+viper.GetString("VoteTokenSymbol"), err)
				}

				hypha, err := getTokenSupply(ctx, api, viper.GetString("RewardToken.Contract"), viper.GetString("RewardToken.Symbol"))
				if err == nil {
<<<<<<< Updated upstream
					hyphaSupply.Set(assetToFloat(&hypha))
=======
					hyphaSupply.Set(hypha)
>>>>>>> Stashed changes
				} else {
					errorCount.Add(1)
					log.Println("Retrieval error: supply: token contract: "+viper.GetString("RewardToken.Contract")+" symbol: "+viper.GetString("RewardToken.Symbol"), err)
				}

				husd, err := getTokenSupply(ctx, api, viper.GetString("Treasury.TokenContract"), viper.GetString("Treasury.Symbol"))
				if err == nil {
<<<<<<< Updated upstream
					husdSupply.Set(assetToFloat(&husd))
=======
					husdSupply.Set(husd)
>>>>>>> Stashed changes
				} else {
					errorCount.Add(1)
					log.Println("Retrieval error: supply: token contract: "+viper.GetString("Treasury.TokenContract")+" symbol: "+viper.GetString("Treasury.Symbol"), err)
				}

				seeds, err := getTokenBalance(ctx, api, viper.GetString("SeedsTokenContract"), viper.GetString("DAOContract"), "SEEDS")
				if err == nil {
					seedsBalance.Set(assetToFloat(&seeds))
				} else {
					errorCount.Add(1)
					log.Println("Retrieval error: balance: "+viper.GetString("DAOContract")+" token contract: "+viper.GetString("SeedsTokenContract")+" symbol: SEEDS", err)
				}

				escrowedSeeds, err := getTokenBalance(ctx, api, viper.GetString("SeedsTokenContract"), viper.GetString("EscrowContract"), "SEEDS")
				if err == nil {
<<<<<<< Updated upstream
					escrowedSeedsBalance.Set(assetToFloat(&escrowedSeeds))
				} else {
					errorCount.Add(1)
					log.Println("Retrieval error: balance: "+viper.GetString("EscrowContract")+" token contract: "+viper.GetString("SeedsTokenContract")+" symbol: SEEDS", err)
				}

				hyphaSeedsAccountSeeds, err := getTokenBalance(ctx, api, viper.GetString("SeedsTokenContract"), viper.GetString("HyphaSeedsAccount"), "SEEDS")
				if err == nil {
					hyphaSeedsAccountBalance.Set(assetToFloat(&hyphaSeedsAccountSeeds))
				} else {
					errorCount.Add(1)
					log.Println("Retrieval error: balance: "+viper.GetString("HyphaSeedsAccount")+" token contract: "+viper.GetString("SeedsTokenContract")+" symbol: SEEDS", err)
				}

				members := models.Members(ctx, api)
				memberCount.Set(float64(len(members)))

				applicants := models.Applicants(ctx, api)
				applicantCount.Set(float64(len(applicants)))

				proposals, err := getLegacyObjects(ctx, api, eos.AN(viper.GetString("DAOContract")), eos.Name("proposal"))
				if err == nil {
					openProposals.Set(float64(len(proposals)))
=======
					voteEventCount.Set(float64(len(results)))
>>>>>>> Stashed changes
				} else {
					errorCount.Add(1)
					log.Println("Retrieval error: an error querying legacy objects from "+viper.GetString("DAOContract")+" scope: proposal", err)
				}

				seedsPriceUsdAsset, err := getSeedsUsdPrice(ctx, api)
				if err == nil {
					seedsPriceUsd.Set(assetToFloat(&seedsPriceUsdAsset))
				} else {
					errorCount.Add(1)
					log.Println("Retrieval error: an error querying the total number of documents from "+viper.GetString("DAOContract"), err)
				}

				docs, err := docgraph.GetAllDocuments(ctx, api, eos.AN(viper.GetString("DAOContract")))
				if err == nil {
					documentCount.Set(float64(len(docs)))
				} else {
					errorCount.Add(1)
					log.Println("Retrieval error: an error querying the total number of documents from "+viper.GetString("DAOContract"), err)
				}

				query := hyperion.NewQuery("castvote", viper.GetString("TelosDecideContract"), "")
				query.After = time.Now().AddDate(0, 0, -1)
				query.Limit = 1000
				results, err := query.Results()
				if err == nil {
<<<<<<< Updated upstream
					voteEventCount.Set(float64(len(results)))
=======
					daoEventCount.Set(float64(len(results)))
>>>>>>> Stashed changes
				} else {
					errorCount.Add(1)
					log.Println("Error updating the vote count: ", err)
				}

				log.Println("Sleeping for : " + viper.GetDuration("ScrapeInterval").String())
				time.Sleep(viper.GetDuration("ScrapeInterval"))
			}
		}()

		r := prometheus.NewRegistry()
		r.MustRegister(errorCount)
		r.MustRegister(hyphaSupply)
		r.MustRegister(hvoiceSupply)
		r.MustRegister(husdSupply)
		r.MustRegister(memberCount)
		r.MustRegister(applicantCount)
		r.MustRegister(openProposals)
		r.MustRegister(seedsBalance)
		r.MustRegister(escrowedSeedsBalance)
		r.MustRegister(hyphaSeedsAccountBalance)
		r.MustRegister(documentCount)
		r.MustRegister(voteEventCount)
		r.MustRegister(seedsPriceUsd)

		http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
		log.Println("Listening for requests on : " + viper.GetString("ServePort"))
		log.Fatal(http.ListenAndServe(":"+viper.GetString("ServePort"), nil))
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
