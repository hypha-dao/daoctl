package cmd

import (
	"context"
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/daoctl/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getRoleCmd = &cobra.Command{
	Use:   "role [role-hash]",
	Short: "retrieve role details",
	Long:  "retrieve the detailed about a role",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()
		contract := eos.AN(viper.GetString("DAOContract"))

		roleDoc, err := util.Get(ctx, api, contract, args[0])
		if err != nil {
			return fmt.Errorf("cannot find document with hash: %v %v", args[0], err)
		}

		role, err := models.NewRole(roleDoc)
		if err != nil {
			return fmt.Errorf("cannot convert document to role type: %v %v", args[0], err)
		}

		fmt.Println("\n\nRole: ", role.Title)
		fmt.Println()
		fmt.Println(role.String())
		fmt.Println()
		return nil
	},
}

func init() {
	getCmd.AddCommand(getRoleCmd)
}
