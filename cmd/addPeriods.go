package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/document-graph/docgraph"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type addPeriod struct {
	Predecessor eos.Checksum256 `json:"predecessor"`
	StartTime   eos.TimePoint   `json:"start_time"`
	Label       string          `json:"label"`
}

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

// moon phase
type MoonPhase struct {
	Timestamp uint64 `json:"timestamp"`
	PhaseName string `json:"phase_name"`
	PhaseTime eos.BlockTimestamp
}

func execTrx(ctx context.Context, api *eos.API, actions []*eos.Action) (string, error) {
	txOpts := &eos.TxOptions{}
	if err := txOpts.FillFromChain(ctx, api); err != nil {
		return string(""), fmt.Errorf("error filling tx opts: %s", err)
	}

	tx := eos.NewTransaction(actions, txOpts)

	_, packedTx, err := api.SignTransaction(ctx, tx, txOpts.ChainID, eos.CompressionNone)
	if err != nil {
		return string(""), fmt.Errorf("error signing transaction: %s", err)
	}

	response, err := api.PushTransaction(ctx, packedTx)
	if err != nil {
		return string(""), fmt.Errorf("error pushing transaction: %s", err)
	}
	trxID := hex.EncodeToString(response.Processed.ID)
	return trxID, nil
}

var addPeriodsCmd = &cobra.Command{
	Use:   "addperiods",
	Short: "addperiods - admin only",

	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := context.Background()
		api := getAPI()
		contract := eos.AN(viper.GetString("DAOContract"))
		period, err := getLastPeriod(ctx, api, contract)

		keyBag := &eos.KeyBag{}
		keyBag.ImportPrivateKey(context.Background(), "5KFCMj1ewfRYPhP7kCp9S6FpHKheRBS9sZLLNTqZu3WHbQiVG9s")
		api.SetSigner(keyBag)

		if err != nil {
			return fmt.Errorf("cannot get latest period: %v", err)
		}

		predecessor := period.Document.Hash
		fmt.Println("Initial predecessor			: " + predecessor.String())

		timestamp := viper.GetInt("addperiods-cmd-start-time")
		periodCount := uint32(viper.GetInt("addperiods-cmd-period-count"))

		fmt.Println("start time: " + strconv.Itoa(timestamp))
		fmt.Println("period count: " + strconv.Itoa(int(periodCount)))

		var phases []MoonPhase
		var request eos.GetTableRowsRequest
		request.Code = "cycle.seeds"
		request.Scope = "cycle.seeds"
		request.Table = "moonphases"
		request.LowerBound = strconv.Itoa(timestamp)
		request.Limit = periodCount
		request.JSON = true
		response, err := api.GetTableRows(ctx, request)
		if err != nil {
			return fmt.Errorf("get table rows %v", err)
		}

		err = response.JSONToStructs(&phases)
		if err != nil {
			return fmt.Errorf("json to structs %v", err)
		}

		if len(phases) == 0 {
			return fmt.Errorf("phases not found %v", err)
		}

		for _, phase := range phases {

			startTime := time.Unix(int64(phase.Timestamp), 0).UTC()
			nodeLabel := "Starting " + startTime.Format("2006 Jan 02")
			startTimePoint := eos.TimePoint(phase.Timestamp * 1000000)

			fmt.Println("Add period: ")
			fmt.Println("	Predecessor			: " + predecessor.String())
			fmt.Println(" 	Start time (int)		: " + strconv.Itoa(int(phase.Timestamp)))
			fmt.Println("	Start time (tp) 		: " + startTimePoint.String())
			fmt.Println(" 	Start time (read)		: " + startTime.Format("2006 Jan 02 15:04:05"))
			fmt.Println("	Label				: " + phase.PhaseName)

			addPeriodAction := []*eos.Action{{
				Account: contract,
				Name:    eos.ActN("addperiod"),
				Authorization: []eos.PermissionLevel{
					{Actor: contract, Permission: eos.PN("active")},
				},
				ActionData: eos.NewActionData(addPeriod{
					Predecessor: predecessor,
					StartTime:   startTimePoint,
					Label:       phase.PhaseName,
				}),
			}}

			trxId, err := execTrx(ctx, api, addPeriodAction)
			if err != nil {
				return fmt.Errorf("cannot execute transaction: %v", err)
			}
			fmt.Println(" Transaction ID: ", trxId)

			fmt.Println("	Node Label			: " + nodeLabel)
			fmt.Println("	Readable Start Date		: " + startTime.Format("2006 Jan 02"))
			fmt.Println("	Readable Start Time		: " + startTime.Format("15:04:05 UTC"))
			fmt.Println()

			updateLastPeriod(ctx, api, contract, "system", "node_label", nodeLabel)

			updateLastPeriod(ctx, api, contract, "system", "readable_start_date", startTime.Format("2006 Jan 02"))

			updateLastPeriod(ctx, api, contract, "system", "readable_start_time", startTime.Format("15:04:05 UTC"))

			lastPeriod, err := docgraph.GetLastDocumentOfEdge(ctx, api, contract, eos.Name("next"))
			if err != nil {
				return fmt.Errorf("cannot get last created period: %v", err)
			}

			predecessor = lastPeriod.Hash
		}

		return nil
	},
}

