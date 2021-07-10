package cmd

import (
	"context"
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type hyphaProposal struct {
	Proposer      eos.AccountName         `json:"proposer"`
	ProposalName  eos.Name                `json:"proposal_name"`
	ContentGroups []docgraph.ContentGroup `json:"content_groups"`
}

type eosioProposal struct {
	Proposer           eos.AccountName       `json:"proposer"`
	ProposalName       eos.Name              `json:"proposal_name"`
	RequestedApprovals []eos.PermissionLevel `json:"requested"`
	Transaction        *eos.Transaction      `json:"trx"`
}

func getProposedDeployments() ([]depProposal, error) {
	var deps []depProposal
	var request eos.GetTableRowsRequest
	request.Code = viper.GetString("MsigContract")
	request.Scope = viper.GetString("MsigContract")
	request.Table = "proposals"
	request.JSON = true
	response, err := getAPI().GetTableRows(context.Background(), request)
	if err != nil {
		return []depProposal{}, fmt.Errorf("get table rows %v", err)
	}

	err = response.JSONToStructs(&deps)
	if err != nil {
		return []depProposal{}, fmt.Errorf("json to structs %v", err)
	}
	return deps, nil
}

var proposeDeploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "manage multisig deployment proposals",
	Long:  "manage multisig deployment proposals",
}

func init() {
	proposeCmd.AddCommand(proposeDeploymentCmd)
}
