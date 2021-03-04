package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "creates a local backup of the environment, including all documents and edges",
	Long:  "creates a local backup of the environment, including all documents and edges",
	// Args:  cobra.RangeArgs(1, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		api := getAPI()
		ctx := context.Background()
		contract := eos.AN(viper.GetString("DAOContract"))

		folderName := viper.GetString("backup-cmd-output-dir") + "/dao-backup-" + time.Now().Format("2006Jan02-150405")

		err := os.Mkdir(folderName, 0777)
		if err != nil {
			fmt.Println("Unable to create folder: ", folderName)
			panic(err)
		}
		fmt.Println("\nBacking up to folder: ", folderName)

		documents, err := docgraph.GetAllDocuments(ctx, api, contract)
		if err != nil {
			return fmt.Errorf("cannot get all documents: %v", err)
		}

		documentsJson, err := json.MarshalIndent(documents, "", "  ")
		if err != nil {
			return fmt.Errorf("cannot marshal documents to json: %v", err)
		}

		err = ioutil.WriteFile("documents.json", documentsJson, 0644)
		if err != nil {
			return fmt.Errorf("cannot documents to documents.json file: %v", err)
		}

		edges, err := docgraph.GetAllEdges(ctx, api, contract)
		if err != nil {
			return fmt.Errorf("cannot get all edges: %v", err)
		}

		edgesJson, err := json.MarshalIndent(edges, "", "  ")
		if err != nil {
			return fmt.Errorf("cannot marshal edges to json: %v", err)
		}

		err = ioutil.WriteFile("edges.json", edgesJson, 0644)
		if err != nil {
			return fmt.Errorf("cannot edges to edges.json file: %v", err)
		}
		return nil
	},
}

func init() {
	backupCmd.Flags().StringP("output-dir", "", "./", "directory location to save the backup folder")
	RootCmd.AddCommand(backupCmd)
}
