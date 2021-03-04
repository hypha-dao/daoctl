package cmd

import (
	"context"
	"fmt"

	"github.com/alexeyco/simpletable"
	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/daoctl/util"
	"github.com/hypha-dao/daoctl/views"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getRolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "retrieve roles",
	Long:  "retrieve all active roles For a json dump, append the argument --json.",
	RunE: func(cmd *cobra.Command, args []string) error {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()
		contract := eos.AN(viper.GetString("DAOContract"))

		gc, err := util.GetCache(ctx, api, contract)
		if err != nil {
			return fmt.Errorf("cannot get cache: %v", err)
		}

		roleDocs := gc.DocsByType["role"]
		roles := make([]models.Role, len(roleDocs))
		for idx, roleHash := range roleDocs {
			roleDoc, found := gc.Cache.Get(roleHash)
			if !found {
				return fmt.Errorf("document cache out of whack, try deleting the cache file: %v %v", roleHash, err)
			}

			roles[idx], err = models.NewRole(roleDoc.(docgraph.Document))
			if err != nil {
				return fmt.Errorf("cannot get convert document to role type: %v", err)
			}
		}

		rolesTable := views.RoleTable(roles)
		rolesTable.SetStyle(simpletable.StyleCompactLite)

		fmt.Println("\n" + rolesTable.String() + "\n\n")

		return nil
	},
}

func init() {
	getCmd.AddCommand(getRolesCmd)
}
