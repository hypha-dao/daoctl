package models

import (
	"context"

	eos "github.com/eoscanada/eos-go"
	"github.com/spf13/viper"
)

// RedemptionRequest is a type that represents a redemption request by a member
type RedemptionRequest struct {
	ID            uint64             `json:"redemption_id"`
	Requestor     eos.Name           `json:"requestor"`
	Requested     eos.Asset          `json:"amount_requested"`
	Paid          eos.Asset          `json:"amount_paid"`
	NotesRaw      []StringKV         `json:"notes"`
	RequestedDate eos.BlockTimestamp `json:"requested_date"`
	UpdatedDate   eos.BlockTimestamp `json:"updated_date"`
	NotesMap      *map[string]string
}

// LoadRequestByID returns a request for the provided redemption ID
func LoadRequestByID(ctx context.Context, api *eos.API, ID uint64) RedemptionRequest {
	var requests []RedemptionRequest
	// var memberAccounts []eos.Name
	var request eos.GetTableRowsRequest
	request.Code = viper.GetString("Treasury.Contract")
	request.Scope = viper.GetString("Treasury.Contract")
	request.Table = "redemptions"
	request.Limit = 1000 // TODO: support dynamic number of members
	request.JSON = true
	response, _ := api.GetTableRows(ctx, request)
	response.JSONToStructs(&requests)

	return requests[0]
}

// Requests returns a list of all redemption requests
func Requests(ctx context.Context, api *eos.API, all bool) []RedemptionRequest {
	var requests []RedemptionRequest
	// var memberAccounts []eos.Name
	var request eos.GetTableRowsRequest
	request.Code = viper.GetString("Treasury.Contract")
	request.Scope = viper.GetString("Treasury.Contract")
	request.Table = "redemptions"
	request.Limit = 1000 // TODO: support dynamic number of members
	request.JSON = true
	request.Index = "3"
	request.KeyType = "i64"
	request.Reverse = true
	response, _ := api.GetTableRows(ctx, request)
	response.JSONToStructs(&requests)

	for index, r := range requests {
		if !all && r.Paid.Amount >= r.Requested.Amount {
			return requests[0:index]
		}
		requests[index].NotesMap = ToMap(r.NotesRaw)
	}

	return requests
}
