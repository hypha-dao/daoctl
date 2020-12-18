package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/views"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getDocumentNextCmd = &cobra.Command{
	Use:   "next [edge-name]",
	Short: "retrieve the next document in the graph based on the edge",
	Long:  "retrieve the next document in the graph based on the edge",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		edge := args[0]

		lastHashBytes, err := ioutil.ReadFile("last-doc.tmp")
		if err != nil {
			panic("Missing last-doc.tmp file")
		}
		lastHash := string(lastHashBytes)

		lastDocument, err := docgraph.LoadDocument(ctx, api, eos.AN(viper.GetString("DAOContract")), lastHash)
		if err != nil {
			panic("Document not found: " + lastHash)
		}

		edges, err := docgraph.GetEdgesFromDocumentWithEdge(ctx, api, eos.AN(viper.GetString("DAOContract")), lastDocument, eos.Name(edge))
		if err != nil || len(edges) <= 0 {
			panic("There are no edges from " + lastHash + " named " + edge)
		}

		document, err := docgraph.LoadDocument(ctx, api, eos.AN(viper.GetString("DAOContract")), edges[0].ToNode.String())
		if err != nil {
			panic("Next document not found: " + edges[0].ToNode.String())
		}

		fmt.Println("Document Details")

		fmt.Println()
		output := []string{
			fmt.Sprintf("ID|%v", strconv.Itoa(int(document.ID))),
			fmt.Sprintf("Hash|%v", document.Hash.String()),
			fmt.Sprintf("Creator|%v", string(document.Creator)),
			fmt.Sprintf("Created Date|%v", document.CreatedDate.Time.Format("2006 Jan 02")),
		}

		fmt.Println(columnize.SimpleFormat(output))
		fmt.Println()
		printContentGroups(&document)

		fromEdges, err := docgraph.GetEdgesFromDocument(ctx, api, eos.AN(viper.GetString("DAOContract")), document)
		if err != nil {
			fmt.Println("ERROR: Cannot get edges from document: ", err)
		}

		fromEdgesTable := views.EdgeTable(fromEdges)
		fromEdgesTable.SetStyle(simpletable.StyleCompactLite)
		fmt.Println("\n" + fromEdgesTable.String() + "\n\n")

		toEdges, err := docgraph.GetEdgesToDocument(ctx, api, eos.AN(viper.GetString("DAOContract")), document)
		if err != nil {
			fmt.Println("ERROR: Cannot get edges to document: ", err)
		}

		toEdgesTable := views.EdgeTable(toEdges)
		toEdgesTable.SetStyle(simpletable.StyleCompactLite)
		fmt.Println("\n" + toEdgesTable.String() + "\n\n")

		err = ioutil.WriteFile("last-doc.tmp", []byte(document.Hash.String()), 0644)
	},
}

func init() {
	getDocumentCmd.AddCommand(getDocumentNextCmd)
}
