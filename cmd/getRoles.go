package cmd

import (
	"context"
	"fmt"

	"github.com/alexeyco/simpletable"
	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/daoctl/views"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getRolesCmd = &cobra.Command{
	Use:   "roles [account name]",
	Short: "retrieve roles",
	Long:  "retrieve all active roles For a json dump, append the argument --json.",
	// Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		periods := models.LoadPeriods(api)
		roles := models.Roles(ctx, api, periods)
		rolesTable := views.RoleTable(roles)
		rolesTable.SetStyle(simpletable.StyleCompactLite)

		fmt.Println("\n\n" + rolesTable.String() + "\n\n")

		if viper.GetBool("get-roles-cmd-include-proposals") == true {
			propRoles := models.ProposedRoles(ctx, api, periods)
			propRolesTable := views.RoleTable(propRoles)
			propRolesTable.SetStyle(simpletable.StyleCompactLite)
			fmt.Println("\n\n" + propRolesTable.String() + "\n\n")
			return
		}
	},
}

func init() {
	getCmd.AddCommand(getRoleCmd)
	getRoleCmd.Flags().BoolP("include-proposals", "i", false, "include proposals in the output")
}
