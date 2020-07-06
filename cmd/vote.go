package cmd

import (
	"context"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var voteCmd = &cobra.Command{
	Use:   "vote ballot_id pass|fail",
	Short: "vote pass or fail on a ballot",
	Long:  "vote on a specific ballot.  Example:  daoctl vote 34 pass",
	Args:  cobra.RangeArgs(2, 2),
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("vote pass <ballot_id> is not yet implemented")

		ballotName := eos.Name(viper.GetString("BallotPrefix") + args[0]) // TODO: this will break; need to make the prefix dynamic
		option := eos.Name(args[1])                                       // only supporting a single option value, will enhance to include multi-value later

		pushEOSCActions(context.Background(), getAPI(), newVoteAction(ballotName, option))
	},
}

// Vote represents a set of options being cast as a vote to Telos Decide
type Vote struct {
	Voter      eos.Name   `json:"voter"`
	BallotName eos.Name   `json:"ballot_name"`
	Options    []eos.Name `json:"options"`
}

func newVoteAction(ballotName, option eos.Name) *eos.Action {

	return &eos.Action{
		Account: eos.AN(viper.GetString("TelosDecideContract")),
		Name:    eos.ActN("castvote"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AN(viper.GetString("DAOUser")), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(&Vote{
			Voter:      eos.Name(viper.GetString("DAOUser")),
			BallotName: ballotName,
			Options:    []eos.Name{option},
		}),
	}
}

func init() {
	RootCmd.AddCommand(voteCmd)
}
