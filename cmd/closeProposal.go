package cmd

import (
  "context"
  "fmt"
  eos "github.com/eoscanada/eos-go"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "strconv"
)

func newCloseAction(proposal_id uint64) *eos.Action {
	return &eos.Action{
		Account: eos.AN(viper.GetString("DAOContract")),
		Name:    eos.ActN("closeprop"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AN(viper.GetString("DAOUser")), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(proposal_id),
	}
}

var closeCmd = &cobra.Command{
	Use:   "close [proposal id]",
	Short: "close a proposal",
	Long:  "close a proposal that is linked to a ballot where the voting period has ended",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {


		proposalID, err := strconv.ParseUint(args[0], 10,64)
		if err != nil {
		  fmt.Println("Error reading proposal_id. ", err)
		  return
    }

		pushEOSCActions(context.Background(), getAPI(), newCloseAction(uint64(proposalID)))
	},
}

func init() {
	RootCmd.AddCommand(closeCmd)
	closeCmd.Flags().IntP("proposal_id", "p", -1, "proposal ID to be closed")
}
