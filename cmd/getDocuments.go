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

var getDocumentsCmd = &cobra.Command{
	Use:   "documents [scope]",
	Short: "query and manage documents",
	Long:  "query and manage documents",
	// Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		docs := models.LoadDocuments(ctx, api, viper.GetString("get-documents-cmd-scope"))

		docsTable := views.DocTable(docs)
		docsTable.SetStyle(simpletable.StyleCompactLite)
		fmt.Println("\n" + docsTable.String() + "\n\n")
	},
}

func init() {
	getCmd.AddCommand(getDocumentsCmd)
	getDocumentsCmd.Flags().StringP("scope", "", "proposal", "document scope used to query the on-chain object table")
}
