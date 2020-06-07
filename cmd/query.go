package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/daoctl/views"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"go.uber.org/zap"
)

// getCmd represents the get command
var queryCmd = &cobra.Command{
	Use:   "query [filter]",
	Short: "Query action history on the DAO for specific users",
	Run: func(cmd *cobra.Command, args []string) {
		// api := getAPI()
		// ctx := context.Background()
		accountName := viper.GetString("query-cmd-account")
		contractName := viper.GetString("query-cmd-contract")
		actionName := viper.GetString("query-cmd-action")
		trxID := viper.GetString("query-cmd-trx")
		limit := viper.GetInt("query-cmd-limit")

		zlog.Debug("Query parameters",
			zap.String("account", accountName),
			zap.String("contract", contractName),
			zap.String("action", actionName),
			zap.Int("limit", limit))

		request := viper.GetString("HyperionEndpoint") + "/history/"
		if trxID != "" {
			request += "get_transaction?id=" + trxID
		} else {
			request += "get_actions?limit=" + strconv.Itoa(limit)
			if actionName != "" {
				request += "&act.name=" + actionName
			}
			if contractName != "" {
				request += "&act.account=" + contractName
			}
			if accountName != "" {
				request += "&account=" + accountName
			}
		}

		zlog.Debug("Query request", zap.String("request", request))
		resp, err := http.Get(request)
		if err != nil {
			zlog.Panic(
				"http get query failed",
				zap.Error(err),
				zap.String("request", request),
			)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			zlog.Panic(
				"reading from get query response body failed",
				zap.Error(err),
				zap.String("response body", string(body)),
			)
		}

		if viper.GetBool("query-cmd-json") {
			jsonRaw := pretty.Color(pretty.Pretty(body), pretty.TerminalStyle)
			if err != nil {
				zlog.Warn(
					"unable to format JSON with color, using default",
					zap.Error(err))
				fmt.Println(string(body))
			}
			fmt.Println(string(jsonRaw))
			return
		}

		var actions []models.QrAction
		result := gjson.Get(string(body), "actions")
		result.ForEach(func(key, value gjson.Result) bool {
			actionTime, err := time.Parse("2006-01-02T15:04:05.000", getString(value, "timestamp"))
			if err != nil {
				zlog.Warn(
					"unable to parse timestamp - using blank",
					zap.Error(err),
					zap.String("timestamp attempted parsing: ", getString(value, "timestamp")),
				)
			}

			action := models.QrAction{
				Timestamp:      actionTime,
				TrxID:          gjson.Get(value.String(), "trx_id").String(),
				ActionContract: getString(value, "act.account"),
				ActionName:     getString(value, "act.name"),
				Data:           getString(value, "act.data"),
			}
			//fmt.Println("Action name: ", gjson.Get(value.String(), "act.name"))
			actions = append(actions, action)
			return true // keep iterating
		})

		actionTable := views.ActionQueryResultTable(actions)
		actionTable.SetStyle(simpletable.StyleCompactLite)

		fmt.Println(actionTable.String())
	},
}

func getString(result gjson.Result, element string) string {

	charsToShow := 45
	suffix := "... <snip>"
	longValue := gjson.Get(result.String(), element).String()

	if len(longValue) < charsToShow {
		charsToShow = len(longValue)
		suffix = ""
	}
	return longValue[:charsToShow] + suffix
}

func init() {
	RootCmd.AddCommand(queryCmd)
	queryCmd.Flags().StringP("account", "", "", "member's account to query")
	queryCmd.Flags().StringP("contract", "", "dao.hypha", "query actions called on this smart contract (defaults to DAO)")
	queryCmd.Flags().StringP("action", "", "", "action name to query on the contract (defaults to create)")
	queryCmd.Flags().StringP("trx", "", "", "transaction ID to query for the full content of that transaction")
	queryCmd.Flags().IntP("limit", "", 10, "maximum number of records to retrieve")
	queryCmd.Flags().BoolP("json", "j", false, "print the results as JSON")

}
