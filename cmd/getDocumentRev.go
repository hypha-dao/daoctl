package cmd

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getDocumentRevCmd = &cobra.Command{
	Use:   "rev [edge-name]",
	Short: "traverse a reverse edge of the document",
	Long:  "traverse a reverse edge of the document",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		lastHashBytes, err := ioutil.ReadFile("last-doc.tmp")
		if err != nil {
			panic("Missing last-doc.tmp file")
		}
		lastHash := string(lastHashBytes)

		lastDocument, err := docgraph.LoadDocument(ctx, api, eos.AN(viper.GetString("DAOContract")), lastHash)
		if err != nil {
			panic("Document not found: " + lastHash)
		}

		var edge eos.Name
		if len(args) == 0 {
			toEdges, err := docgraph.GetEdgesToDocument(ctx, api, eos.AN(viper.GetString("DAOContract")), lastDocument)
			if err != nil {
				fmt.Println("ERROR: Cannot get edges from node: ", err)
			}
			colorCyan := "\033[36m"
			colorReset := "\033[0m"
			fmt.Println(string(colorCyan), "NOTE: <edge-name> argument not provided; defaulting to first one in list: "+string(toEdges[0].EdgeName))
			fmt.Println(string(colorReset))
			edge = toEdges[0].EdgeName
		} else {
			edge = eos.Name(args[0])
		}

		edges, err := docgraph.GetEdgesToDocumentWithEdge(ctx, api, eos.AN(viper.GetString("DAOContract")), lastDocument, eos.Name(edge))
		if err != nil || len(edges) <= 0 {
			panic("There are no edges from " + lastHash + " named " + string(edge))
		}

		document, err := docgraph.LoadDocument(ctx, api, eos.AN(viper.GetString("DAOContract")), edges[0].FromNode.String())
		if err != nil {
			panic("Reverse node not found: " + edges[0].FromNode.String())
		}

		// printDocument(ctx, api, document)

		err = ioutil.WriteFile("last-doc.tmp", []byte(document.Hash.String()), 0644)
	},
}

func init() {
	getDocumentCmd.AddCommand(getDocumentRevCmd)
}
