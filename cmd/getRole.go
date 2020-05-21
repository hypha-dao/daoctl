package cmd

import (
  "context"
  "fmt"
  "strconv"

  eos "github.com/eoscanada/eos-go"
  "github.com/hypha-dao/daoctl/models"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

var getRoleCmd = &cobra.Command{
	Use:   "role [role id]",
	Short: "retrieve role details",
	Long:  "retrieve the detailed about a role",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()
		//ac := accounting.NewAccounting("", 0, ",", ".", "%s %v", "%s (%v)", "%s --") // TODO: make this configurable

		roleID, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
		  fmt.Println("Parse error: Role id must be a positive integer (uint64)")
		  return;
    }
		periods := models.LoadPeriods(api)
		role := models.NewRoleByID(ctx, api, periods, roleID)

		fmt.Println("\n\nRole: ", role.Title, "\n")
		fmt.Println(role.String())
		fmt.Println()
	},
}

func init() {
	getCmd.AddCommand(getRoleCmd)
}
