package cmd

import (
	"context"
	"fmt"
	"strconv"

	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// type updateDoc struct {
// 	Hash    eos.Checksum256    `json:"hash"`
// 	Updater eos.AccountName    `json:"updater"`
// 	Group   string             `json:"group"`
// 	Key     string             `json:"key"`
// 	Value   docgraph.FlexValue `json:"value"`
// }

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "temporary action",
	Long:  "temporary action",
	RunE: func(cmd *cobra.Command, args []string) error {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()
		contract := eos.AN(viper.GetString("DAOContract"))

		docs, err := getAssignments(ctx, api, contract)
		if err != nil {
			return fmt.Errorf("cannot get all documents: %v", err)
		}

		ignoreCount := 0
		updateCount := 0
		investigateCount := 0

		for i, doc := range docs {

			fmt.Println("Assignment document (", strconv.Itoa(i), " of ", strconv.Itoa(len(docs)), "):  ", doc.Hash.String())
			edges, err := docgraph.GetEdgesToDocumentWithEdge(ctx, api, contract, doc, eos.Name("assignment"))
			if err != nil {
				return fmt.Errorf("cannot get role edges from document: %v", err)
			}

			roleFromOriginalProposal, err := doc.GetContentFromGroup("details", "role")
			if err != nil {
				return fmt.Errorf("cannot get role edges from document: %v", err)
			}
			fmt.Println("- role	from details.role	: ", roleFromOriginalProposal.String())

			if len(edges) == 0 {
				fmt.Println("- no edges from role to assignment	: ", doc.Hash.String())
				investigateCount++
			} else {
				fmt.Println("- role	from edge		: ", edges[0].ToNode.String())
				if edges[0].ToNode.String() != roleFromOriginalProposal.String() {
					fmt.Println("- details.role does not match inbound assignment edge role: OVERWRITE")
					updateCount++
				} else {
					fmt.Println("- hashes match")
					ignoreCount++
				}

			}

			fmt.Println()

			// action := eos.Action{
			// 	Account: eos.AN(viper.GetString("DAOContract")),
			// 	Name:    eos.ActN("updatedoc"),
			// 	Authorization: []eos.PermissionLevel{
			// 		{Actor: eos.AN(viper.GetString("DAOContract")), Permission: eos.PN("active")},
			// 	},
			// 	ActionData: eos.NewActionData(updateDoc{
			// 		Hash:    doc.Hash,
			// 		Updater: eos.AN(viper.GetString("DAOContract")),
			// 		Group:   "details",
			// 		Key:     "role",
			// 		Value: docgraph.FlexValue{
			// 			BaseVariant: eos.BaseVariant{
			// 				TypeID: docgraph.GetVariants().TypeID("checksum256"),
			// 				Impl:   edges[0].ToNode,
			// 			},
			// 		}}),
			// }

			//pushEOSCActions(ctx, api, &action)
		}

		fmt.Println("Total : ", strconv.Itoa(len(docs)))
		fmt.Println("Update	: ", strconv.Itoa(updateCount))
		fmt.Println("Ignore	: ", strconv.Itoa(ignoreCount))
		fmt.Println("Investigate: ", strconv.Itoa(investigateCount))

		return nil
	},
}

func init() {
	RootCmd.AddCommand(fixCmd)
}

func getAssignments(ctx context.Context, api *eos.API, contract eos.AccountName) ([]docgraph.Document, error) {

	allDocuments, err := docgraph.GetAllDocuments(ctx, api, eos.AN(viper.GetString("DAOContract")))
	if err != nil {
		return []docgraph.Document{}, fmt.Errorf("cannot get all documents: %v", err)
	}

	var filteredDocs []docgraph.Document
	for _, doc := range allDocuments {

		typeFV, err := doc.GetContent("type")
		if err == nil &&
			typeFV.Impl.(eos.Name) == eos.Name("assignment") {
			filteredDocs = append(filteredDocs, doc)
		}
	}
	return filteredDocs, nil
}
