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

var editCmd = &cobra.Command{
	Use:   "edit -f [filename]",
	Short: "propose an edit to an object",
	Long:  "propose an edit to an object based on JSON or a JSON file",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		var data []byte
		var err error
		if len(viper.GetString("edit-cmd-file")) > 0 {
			data, err = ioutil.ReadFile(viper.GetString("edit-cmd-file"))
			if err != nil {
				fmt.Println("Unable to read file: ", viper.GetString("edit-cmd-file"))
				return
			}
		} else if len(viper.GetString("edit-cmd-json")) > 0 {
			data = []byte(viper.GetString("edit-cmd-json"))
		} else {
			fmt.Println("Unable to read either --json or --file parameter.")
			return
		}

		contract := toAccount(viper.GetString("DAOContract"), "contract")
		action := toActionName("edit", "action")

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
			{
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
	RootCmd.AddCommand(editCmd)
	editCmd.Flags().StringP("json", "j", "", "JSON content with edit contents")
	editCmd.Flags().StringP("file", "f", "", "file containing JSON to represent the edit")
}
