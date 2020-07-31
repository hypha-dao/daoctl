package models

import (
	"context"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/viper"
)

// DAOPayment represents a reimbursement on a redemption request
type DAOPayment struct {
	ID           uint64             `json:"payment_id"`
	Recipient    eos.Name           `json:"recipient"`
	Amount       eos.Asset          `json:"amount"`
	PaymentDate  eos.BlockTimestamp `json:"payment_date"`
	AssignmentID uint64             `json:"assignment_id"`
	Memo         string             `json:"memo"`
}

// DAOPayments returns a list of all redemption payments
func DAOPayments(ctx context.Context, api *eos.API) []Payment {
	var payments []Payment
	var request eos.GetTableRowsRequest
	request.Code = viper.GetString("DAOContract")
	request.Scope = viper.GetString("DAOContract")
	request.Table = "payments"
	request.Limit = 1000 // TODO: support dynamic number of results
	request.JSON = true
	response, _ := api.GetTableRows(ctx, request)
	response.JSONToStructs(&payments)

	return payments
}
