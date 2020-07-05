package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "creates a local backup of current DAO data",
	Long:  "creates a new directory containing JSON files for all DAO object types",
	// Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		api := getAPI()
		ctx := context.Background()

		folderName := viper.GetString("backup-cmd-output-dir") + "/dao-backup-" + time.Now().Format("2006Jan02-150405")

		err := os.Mkdir(folderName, 0777)
		if err != nil {
			fmt.Println("Unable to create folder: ", folderName)
			panic(err)
		}
		fmt.Println("\nBacking up to folder: ", folderName)

		var request eos.GetTableByScopeRequest
		request.Code = viper.GetString("DAOContract")
		request.Table = "objects"
		request.Limit = 500 // maximum number of scopes
		response, err := api.GetTableByScope(context.Background(), request)
		errorCheck("get table by scope", err)

		var scopes []models.Scope
		json.Unmarshal(response.Rows, &scopes)

		for _, scope := range scopes {
			saveObjects(ctx, api, folderName, string(scope.Scope), "objects")
		}

		saveObjects(ctx, api, folderName, viper.GetString("DAOContract"), "periods")
		saveObjects(ctx, api, folderName, viper.GetString("DAOContract"), "config")
		saveObjects(ctx, api, folderName, viper.GetString("DAOContract"), "applicants")
		saveObjects(ctx, api, folderName, viper.GetString("DAOContract"), "members")
		saveObjects(ctx, api, folderName, viper.GetString("DAOContract"), "payments")
	},
}

func saveObjects(ctx context.Context, api *eos.API, folderName, scope, table string) {

	filename := folderName + "/" + scope + "_" + table + ".json"

	var request eos.GetTableRowsRequest
	request.Code = viper.GetString("DAOContract")
	request.Scope = scope
	request.Table = table
	request.Limit = 1000
	request.JSON = true
	response, err := api.GetTableRows(ctx, request)
	if err != nil {
		fmt.Println("Unable to retrieve rows: ", scope)
		panic(err)
	}

	data, err := response.Rows.MarshalJSON()
	if err != nil {
		fmt.Println("Unable to backup scope: ", scope, ", table: ", table)
		panic(err)
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println("Unable to write file: ", filename, " for scope: ", scope, ", table: ", table)
		panic(err)
	}
}

func init() {
	backupCmd.Flags().StringP("output-dir", "", "./", "directory location to save the backup folder")
	RootCmd.AddCommand(backupCmd)
}
