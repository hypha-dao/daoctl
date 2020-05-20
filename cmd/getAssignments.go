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
	Use:   "assignments",
	Short: "retrieve assignments",
	Long:  "retrieve all active assignments",
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		periods := models.LoadPeriods(api)
		roles := models.Roles(ctx, api, periods)

		if viper.GetBool("get-assignments-cmd-active") == true {
			printAssignmentTable(ctx, api, roles, periods, "Active Assignment", "assignment")
		}

		if viper.GetBool("get-assignments-cmd-include-proposals") == true {
			printAssignmentTable(ctx, api, roles, periods, "Current Assignment Proposals", "proposal")
		}

		if viper.GetBool("get-assignments-cmd-failed-proposals") == true {
			printAssignmentTable(ctx, api, roles, periods, "Failed Assignment Proposals", "failedprops")
		}

		if viper.GetBool("get-assignments-cmd-include-archive") == true {
			printAssignmentTable(ctx, api, roles, periods, "Archive of Assignment Proposals", "proparchive")
		}
	},
}

func printAssignmentTable(ctx context.Context, api *eos.API, roles []models.Role, periods []models.Period, title, scope string) {
	fmt.Println("\n", title)
	assignments := models.Assignments(ctx, api, roles, periods, scope)
	assignmentsTable := views.AssignmentTable(assignments)
	assignmentsTable.SetStyle(simpletable.StyleCompactLite)
	fmt.Println("\n" + assignmentsTable.String() + "\n\n")
}

func init() {
	getCmd.AddCommand(getAssignmentCmd)
	getAssignmentCmd.Flags().BoolP("include-proposals", "i", false, "include a table with proposals in the output")
	getAssignmentCmd.Flags().BoolP("failed-proposals", "f", false, "include a table with failed proposals")
	getAssignmentCmd.Flags().BoolP("active", "a", true, "show active assignments")
	getAssignmentCmd.Flags().BoolP("include-archive", "o", true, "include a table with the archive of assignment proposals")

}
