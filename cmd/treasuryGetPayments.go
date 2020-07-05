package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var treasuryGetPaymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "view a table of payments",
	//Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		// ctx := context.Background()
		// contract := toAccount(viper.GetString("TelosDecideContract"), "contract")
		// action := toActionName("castvote", "action")
		fmt.Println("not yet implemented")

	},
}

func init() {
	treasuryGetCmd.AddCommand(treasuryGetPaymentsCmd)
}
