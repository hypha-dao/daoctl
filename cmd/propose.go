package cmd

import (
	"github.com/spf13/cobra"
)

var proposeCmd = &cobra.Command{
	Use:   "propose [role]",
	Short: "propose a role",
}

func init() {
	RootCmd.AddCommand(proposeCmd)
	proposeCmd.PersistentFlags().StringP("file", "f", "", "filename of document's JSON file")
}
