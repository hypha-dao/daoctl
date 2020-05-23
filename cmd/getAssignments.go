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
	Long:  "retrieve and print assignment tables",
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		periods := models.LoadPeriods(api)
		roles := models.Roles(ctx, api, periods, "role")

		if viper.GetBool("global-active") == true {
			printAssignmentTable(ctx, api, roles, periods, "Active Assignment", "assignment")
		}

		if viper.GetBool("global-include-proposals") == true {
			printAssignmentTable(ctx, api, roles, periods, "Current Assignment Proposals", "proposal")
		}

		if viper.GetBool("global-failed-proposals") == true {
			printAssignmentTable(ctx, api, roles, periods, "Failed Assignment Proposals", "failedprops")
		}

		if viper.GetBool("global-include-archive") == true {
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
}
