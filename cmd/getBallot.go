package cmd

import (
	"context"
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/daoctl/views"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getBallotCmd = &cobra.Command{
	Use:   "ballot [ballot name]",
	Short: "retrieves ballot details",
	Long:  "retrieves the ballot times, voters, voting selections, and quorum info",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		ballotName := eos.Name("hypha1....." + args[0])

		ballot, err := models.NewBallot(ctx, api, ballotName)
		if err != nil {
			panic("Cannot read ballot: " + args[0])
		}

		fmt.Println("\n\n" + views.BallotHeader(*ballot) + "\n\n")
		fmt.Println(views.VotesTable(ballot.Votes))
	},
}

func init() {
	getCmd.AddCommand(getBallotCmd)
	// getRoleCmd.Flags().BoolP("include-proposals", "i", false, "include proposals in the output")
}
