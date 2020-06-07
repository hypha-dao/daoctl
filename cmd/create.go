package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	eos "github.com/eoscanada/eos-go"

	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:   "create -f [filename]",
	Short: "create an object based on the JSON file",
	Long:  "create an object based on the JSON file",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		data, err := ioutil.ReadFile(viper.GetString("create-cmd-file"))
		if err != nil {
			fmt.Println("Unable to read file: ", viper.GetString("create-cmd-file"))
			return
		}

		contract := toAccount(viper.GetString("DAOContract"), "contract")
		action := toActionName("create", "action")

		var dump map[string]interface{}
		err = json.Unmarshal(data, &dump)
		if err != nil {
			fmt.Println("Unable to unmarshal json: ", err)
			return
		}

		api := getAPI()
		actionBinary, err := api.ABIJSONToBin(ctx, contract, eos.Name(action), dump)
		errorCheck("unable to retrieve action binary from JSON via API", err)

		actions := []*eos.Action{
			&eos.Action{
				Account: contract,
				Name:    action,
				Authorization: []eos.PermissionLevel{
					{Actor: eos.AN(viper.GetString("DAOUser")), Permission: eos.PN("active")},
				},
				ActionData: eos.NewActionDataFromHexData([]byte(actionBinary)),
			}}

		pushEOSCActions(context.Background(), getAPI(), actions[0])
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("file", "f", "", "filename of object's JSON file")
}
