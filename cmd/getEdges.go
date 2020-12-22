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

var getEdgesCmd = &cobra.Command{
	Use:   "edges",
	Short: "query edges",
	Long:  "query edges",
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		edges, err := docgraph.GetAllEdges(ctx, api, eos.AN(viper.GetString("DAOContract")))
		if err != nil {
			panic(fmt.Errorf("cannot get all edges: %v", err))
		}

		edgesTable := views.EdgeTable(edges, false, false)
		edgesTable.SetStyle(simpletable.StyleCompactLite)
		fmt.Println("\n" + edgesTable.String() + "\n\n")
	},
}

func init() {
	getCmd.AddCommand(getEdgesCmd)
}
