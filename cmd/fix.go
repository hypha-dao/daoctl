package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/document-graph/docgraph"

	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type period struct {
	Proposer      eos.AccountName         `json:"proposer"`
	ProposalType  eos.Name                `json:"proposal_type"`
	ContentGroups []docgraph.ContentGroup `json:"content_groups"`
}

// Period represents a period of time aligning to a payroll period, typically a week
type Period struct {
	PeriodID  uint64             `json:"period_id"`
	StartTime eos.BlockTimestamp `json:"start_date"`
	EndTime   eos.BlockTimestamp `json:"end_date"`
	Phase     string             `json:"phase"`
}

var fixCmd = &cobra.Command{
	Use:   "fix [doc-hash]",
	Short: "fix a document - admin only",

	Run: func(cmd *cobra.Command, args []string) {

		ctx := context.Background()
		api := getAPI()
		contract := eos.AN(viper.GetString("DAOContract"))

		data, err := ioutil.ReadFile("dao-backup-migration/dao.hypha_periods.json")
		if err != nil {
			fmt.Println("Unable to read file: dao-backup-migration/dao.hypha_periods.json")
			return
		}

		var periods []models.Period
		err = json.Unmarshal([]byte(data), &periods)
		if err != nil {
			panic(err)
		}

		documents, err := docgraph.GetAllDocuments(ctx, api, contract)
		if err != nil {
			panic(err)
		}

		for _, document := documents {

			typeFV, err := document.GetContent("type")
			docType := typeFV.Impl.(eos.Name)

			if docType == eos.Name("assignment") {
				claimedEdges, err := docgraph.GetEdgesFromDocumentWithEdge(ctx, api, contract, document, eos.Name("claimed") )
				if err != nil {
					panic(err)
				}

				periodCountFV, err := document.GetContent("type")
				periodCount := typeFV.Impl.(uint64)
				claimCount := len(claimedEdges)
				

			}

			

			
			periodCount
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
