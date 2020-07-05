package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var setCmd = &cobra.Command{
	Use:   "set [options] [type]",
	Short: "set data on the DAO on-chain smart contract",
}

func init() {
	RootCmd.AddCommand(setCmd)
}
