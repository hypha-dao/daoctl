package cmd

import (
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/spf13/cobra"
)

type deployment struct {
	Proposer           eos.AccountName         `json:"proposer"`
	ProposalName       eos.Name                `json:"proposal_name"`
	RequestedApprovals []eos.PermissionLevel   `json:"requested"`
	ContentGroups      []docgraph.ContentGroup `json:"content_groups"`
	Transaction        *eos.Transaction        `json:"trx"`
}

var proposeDeploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "manage multisig deployment proposals",
	Long:  "manage multisig deployment proposals",
}

func init() {
	proposeCmd.AddCommand(proposeDeploymentCmd)
}
