package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [options] [type]",
	Short: "get objects from the DAO on-chain smart contract",
}

func init() {
	RootCmd.AddCommand(getCmd)

	//  getCmd.Flags().StringP("type", "t", "role", "Type of object to retrieve")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
