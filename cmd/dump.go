package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tidwall/pretty"
)

var dumpCmd = &cobra.Command{
	Use:   "dump [hash]",
	Short: "raw dump of the documents json",
	Long:  "raw dump of the documents json",
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()
		contract := eos.AN(viper.GetString("DAOContract"))

		document, err := util.Get(ctx, api, contract, args[0])
		if err != nil {
			return fmt.Errorf("cannot find document with hash: %v %v", args[0], err)
		}

		docJson, err := json.Marshal(document)
		if err != nil {
			return fmt.Errorf("cannot marshall document to JSON: %v %v", args[0], err)
		}

		fmt.Println(string(pretty.Color(pretty.Pretty(docJson), nil)))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(dumpCmd)
}
