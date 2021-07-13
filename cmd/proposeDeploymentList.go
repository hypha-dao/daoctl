package cmd

import (
	"context"
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type depProposal struct {
	Proposer     eos.Name        `json:"proposer"`
	ProposalName eos.Name        `json:"proposal_name"`
	DocumentHash eos.Checksum256 `json:"document_hash"`
}

var proposeDeploymentListCmd = &cobra.Command{
	Use:   "list",
	Short: "list open deployment proposals",
	Long:  "list open deployment proposals",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		deps, err := getProposedDeployments()
		if err != nil {
			return fmt.Errorf("cannot get list of proposed deployments %v", err)
		}

		dT, err := depsTable(ctx, getAPI(), deps)
		if err != nil {
			return fmt.Errorf("cannot construct deployment proposal table %v", err)
		}

		fmt.Println(dT.String())
		fmt.Println()
		fmt.Println()
		return nil
	},
}

func init() {
	proposeDeploymentCmd.AddCommand(proposeDeploymentListCmd)
}

func depsHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Proposal Name"},
			{Align: simpletable.AlignCenter, Text: "Developer"},
			{Align: simpletable.AlignCenter, Text: "Proposer"},
			{Align: simpletable.AlignCenter, Text: "Commit"},
			{Align: simpletable.AlignCenter, Text: "Notes"},
		},
	}
}

func depsTable(ctx context.Context, api *eos.API, deps []depProposal) (*simpletable.Table, error) {
	table := simpletable.New()
	table.Header = depsHeader()

	for _, dep := range deps {

		propDoc, err := docgraph.LoadDocument(ctx, api, eos.AN(viper.GetString("MsigContract")), dep.DocumentHash.String())
		if err != nil {
			return &simpletable.Table{}, fmt.Errorf("error retrieving document hash: %v", err)
		}

		devFv, err := propDoc.ContentGroups[0].GetContent("developer")
		if err != nil {
			return &simpletable.Table{}, fmt.Errorf("error converting flex value to string: %v", err)
		}

		commitFv, err := propDoc.ContentGroups[0].GetContent("github_commit")
		if err != nil {
			return &simpletable.Table{}, fmt.Errorf("error converting flex value to string: %v", err)
		}

		notesFv, err := propDoc.ContentGroups[0].GetContent("notes")
		if err != nil {
			return &simpletable.Table{}, fmt.Errorf("error converting flex value to string: %v", err)
		}

		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: string(dep.ProposalName)},
			{Align: simpletable.AlignLeft, Text: devFv.String()},
			{Align: simpletable.AlignLeft, Text: string(dep.Proposer)},
			{Align: simpletable.AlignLeft, Text: commitFv.String()},
			{Align: simpletable.AlignLeft, Text: notesFv.String()},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.SetStyle(simpletable.StyleCompactLite)
	return table, nil
}
