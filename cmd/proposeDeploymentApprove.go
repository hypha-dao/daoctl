package cmd

import (
	"context"
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type approveProposal struct {
	Proposer     eos.Name            `json:"proposer"`
	ProposalName eos.Name            `json:"proposal_name"`
	Permission   eos.PermissionLevel `json:"level"`
}

var proposeDeploymentApproveCmd = &cobra.Command{
	Use:   "approve [proposal-name]",
	Short: "approve an existing multisig deployment proposal",
	Long:  "approve an existing multisig deployment proposal",
	RunE: func(cmd *cobra.Command, args []string) error {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		proposalName, err := grabInput("propose-deployment-approve-cmd-proposal-name", proposalNamePromptLabel)
		if err != nil {
			return fmt.Errorf("cannot read proposal-name: %v %v", proposalNamePromptLabel, err)
		}

		action := eos.Action{
			Account: eos.AN("eosio.msig"),
			Name:    eos.ActN("approve"),
			Authorization: []eos.PermissionLevel{
				{Actor: eos.AN(viper.GetString("DAOUser")), Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(approveProposal{
				Proposer:     eos.Name(viper.GetString("DAOUser")),
				ProposalName: eos.Name(proposalName),
				Permission: eos.PermissionLevel{
					Actor:      eos.AN(viper.GetString("DAOUser")),
					Permission: eos.PN("active")},
			}),
		}

		pushEOSCActions(ctx, api, &action)
		return nil
	},
}

func init() {
	proposeDeploymentCmd.AddCommand(proposeDeploymentApproveCmd)
	proposeDeploymentApproveCmd.Flags().StringP("proposal-name", "", "", proposalNamePromptLabel)
}
