package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func configFileAction(ctx context.Context, configFilename string) *eos.Action {
	data, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Panicln("Unable to read file: ", viper.GetString("create-cmd-file"))
		return nil
	}

	contract := toAccount(viper.GetString("DAOContract"), "contract")
	actionName := toActionName("setconfig", "action")

	var dump map[string]interface{}
	err = json.Unmarshal(data, &dump)
	if err != nil {
		log.Panicln("Unable to unmarshal config json: ", err)
		return nil
	}

	api := getAPI()
	actionBinary, err := api.ABIJSONToBin(ctx, contract, eos.Name(actionName), dump)
	errorCheck("unable to retrieve action binary from JSON via API", err)

	action := eos.Action{
		Account: contract,
		Name:    actionName,
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AN(viper.GetString("DAOContract")), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionDataFromHexData([]byte(actionBinary)),
	}
	return &action
}

var setConfigCmd = &cobra.Command{
	Use:   "config [-f filename] [-k config-key] [-v config-value]",
	Short: "OLD - DO NOT USE - set the configuration based on a file OR a configuration key and value",
	Long:  "OLD - DO NOT USE - set the configuration based on a file OR a configuration key and value",

	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		if len(viper.GetString("set-config-cmd-file")) > 0 {
			action := configFileAction(ctx, viper.GetString("set-config-cmd-file"))
			pushEOSCActions(context.Background(), getAPI(), action)
		} else {
			zlog.Panic("only setting via configuration file is currently supported")
		}
	},
}

func init() {
	setCmd.AddCommand(setConfigCmd)
	setConfigCmd.Flags().StringP("config-key", "k", "", "key of the key-value pair to set in the DAO's configuration")
	setConfigCmd.Flags().StringP("config-value", "v", "", "value of the key-value pair to set in the DAO's configuration")
	setConfigCmd.Flags().StringP("file", "f", "", "filename containing JSON configuration")

}
