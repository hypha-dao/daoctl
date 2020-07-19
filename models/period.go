package models

import (
	"context"
	"fmt"
	"time"

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
func LoadPeriods(api *eos.API, includePast, includeFuture bool) []Period {

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

	var returnPeriods []Period
	currentPeriod, err := CurrentPeriod(&periods)
	if (includePast || includeFuture) && err != nil {
		panic(err)
	}

	for _, period := range periods {
		if includePast || includeFuture {
			if includePast && period.PeriodID <= uint64(currentPeriod) {
				returnPeriods = append(returnPeriods, period)
			} else if includeFuture && period.PeriodID >= uint64(currentPeriod) {
				returnPeriods = append(returnPeriods, period)
			}
		}
	}
	return returnPeriods
}

// CurrentPeriod provides the period ID for the current date and time
func CurrentPeriod(periods *[]Period) (int64, error) {
	now := time.Now()

	// assume that periods are in sorted
	for _, period := range *periods {
		if now.After(period.StartTime.Time) && now.Before(period.EndTime.Time) {
			return int64(period.PeriodID), nil
		}
	}
	return -1, fmt.Errorf("current time does not fall within a period")
}
