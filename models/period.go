package models

import (
	"context"

	eos "github.com/eoscanada/eos-go"
)

// Period represents a period of time aligning to a payroll period, typically a week
type Period struct {
	PeriodID  uint64             `json:"period_id"`
	StartTime eos.BlockTimestamp `json:"start_date"`
	EndTime   eos.BlockTimestamp `json:"end_date"`
	Phase     string             `json:"phase"`
}

// LoadPeriods loads the period data from the blockchain
func LoadPeriods(api *eos.API) []Period {
	var periods []Period
	var periodRequest eos.GetTableRowsRequest
	periodRequest.Code = "dao.hypha"
	periodRequest.Scope = "dao.hypha"
	periodRequest.Table = "periods"
	periodRequest.Limit = 1000
	periodRequest.JSON = true

	periodResponse, err := api.GetTableRows(context.Background(), periodRequest)
	if err != nil {
		panic(err)
	}
	periodResponse.JSONToStructs(&periods)
	return periods
}
