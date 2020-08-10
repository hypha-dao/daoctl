package cmd

import (
	"context"
	"fmt"
	"strconv"

	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getDocumentCmd = &cobra.Command{
	Use:   "document [document id]",
	Short: "retrieve document details",
	Long:  "retrieve the detailed about a document",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()
		//ac := accounting.NewAccounting("", 0, ",", ".", "%s %v", "%s (%v)", "%s --") // TODO: make this configurable

		documentID, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			fmt.Println("Parse error: Document id must be a positive integer (uint64)")
			return
		}
		document := models.LoadDocument(ctx, api, viper.GetString("get-document-cmd-scope"), documentID)

		fmt.Println("\n\nDocument Details")
		fmt.Println("Scope: ", viper.GetString("get-document-cmd-scope"), "; ID: ", documentID)
		fmt.Println()
		fmt.Println(document)
		fmt.Println()
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
// 		fmt.Sprintf("Ballot ID|%v", string(r.BallotName)[11:]),
// 		fmt.Sprintf("Description|%v", r.Description),
// 	}
// 	return columnize.SimpleFormat(output)
// }

func init() {
	getDocumentCmd.Flags().StringP("scope", "", "proposal", "document scope used to query the on-chain object table")

	getCmd.AddCommand(getDocumentCmd)
}
