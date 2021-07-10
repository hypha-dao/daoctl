package cmd

import (
	"context"
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type hyphaCancelProposal struct {
	ProposalName eos.Name `json:"proposal_name"`
	Canceler     eos.Name `json:"canceler"`
}

type eosioCancelProposal struct {
	Proposer     eos.Name `json:"proposer"`
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

		deps, err := getProposedDeployments()
		if err != nil {
			return fmt.Errorf("cannot get list of proposed deployments %v", err)
		}

		found := false
		var proposer eos.Name

		for _, dep := range deps {
			if string(dep.ProposalName) == proposalName {
				found = true
				proposer = dep.Proposer
				break
			}
		}

		if !found {
			zlog.Error("cannot find proposal in the hypha msig contract",
				zap.String("contract", viper.GetString("MsigContract")),
				zap.String("proposal-name", proposalName))
		}

		eosioAction := eos.Action{
			Account: eos.AN("eosio.msig"),
			Name:    eos.ActN("cancel"),
			Authorization: []eos.PermissionLevel{
				{Actor: eos.AN(viper.GetString("DAOUser")), Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(eosioCancelProposal{
				Proposer:     proposer,
				ProposalName: eos.Name(proposalName),
				Canceler:     eos.Name(viper.GetString("DAOUser"))}),
		}

		hyphaAction := eos.Action{
			Account: eos.AN(viper.GetString("MsigContract")),
			Name:    eos.ActN("cancel"),
			Authorization: []eos.PermissionLevel{
				{Actor: eos.AN(viper.GetString("DAOUser")), Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(hyphaCancelProposal{
				ProposalName: eos.Name(proposalName),
				Canceler:     eos.Name(viper.GetString("DAOUser"))}),
		}

		pushEOSCActions(ctx, api, &eosioAction, &hyphaAction)
		return nil
	},
}

func init() {
	proposeDeploymentCmd.AddCommand(proposeDeploymentCancelCmd)
	proposeDeploymentCancelCmd.Flags().StringP("proposal-name", "", "", proposalNamePromptLabel)
}
