package cmd

import (
	"context"
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/daoctl/views"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getAssignmentCmd = &cobra.Command{
	Use:   "assignments [account name]",
	Short: "retrieve assignments",
	Long:  "retrieve all active assignments For a json dump, append the argument --json.",
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		periods := models.LoadPeriods(api)
		roles := models.Roles(ctx, api, periods)

		if viper.GetBool("get-assignments-cmd-failed-proposals") == true {
			assignments := models.Assignments(ctx, api, roles, periods, "failedprops")
			assignmentsTable := views.AssignmentTable(assignments)
			assignmentsTable.SetStyle(simpletable.StyleCompactLite)
			fmt.Println("\n\n" + assignmentsTable.String() + "\n\n")
		} else {
			assignments := models.Assignments(ctx, api, roles, periods, "assignment")
			assignmentsTable := views.AssignmentTable(assignments)
			assignmentsTable.SetStyle(simpletable.StyleCompactLite)
			fmt.Println("\n\n" + assignmentsTable.String() + "\n\n")

			if viper.GetBool("get-assignments-cmd-include-proposals") == true {
				propAssignments := models.Assignments(ctx, api, roles, periods, "proposal")
				propAssignmentsTable := views.AssignmentTable(propAssignments)
				propAssignmentsTable.SetStyle(simpletable.StyleCompactLite)
				fmt.Println("\n\n" + propAssignmentsTable.String() + "\n\n")
				return
			}
		}
	},
}

func init() {
	getCmd.AddCommand(getAssignmentCmd)
	getAssignmentCmd.Flags().BoolP("include-proposals", "i", false, "include proposals in the output")
	getAssignmentCmd.Flags().BoolP("failed-proposals", "f", false, "show only failed proposals")
}
