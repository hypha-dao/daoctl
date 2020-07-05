package cmd

import (
	"github.com/spf13/cobra"
)

var voteCmd = &cobra.Command{
	Use:   "vote [ballot_id] [pass | fail]",
	Short: "vote pass or fail on a ballot",
}

func init() {
	RootCmd.AddCommand(voteCmd)
}
