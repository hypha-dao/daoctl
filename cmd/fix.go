package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	eos "github.com/eoscanada/eos-go"
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

type Hash struct {
	hash eos.Checksum256 `json:"hash"`
}

func ExecWithRetry(ctx context.Context, api *eos.API, actions []*eos.Action) (string, error) {
	trxId, err := ExecTrx(ctx, api, actions)

	if err != nil {
		if !strings.Contains(err.Error(), "deadline exceeded") {
			return string(""), err
		} else {
			attempts := 1
			for attempts < 3 {
				trxId, err = ExecTrx(ctx, api, actions)
				if err == nil {
					return trxId, nil
				}
				attempts++
			}
		}
		return string(""), err
	}
	return trxId, nil
}

// ExecTrx executes a list of actions
func ExecTrx(ctx context.Context, api *eos.API, actions []*eos.Action) (string, error) {
	txOpts := &eos.TxOptions{}
	if err := txOpts.FillFromChain(ctx, api); err != nil {
		return "error", fmt.Errorf("error filling tx opts: %s", err)
	}

	tx := eos.NewTransaction(actions, txOpts)

	_, packedTx, err := api.SignTransaction(ctx, tx, txOpts.ChainID, eos.CompressionNone)
	if err != nil {
		return "error", fmt.Errorf("error signing transaction: %s", err)
	}

	response, err := api.PushTransaction(ctx, packedTx)
	if err != nil {
		return "error", fmt.Errorf("error pushing transaction: %s", err)
	}
	trxID := hex.EncodeToString(response.Processed.ID)
	return trxID, nil
}

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "fix a document - admin only",

	RunE: func(cmd *cobra.Command, args []string) error {

		ctx := context.Background()
		api := getAPI()
		contract := eos.AN(viper.GetString("DAOContract"))

		keyBag := &eos.KeyBag{}

		// this is the testnet key, in prod, I will use one of my low-security accounts
		keyBag.ImportPrivateKey(ctx, "5HwnoWBuuRmNdcqwBzd1LABFRKnTk2RY2kUMYKkZfF8tKodubtK")
		api.SetSigner(keyBag)

		docs, err := docgraph.GetAllDocuments(ctx, api, contract)
		if err != nil {
			return fmt.Errorf("cannot get all documents: %v", err)
		}

		for _, doc := range docs {

			typeFV, err := doc.GetContent("type")
			if err == nil && typeFV.Impl.(eos.Name) == eos.Name("assignment") {

				fmt.Println("Run fix on ", doc.Hash.String())
				actions := []*eos.Action{{
					Account: contract,
					Name:    eos.ActN("fix"),
					Authorization: []eos.PermissionLevel{
						{Actor: eos.AN(viper.GetString("DAOUser")), Permission: eos.PN("active")},
					},
					ActionData: eos.NewActionData(Hash{
						doc.Hash,
					})},
				}

				ExecWithRetry(ctx, api, actions)
			}

		}

		fmt.Println("Number of assignments fixed: ", len(docs))
		return nil
	},
}

func init() {
	RootCmd.AddCommand(fixCmd)
}
