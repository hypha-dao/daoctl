package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/document-graph/docgraph"

	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getDefaultPeriod() docgraph.Document {

	ctx := context.Background()

	root, err := docgraph.LoadDocument(ctx, getAPI(), eos.AN(viper.GetString("DAOContract")), "0f374e7a9d8ab17f172f8c478744cdd4016497e15229616f2ffd04d8002ef64a")
	if err != nil {
		panic(err)
	}

	edges, err := docgraph.GetEdgesFromDocumentWithEdge(ctx, getAPI(), eos.AN(viper.GetString("DAOContract")), root, eos.Name("start"))
	if err != nil || len(edges) <= 0 {
		panic("There are no next edges")
	}

	lastDocument, err := docgraph.LoadDocument(ctx, getAPI(), eos.AN(viper.GetString("DAOContract")), edges[0].ToNode.String())
	if err != nil {
		panic("Next document not found: " + edges[0].ToNode.String())
	}

	index := 1
	for index < 6 {

		edges, err := docgraph.GetEdgesFromDocumentWithEdge(ctx, getAPI(), eos.AN(viper.GetString("DAOContract")), lastDocument, eos.Name("next"))
		if err != nil || len(edges) <= 0 {
			panic("There are no next edges")
		}

		lastDocument, err = docgraph.LoadDocument(ctx, getAPI(), eos.AN(viper.GetString("DAOContract")), edges[0].ToNode.String())
		if err != nil {
			panic("Next document not found: " + edges[0].ToNode.String())
		}

		index++
	}
	return lastDocument
}

var proposeAssignmentCmd = &cobra.Command{
	Use:   "assignment -f [filename] --role <hash>",
	Short: "propose an assignment",

	Run: func(cmd *cobra.Command, args []string) {

		ctx := context.Background()
		contract := toAccount(viper.GetString("DAOContract"), "contract")

		var role docgraph.Document
		var err error
		if len(viper.GetString("propose-assignment-cmd-role")) > 0 {
			role, err = docgraph.LoadDocument(ctx, getAPI(), contract, viper.GetString("propose-assignment-cmd-role"))
		} else {
			role, err = docgraph.GetLastDocumentOfEdge(ctx, getAPI(), contract, eos.Name("role"))
		}
		if err != nil {
			panic(err)
		}

		data, err := ioutil.ReadFile(viper.GetString("propose-cmd-file"))
		if err != nil {
			fmt.Println("Unable to read file: ", viper.GetString("propose-cmd-file"))
			return
		}

		var proposalDoc docgraph.Document
		err = json.Unmarshal([]byte(data), &proposalDoc)
		if err != nil {
			panic(err)
		}

		// inject the role hash in the first content group of the document
		proposalDoc.ContentGroups[0] = append(proposalDoc.ContentGroups[0], docgraph.ContentItem{
			Label: "role",
			Value: &docgraph.FlexValue{
				BaseVariant: eos.BaseVariant{
					TypeID: docgraph.GetVariants().TypeID("checksum256"),
					Impl:   role.Hash,
				}},
		})

		// inject the period hash in the first content group of the document
		proposalDoc.ContentGroups[0] = append(proposalDoc.ContentGroups[0], docgraph.ContentItem{
			Label: "start_period",
			Value: &docgraph.FlexValue{
				BaseVariant: eos.BaseVariant{
					TypeID: docgraph.GetVariants().TypeID("checksum256"),
					Impl:   getDefaultPeriod().Hash,
				}},
		})

		// inject the assignee in the first content group of the document
		proposalDoc.ContentGroups[0] = append(proposalDoc.ContentGroups[0], docgraph.ContentItem{
			Label: "assignee",
			Value: &docgraph.FlexValue{
				BaseVariant: eos.BaseVariant{
					TypeID: docgraph.GetVariants().TypeID("name"),
					Impl:   eos.Name(viper.GetString("DAOUser")),
				}},
		})

		action := eos.ActN("propose")
		actions := eos.Action{
			Account: contract,
			Name:    action,
			Authorization: []eos.PermissionLevel{
				{Actor: eos.AN(viper.GetString("DAOUser")), Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(proposal{
				Proposer:      eos.AN(viper.GetString("DAOUser")),
				ProposalType:  eos.Name("assignment"),
				ContentGroups: proposalDoc.ContentGroups,
			})}

		pushEOSCActions(context.Background(), getAPI(), &actions)
	},
}

func init() {
	proposeCmd.AddCommand(proposeAssignmentCmd)
	proposeAssignmentCmd.PersistentFlags().StringP("role", "", "", "hash of role to apply to")

}
