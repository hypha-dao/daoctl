package models

import (
	"context"
	"errors"
	"fmt"
	"strconv"

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

// LoadPaymentByID returns a request for the provided redemption ID
func LoadPaymentByID(ctx context.Context, api *eos.API, ID uint64) (Payment, error) {
	var payments []Payment
	var err error
	// var memberAccounts []eos.Name
	var request eos.GetTableRowsRequest
	request.Code = viper.GetString("Treasury.Contract")
	request.Scope = viper.GetString("Treasury.Contract")
	request.Table = "payments"
	request.Limit = 1
	request.LowerBound = strconv.Itoa(int(ID))
	request.UpperBound = strconv.Itoa(int(ID))
	request.JSON = true
	response, _ := api.GetTableRows(ctx, request)
	response.JSONToStructs(&payments)

	if len(payments) >= 1 {
		payments[0].NotesMap = ToMap(payments[0].NotesRaw)
		payments[0].Request, err = LoadRequestByID(ctx, api, payments[0].RequestID)
		if err != nil {
			fmt.Println("Warning: this payment's corresponding request is not found - this should not happen")
		}
		return payments[0], nil
	}

	return Payment{}, errors.New("Payment not found")
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
