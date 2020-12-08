package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/views"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getDocumentCmd = &cobra.Command{
	Use:   "document [hash]",
	Short: "retrieve document details",
	Long:  "retrieve the detailed content within a document",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		hash := args[0]

		document, err := docgraph.LoadDocument(ctx, api, eos.AN(viper.GetString("DAOContract")), hash)
		if err != nil {
			panic("Document not found: " + hash)
		}

		jsonDoc, _ := json.MarshalIndent(document, "", "  ")

		fmt.Println("\nDocument Details")
		fmt.Println()
		fmt.Println(string(jsonDoc))
		fmt.Println()

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
	},
}

// func toString(d *models.Document) string {

// 	var assetStr []string
// 	for key, element := range d.Assets {
// 		assetStr = append(assetStr, fmt.Sprintf(key, "|%v", util.FormatAsset(&element, 2)))
// 	}

// 	output := []string{
// 		fmt.Sprintf("Doc ID|%v", strconv.Itoa(int(d.ID))),
// 		fmt.Sprintf("Prior ID|%v", strconv.Itoa(int(r.PriorID))),
// 		fmt.Sprintf("Owner|%v", string(d.Names["owner"])),
// 		fmt.Sprintf("Title|%v", string(d.Strings["title"])),
// 		fmt.Sprintf("Ballot|%v", string(d.Names["ballot_id"])),
// 		fmt.Sprintf("Assets|%v", assetStr),

// 		fmt.Sprintf("Minimum Time Commitment|%v", strconv.FormatFloat(r.MinTime*100, 'f', -1, 64)),
// 		fmt.Sprintf("Minimum Deferred Pay|%v", strconv.FormatFloat(r.MinDeferred*100, 'f', -1, 64)),
// 		fmt.Sprintf("Full Time Capacity|%v", strconv.FormatFloat(r.FullTimeCapacity, 'f', 1, 64)),
// 		fmt.Sprintf("FTE Cap Cost|%v", util.FormatAsset(&fteCapCost, 2)),
// 		fmt.Sprintf("Start Period|%v", r.StartPeriod.StartTime.Time.Format("2006 Jan 02 15:04:05")),
// 		fmt.Sprintf("End Period|%v", r.EndPeriod.EndTime.Time.Format("2006 Jan 02 15:04:05")),
// 		fmt.Sprintf("Created Date|%v", r.CreatedDate.Time.Format("2006 Jan 02 15:04:05")),
// 		fmt.Sprintf("Ballot ID|%v", string(r.BallotName)[10:]),
// 		fmt.Sprintf("Description|%v", r.Description),
// 	}
// 	return columnize.SimpleFormat(output)
// }

func init() {
	getCmd.AddCommand(getDocumentCmd)
}
