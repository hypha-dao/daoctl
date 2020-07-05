package models

import (
	"context"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/viper"
)

// Payment represents a reimbursement on a redemption request
type Payment struct {
	ID            uint64                          `json:"payment_id"`
	Creator       eos.Name                        `json:"creator"`
	Amount        eos.Asset                       `json:"amount_paid"`
	CreatedDate   eos.BlockTimestamp              `json:"created_date"`
	ConfirmedDate eos.BlockTimestamp              `json:"confirmed_date"`
	Attestations  map[eos.Name]eos.BlockTimestamp `json:"attestations"`
	Notes         map[string]string               `json:"notes"`
	Request       RedemptionRequest
}

// Payments returns a list of all redemption payments
func Payments(ctx context.Context, api *eos.API) []Payment {
	var payments []Payment
	var request eos.GetTableRowsRequest
	request.Code = viper.GetString("TreasuryContract")
	request.Scope = viper.GetString("TreasuryContract")
	request.Table = "payments"
	request.Limit = 1000 // TODO: support dynamic number of results
	request.JSON = true
	response, _ := api.GetTableRows(ctx, request)
	response.JSONToStructs(&payments)

	return payments
}
