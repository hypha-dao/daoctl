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

var getDocumentFwdCmd = &cobra.Command{
	Use:   "fwd [edge-name]",
	Short: "traverse a forward edge of this document",
	Long:  "traverse a forward edge of this document",
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

		colorCyan := "\033[36m"
		colorReset := "\033[0m"
		// colorRed := "\033[31m"

		var edge eos.Name
		if len(args) == 0 {
			fromEdges, err := docgraph.GetEdgesFromDocument(ctx, api, eos.AN(viper.GetString("DAOContract")), lastDocument)
			if err != nil {
				fmt.Println("ERROR: Cannot get edges from document: ", err)
			}
			fmt.Println(string(colorCyan), "NOTE: <edge-name> argument not provided; defaulting to first one in list: "+string(fromEdges[0].EdgeName))
			fmt.Println(string(colorReset))
			edge = fromEdges[0].EdgeName
		} else {
			edge = eos.Name(args[0])
		}

		edges, err := docgraph.GetEdgesFromDocumentWithEdge(ctx, api, eos.AN(viper.GetString("DAOContract")), lastDocument, eos.Name(edge))
		if err != nil || len(edges) <= 0 {
			panic("There are no edges from " + lastHash + " named " + string(edge))
		}

		document, err := docgraph.LoadDocument(ctx, api, eos.AN(viper.GetString("DAOContract")), edges[0].ToNode.String())
		if err != nil {
			panic("Next document not found: " + edges[0].ToNode.String())
		}

		// printDocument(ctx, api, document)

		err = ioutil.WriteFile("last-doc.tmp", []byte(document.Hash.String()), 0644)
	},
}

func init() {
	getDocumentCmd.AddCommand(getDocumentFwdCmd)
}
