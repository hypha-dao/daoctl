package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/daoctl/views"
	"github.com/spf13/cobra"
)

var getAssignmentsCmd = &cobra.Command{
	Use:   "assignments",
	Short: "retrieve assignments",
	Long:  "retrieve and print assignment tables",
	Run: func(cmd *cobra.Command, args []string) {
		// api := eos.New(viper.GetString("EosioEndpoint"))
		// ctx := context.Background()

		// periods := models.LoadPeriods(api, true, true)
		// roles := models.Roles(ctx, api, "role")
		// includeExpired := viper.GetBool("get-assignments-cmd-expired")

		// if viper.GetBool("global-csv") {
		// 	assignmentsTable := getAssignmentTable(ctx, api, roles, periods, "Active Assignment", "assignment", includeExpired)
		// 	csvData := models.TableToData(assignmentsTable)

		// 	file, err := os.Create(viper.GetString("global-output-file"))
		// 	if err != nil {
		// 		log.Fatalln("error writing csv:", err)
		// 	}

		// 	defer file.Close()

		// 	w := csv.NewWriter(file)
		// 	w.WriteAll(csvData) // calls Flush internally

		// 	if err := w.Error(); err != nil {
		// 		log.Fatalln("error writing csv:", err)
		// 	}
		// } else {
		// 	if viper.GetBool("global-active") == true {
		// 		printAssignmentTable(ctx, api, roles, periods, "Active Assignment", "assignment", includeExpired)
		// 	}

		// 	if viper.GetBool("global-include-proposals") == true {
		// 		printAssignmentTable(ctx, api, roles, periods, "Current Assignment Proposals", "proposal", includeExpired)
		// 	}

		// 	if viper.GetBool("global-failed-proposals") == true {
		// 		printAssignmentTable(ctx, api, roles, periods, "Failed Assignment Proposals", "failedprops", includeExpired)
		// 	}

		// 	if viper.GetBool("global-include-archive") == true {
		// 		printAssignmentTable(ctx, api, roles, periods, "Archive of Assignment Proposals", "proparchive", includeExpired)
		// 	}
		// }
	},
}

func getAssignmentTable(ctx context.Context, api *eos.API, roles []models.Role, periods []models.Period, title, scope string, includeExpired bool) *simpletable.Table {
	assignments, err := models.Assignments(ctx, api, roles, periods, scope, includeExpired)
	if err != nil {
		fmt.Println("Cannot get list of assignments: " + err.Error())
		os.Exit(-1)
	}
	return views.AssignmentTable(assignments)
}

func printAssignmentTable(ctx context.Context, api *eos.API, roles []models.Role, periods []models.Period, title, scope string, includeExpired bool) {
	fmt.Println("\n", title)
	assignmentsTable := getAssignmentTable(ctx, api, roles, periods, title, scope, includeExpired)
	assignmentsTable.SetStyle(simpletable.StyleCompactLite)
	fmt.Println("\n" + assignmentsTable.String() + "\n\n")
}

func init() {
	getCmd.AddCommand(getAssignmentsCmd)
	getAssignmentsCmd.Flags().BoolP("expired", "", false, "include expired assignments in the list")
}
