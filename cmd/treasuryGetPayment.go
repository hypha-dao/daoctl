package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var treasuryGetPaymentCmd = &cobra.Command{
	Use:   "payment <redemption_id>",
	Short: "view the details of a specific treasury payment",
	//Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		// ctx := context.Background()
		// contract := toAccount(viper.GetString("TelosDecideContract"), "contract")
		// action := toActionName("castvote", "action")
		fmt.Println("not yet implemented")

	},
}

func init() {
	treasuryGetCmd.AddCommand(treasuryGetPaymentCmd)
}
