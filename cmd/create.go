package cmd

import (
  "encoding/json"
  "fmt"
  "context"
  eos "github.com/eoscanada/eos-go"

  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "io/ioutil"
)

type createPayload struct {
	Data struct {
		Scope string `json:"scope"`
		Names []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"names"`
		Strings []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"strings"`
		Assets []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"assets"`
		TimePoints []interface{} `json:"time_points"`
		Ints       []struct {
			Key   string `json:"key"`
			Value uint64    `json:"value"`
		} `json:"ints"`
		Floats []interface{} `json:"floats"`
		Trxs   []interface{} `json:"trxs"`
	} `json:"data"`
}

func newCreateAction(payload createPayload) *eos.Action {
	return &eos.Action{
		Account: eos.AN(viper.GetString("DAOContract")),
		Name:    eos.ActN("create"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AN(viper.GetString("DAOUser")), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(payload.Data),
	}
}

var createCmd = &cobra.Command{
	Use:   "create [file name]",
	Short: "create an object based on the JSON file",
	Long:  "create an object based on the JSON file",
	// Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		//api := eos.New(viper.GetString("EosioEndpoint"))
		//ctx := context.Background()

		data, err := ioutil.ReadFile(viper.GetString("create-cmd-file"))
		if err != nil {
		  fmt.Println("Unable to read file: ", viper.GetString("create-cmd-file"))
		  return
		}

		var daoObject createPayload
		err = json.Unmarshal(data, &daoObject)
		if err != nil {
			fmt.Println("Unable to unmarshal json: ", err)
			return
		}
		backToJSON, _ := json.MarshalIndent(daoObject, "", "  ")
		fmt.Println(string(backToJSON))

		action := newCreateAction(daoObject)
		pushEOSCActions(context.Background(), getAPI(), action)
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("file", "f", "", "filename of object's JSON file")
}
