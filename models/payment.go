package models

import (
	"context"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/viper"
)

// Attestation that a particular payment is valid and true
type Attestation struct {
	Key   eos.Name           `json:"key"`
	Value eos.BlockTimestamp `json:"value"`
}

// Payment represents a reimbursement on a redemption request
type Payment struct {
	ID            uint64             `json:"payment_id"`
	RequestID     uint64             `json:"redemption_id"`
	Creator       eos.Name           `json:"creator"`
	Amount        eos.Asset          `json:"amount_paid"`
	CreatedDate   eos.BlockTimestamp `json:"created_date"`
	ConfirmedDate eos.BlockTimestamp `json:"confirmed_date"`
	Attestations  []Attestation      `json:"attestations"`
	NotesRaw      []StringKV         `json:"notes"`
	NotesMap      *map[string]string
	Request       RedemptionRequest
	//Attestations  map[eos.Name]eos.BlockTimestamp `json:"attestations"`
}

// Payments returns a list of all redemption payments
func Payments(ctx context.Context, api *eos.API) []Payment {
	var payments []Payment
	var request eos.GetTableRowsRequest
	request.Code = viper.GetString("Treasury.Contract")
	request.Scope = viper.GetString("Treasury.Contract")
	request.Table = "payments"
	request.Limit = 1000 // TODO: support dynamic number of results
	request.JSON = true
	response, _ := api.GetTableRows(ctx, request)
	response.JSONToStructs(&payments)

	for index, p := range payments {
		payments[index].NotesMap = ToMap(p.NotesRaw)
	}

	return payments
}
