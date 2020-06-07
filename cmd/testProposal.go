package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var testProposalCmd = &cobra.Command{
	Use:   "",
	Short: "test the proposal object",
	Long:  "run some tests on the smart contract",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running tests for the proposal")
		return
	},
}

func init() {
	testCmd.AddCommand(testProposalCmd)
}
