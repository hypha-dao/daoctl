package cmd

import (
	"context"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newCloseAction(proposalID string) *eos.Action {
	return &eos.Action{
		Account: eos.AN(viper.GetString("DAOContract")),
		Name:    eos.ActN("closedocprop"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AN(viper.GetString("DAOUser")), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(proposalID),
	}
}

var closeCmd = &cobra.Command{
	Use:   "close [hash]",
	Short: "close a proposal",
	Long:  "close a proposal that is linked to a ballot where the voting period has ended",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {

		// proposalID, err := strconv.ParseUint(args[0], 10, 64)
		// if err != nil {
		// 	fmt.Println("Error reading proposal_id. ", err)
		// 	return
		// }

		pushEOSCActions(context.Background(), getAPI(), newCloseAction(args[0]))
	},
}

func init() {
	RootCmd.AddCommand(closeCmd)
	closeCmd.Flags().BoolP("all", "a", false, "attempt to close all open proposals")
}
