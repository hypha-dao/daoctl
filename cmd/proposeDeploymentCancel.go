package cmd

import (
	"context"
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type cancelProposal struct {
	ProposalName eos.Name `json:"proposal_name"`
	Canceler     eos.Name `json:"canceler"`
}

var proposeDeploymentCancelCmd = &cobra.Command{
	Use:   "cancel [proposal-name]",
	Short: "cancel an existing multisig deployment proposal",
	Long:  "cancel an existing multisig deployment proposal",
	RunE: func(cmd *cobra.Command, args []string) error {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		proposalName, err := grabInput("propose-deployment-cancel-cmd-proposal-name", proposalNamePromptLabel)
		if err != nil {
			return fmt.Errorf("cannot clone repo: %v", err)
		}

		action := eos.Action{
			Account: eos.AN(viper.GetString("MsigContract")),
			Name:    eos.ActN("cancel"),
			Authorization: []eos.PermissionLevel{
				{Actor: eos.AN(viper.GetString("DAOUser")), Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(cancelProposal{
				ProposalName: eos.Name(proposalName),
				Canceler:     eos.Name(viper.GetString("DAOUser")),
			}),
		}

		pushEOSCActions(ctx, api, &action)
		return nil
	},
}

func init() {
	proposeDeploymentCmd.AddCommand(proposeDeploymentCancelCmd)
	proposeDeploymentCancelCmd.Flags().StringP("proposal-name", "", "", proposalNamePromptLabel)
}
