package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type paymentDocRecord struct {
	Recipient     string
	Amount        string
	Memo          string
	FromNodeTitle string
	PeriodStart   time.Time
	PaymentDate   time.Time
	EarliestDate  string
	docgraph.Document
}

var getPaymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "retrieve list of payments",
	// Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		api := eos.New(viper.GetString("EosioEndpoint"))
		ctx := context.Background()
		contract := eos.AN(viper.GetString("DAOContract"))

		if viper.GetBool("get-payments-cmd-documents") {
			paymentDocs, err := getDocumentsOfType(ctx, api, eos.Name("payment"))
			if err != nil {
				return fmt.Errorf("cannot get all documents: %v", err)
			}

			zap.S().Debugf("retrieved payment documents from chain: %v", len(paymentDocs))

			typesOfFromNodes := make(map[eos.Name]int)
			// var paymentRecords []paymentDocRecord
			paymentRecords := make([]paymentDocRecord, 0)
			for count, payment := range paymentDocs {
				if count > 50000 {
					break
				}
				pdr := paymentDocRecord{}
				pdr.Document = payment

				recipientFV, _ := payment.GetContentFromGroup("details", "recipient")
				pdr.Recipient = recipientFV.String()

				memoFV, _ := payment.GetContentFromGroup("details", "memo")
				pdr.Memo = memoFV.String()

				amountFV, _ := payment.GetContentFromGroup("details", "amount")
				pdr.Amount = amountFV.String()

				paymentDateFV, err := payment.GetContentFromGroup("details", "payment_date")
				if err == nil {
					paymentDateTimePoint, _ := paymentDateFV.TimePoint()
					pdr.PaymentDate = time.Unix(int64(paymentDateTimePoint)/1000000, 0).UTC()
				}

				edgesTo, err := docgraph.GetEdgesToDocumentWithEdge(ctx, api, contract, payment, eos.Name("payment"))
				if err != nil {
					return fmt.Errorf("cannot edges to payment with payment edge: %v", err)
				}
				zap.S().Debugf("loaded from edges for document: %v; edge count: %v", payment.Hash.String(), len(edgesTo))

				for _, edge := range edgesTo {
					docFrom, err := docgraph.LoadDocument(ctx, api, contract, edge.FromNode.String())
					if err != nil {
						return fmt.Errorf("cannot get document pointing to payment: %v", err)
					}

					docType, _ := docFrom.GetType()
					typesOfFromNodes[docType]++
					if docType == eos.Name("period") {
						period, err := models.NewSinglePeriod(ctx, api, contract, docFrom)
						if err != nil {
							return fmt.Errorf("unable to load period: %v", err)
						}
						pdr.PeriodStart = period.StartTime
					} else { //if docType == eos.Name("payout") {
						pdr.FromNodeTitle = docFrom.GetNodeLabel()
					}
				}
				paymentRecords = append(paymentRecords, pdr)
			}

			for docType, qty := range typesOfFromNodes {
				fmt.Println("type of from node: ", docType, ":	", strconv.Itoa(qty))
			}

			// Create a csv file
			f, err := os.Create("./doc-payments.csv")
			if err != nil {
				fmt.Println(err)
			}
			defer f.Close()
			// Write Unmarshaled json data to CSV file
			// now := time.Now()
			// twoYears, _ := time.ParseDuration("700d")
			floor := time.Now().AddDate(-2, 0, 0)
			fmt.Println("floor: ", floor)

			w := csv.NewWriter(f)
			header := []string{"payment_label", "from_node_label", "recipient", "recognition_date", "year", "month", "day", "amount", "token", "memo", "payment_date", "period_start", "created_date", "asset", "hash"}
			w.Write(header)
			for _, p := range paymentRecords {
				var record []string
				record = append(record, p.GetNodeLabel())
				record = append(record, p.FromNodeTitle)
				record = append(record, string(p.Recipient))

				earliestDate := p.CreatedDate.Time
				if p.PeriodStart.After(floor) && p.PeriodStart.Before(earliestDate) {
					earliestDate = p.PeriodStart
				}
				if p.PaymentDate.After(floor) && p.PaymentDate.Before(earliestDate) {
					earliestDate = p.PaymentDate
				}

				record = append(record, earliestDate.Format("2006 Jan 02"))
				record = append(record, fmt.Sprint(earliestDate.Year()))
				record = append(record, fmt.Sprint(earliestDate.Month()))
				record = append(record, fmt.Sprint(earliestDate.Day()))

				record = append(record, strings.Fields(p.Amount)[0])
				record = append(record, strings.Fields(p.Amount)[1])
				record = append(record, p.Memo)

				record = append(record, p.PaymentDate.Format("2006 Jan 02"))
				record = append(record, p.PeriodStart.Format("2006 Jan 02"))
				record = append(record, p.CreatedDate.Format("2006 Jan 02"))

				record = append(record, p.Amount)
				record = append(record, p.Hash.String())

				w.Write(record)
			}
			w.Flush()
		} else {
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
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(getPaymentsCmd)
	getPaymentsCmd.Flags().BoolP("documents", "d", true, "use the documents table rather than the now deprecated payments table")
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

func getDocumentsOfType(ctx context.Context, api *eos.API, docType eos.Name) ([]docgraph.Document, error) {

	docs, err := docgraph.GetAllDocuments(ctx, api, eos.AN(viper.GetString("DAOContract")))
	if err != nil {
		return []docgraph.Document{}, fmt.Errorf("cannot get all documents: %v", err)
	}

	var filteredDocs []docgraph.Document
	for _, doc := range docs {

		typeFV, err := doc.GetContent("type")
		if err == nil &&
			typeFV.Impl.(eos.Name) == docType {
			filteredDocs = append(filteredDocs, doc)
		}
	}
	return filteredDocs, nil
}
