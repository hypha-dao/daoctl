package cmd

import (
	"context"
	"fmt"
	"strconv"

	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/util"
	"github.com/hypha-dao/document-graph/docgraph"

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
	Use:   "fix",
	Short: "fix a document - admin only",

	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := context.Background()
		api := getAPI()
		contract := eos.AN(viper.GetString("DAOContract"))

		gc, err := util.GetCache(ctx, api, contract)
		if err != nil {
			return fmt.Errorf("cannot get cache: %v", err)
		}
		fmt.Println("Number of items in the cache: " + strconv.Itoa(gc.Cache.ItemCount()))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(fixCmd)
}
