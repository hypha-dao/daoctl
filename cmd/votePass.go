package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var votePassCmd = &cobra.Command{
	Use:   "<ballot_id>",
	Short: "vote pass on a ballot",
	Long:  "vote pass on a ballot",
	// Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		// ctx := context.Background()
		// contract := toAccount(viper.GetString("TelosDecideContract"), "contract")
		// action := toActionName("castvote", "action")
		fmt.Println("vote pass <ballot_id> is not yet implemented")

	},
}

func init() {
	voteCmd.AddCommand(createCmd)
}
