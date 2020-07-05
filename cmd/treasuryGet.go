package cmd

import (
	"github.com/spf13/cobra"
)

// treasuryCmd represents the treasury
var treasuryGetCmd = &cobra.Command{
	Use:   "get",
	Short: "retrieve and display treasury objects",
}

func init() {
	treasuryCmd.AddCommand(treasuryGetCmd)
}
