package cmd

import (
	"context"
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/views"
	"github.com/hypha-dao/document-graph/docgraph"
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

		docs, err := docgraph.GetAllDocuments(ctx, api, eos.AN(viper.GetString("DAOContract")))
		if err != nil {
			panic(fmt.Errorf("cannot get all documents: %v", err))
		}

		docsTable := views.DocTable(docs)
		docsTable.SetStyle(simpletable.StyleCompactLite)
		fmt.Println("\n" + docsTable.String() + "\n\n")
	},
}

func init() {
	getCmd.AddCommand(getDocumentsCmd)
}
