package cmd

import (
	"github.com/spf13/cobra"
)

// treasuryCmd represents the treasury
var treasuryCmd = &cobra.Command{
	Use:   "treasury",
	Short: "multi-chain token/value redemption module",
}

func init() {
	RootCmd.AddCommand(treasuryCmd)
	treasuryCmd.Flags().StringP("memo", "m", "", "memo to be added to the payment record on chain")
}