func init() {

	RootCmd.AddCommand(addPeriodsCmd)

	addPeriodsCmd.Flags().IntP("start-time", "s", 0, "the start time (moment) of the period that matches the timestamp column in moonphases table")
	addPeriodsCmd.Flags().IntP("period-count", "p", 0, "the number of periods to add from the moonphases table")
}

func getLastPeriod(ctx context.Context, api *eos.API, contract eos.AccountName) (models.Period, error) {
	// rootDocument, err := docgraph.LoadDocument(ctx, api, contract, viper.GetString("RootNode"))
	// if err != nil {
	// 	return models.Period{}, fmt.Errorf("cannot load root document: %v", err)
	// }

	// startEdges, err := docgraph.GetEdgesFromDocumentWithEdge(ctx, api, contract, rootDocument, eos.Name("start"))
	// if err != nil {
	// 	return models.Period{}, fmt.Errorf("error while retrieving start edge: %v", err)
	// }
	// fmt.Println(startEdges)
	// if len(startEdges) == 0 {
	// 	return models.Period{}, fmt.Errorf("no start edge from the root node exists: %v", err)
	// }

	startPeriodDoc, err := docgraph.LoadDocument(ctx, api, contract, viper.GetString("CalendarStart")) //startEdges[0].ToNode.String())
	if err != nil {
		return models.Period{}, fmt.Errorf("error loading the start period document: %v", err)
	}

	period, err := models.NewPeriod(ctx, api, contract, startPeriodDoc)
	if err != nil {
		return models.Period{}, fmt.Errorf("cannot convert document to period type: %v", err)
	}

	for period.Next != nil {
		period = *period.Next
	}

	return period, nil
}

type updateDoc struct {
	Hash    eos.Checksum256    `json:"hash"`
	Updater eos.AccountName    `json:"updater"`
	Group   string             `json:"group"`
	Key     string             `json:"key"`
	Value   docgraph.FlexValue `json:"value"`
}

func updateLastPeriod(ctx context.Context, api *eos.API, contract eos.AccountName, group, key, value string) error {

	lastPeriod, err := docgraph.GetLastDocumentOfEdge(ctx, api, contract, eos.Name("next"))
	if err != nil {
		return fmt.Errorf("cannot get last created period: %v", err)
	}

	fmt.Println("last period: " + lastPeriod.Hash.String())

	actions := []*eos.Action{{
		Account: contract,
		Name:    eos.ActN("updatedoc"),
		Authorization: []eos.PermissionLevel{
			{Actor: contract, Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(updateDoc{
			Hash:    lastPeriod.Hash,
			Updater: eos.AN("dao.hypha"),
			Group:   group,
			Key:     key,
			Value: docgraph.FlexValue{
				BaseVariant: eos.BaseVariant{
					TypeID: docgraph.GetVariants().TypeID("string"),
					Impl:   value,
				},
			}}),
	}}

	execTrx(ctx, api, actions)
	return nil
}
