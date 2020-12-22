package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/document-graph/docgraph"

	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type proposal struct {
	Proposer      eos.AccountName         `json:"proposer"`
	ProposalType  eos.Name                `json:"proposal_type"`
	ContentGroups []docgraph.ContentGroup `json:"content_groups"`
}

var proposeRoleCmd = &cobra.Command{
	Use:   "role -f [filename]",
	Short: "propose a role",

	Run: func(cmd *cobra.Command, args []string) {

		data, err := ioutil.ReadFile(viper.GetString("propose-cmd-file"))
		if err != nil {
			fmt.Println("Unable to read file: ", viper.GetString("propose-cmd-file"))
			return
		}

		var proposalDoc docgraph.Document
		err = json.Unmarshal([]byte(data), &proposalDoc)
		if err != nil {
			panic(err)
		}

		contract := toAccount(viper.GetString("DAOContract"), "contract")
		action := eos.ActN("propose")
		actions := eos.Action{
			Account: contract,
			Name:    action,
			Authorization: []eos.PermissionLevel{
				{Actor: eos.AN(viper.GetString("DAOUser")), Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(proposal{
				Proposer:      eos.AN(viper.GetString("DAOUser")),
				ProposalType:  eos.Name("role"),
				ContentGroups: proposalDoc.ContentGroups,
			})}

		pushEOSCActions(context.Background(), getAPI(), &actions)
	},
}

func init() {
	proposeCmd.AddCommand(proposeRoleCmd)
}
