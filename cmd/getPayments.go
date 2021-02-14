package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getPaymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "retrieve list of payments",
	// Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()

		payments, err := getAllPayments(ctx, api, eos.AN(viper.GetString("DAOContract")))
		if err != nil {
			panic(fmt.Errorf("cannot get all documents: %v", err))
		}

		fmt.Println("Number of payments: " + strconv.Itoa(len(payments)))

		// Create a csv file
		f, err := os.Create("./all-payments.csv")
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		// Write Unmarshaled json data to CSV file
		w := csv.NewWriter(f)
		for _, p := range payments {
			var record []string
			record = append(record, strconv.Itoa(int(p.ID)))
			record = append(record, p.PaymentDate.Time.Format("2006 Jan 02"))
			record = append(record, strconv.Itoa(int(p.PeriodID)))
			record = append(record, strconv.Itoa(int(p.AssignmentID)))
			record = append(record, string(p.Recipient))
			record = append(record, p.Amount.String())
			record = append(record, p.Memo)
			w.Write(record)
		}
		w.Flush()
	},
}

func init() {
	getCmd.AddCommand(getPaymentsCmd)
}

type payment struct {
	ID           uint64             `json:"payment_id"`
	PaymentDate  eos.BlockTimestamp `json:"payment_date"`
	PeriodID     eos.Uint64         `json:"period_id"`
	AssignmentID eos.Uint64         `json:"assignment_id"`
	Recipient    eos.Name           `json:"recipient"`
	Amount       eos.Asset          `json:"amount"`
	Memo         string             `json:"memo"`
}

func getRange(ctx context.Context, api *eos.API, contract eos.AccountName, id, count int) ([]payment, bool, error) {
	var documents []payment
	var request eos.GetTableRowsRequest
	if id > 0 {
		request.LowerBound = strconv.Itoa(id)
	}
	request.Code = string(contract)
	request.Scope = string(contract)
	request.Table = "payments"
	request.Limit = uint32(count)
	request.JSON = true
	response, err := api.GetTableRows(ctx, request)
	if err != nil {
		return []payment{}, false, fmt.Errorf("get table rows %v", err)
	}

	err = response.JSONToStructs(&documents)
	if err != nil {
		return []payment{}, false, fmt.Errorf("json to structs %v", err)
	}
	return documents, response.More, nil
}

func getAllPayments(ctx context.Context, api *eos.API, contract eos.AccountName) ([]payment, error) {

	var allPayments []payment

	batchSize := 150

	batch, more, err := getRange(ctx, api, contract, 0, batchSize)
	if err != nil {
		return []payment{}, fmt.Errorf("json to structs %v", err)
	}
	allPayments = append(allPayments, batch...)

	for more {
		batch, more, err = getRange(ctx, api, contract, int(batch[len(batch)-1].ID), batchSize)
		if err != nil {
			return []payment{}, fmt.Errorf("json to structs %v", err)
		}
		allPayments = append(allPayments, batch...)
	}

	return allPayments, nil
}
