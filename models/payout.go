package models

import (
	"context"

	eos "github.com/eoscanada/eos-go"
)

// Payout represents a person assigned to a role for a specific period of time
type Payout struct {
	ID              uint64
	Approved        bool
	Receiver        eos.Name
	BallotName      eos.Name
	Title           string
	Description     string
	Husd            eos.Asset
	Hypha           eos.Asset
	Hvoice          eos.Asset
	SeedsEscrow     eos.Asset
	SeedsLiquid     eos.Asset
	DeferredPay     float64
	InstantHusdPerc float64
	StartPeriod     Period
	EndPeriod       Period
	CreatedDate     eos.BlockTimestamp
}

// NewPayout converts a generic DAO Object to a typed Payout
func NewPayout(daoObj DAOObject, periods []Period) Payout {
	var a Payout
	a.ID = daoObj.ID
	a.Receiver = daoObj.Names["recipient"]
	a.Title = daoObj.Strings["title"]
	a.BallotName = daoObj.Names["ballot_id"]
	a.Husd = daoObj.Assets["husd_amount"]

	if daoObj.Assets["husd_amount"].Amount == 0 {
		a.Husd, _ = eos.NewAssetFromString("0.00 HUSD")
	} else {
		a.Husd = daoObj.Assets["husd_amount"]
	}
	if daoObj.Assets["seeds_instant_amount"].Amount == 0 {
		a.SeedsLiquid, _ = eos.NewAssetFromString("0.0000 SEEDS")
	} else {
		a.SeedsLiquid = daoObj.Assets["seeds_instant_amount"]
	}
	a.Hypha = daoObj.Assets["hypha_amount"]
	a.Hvoice = daoObj.Assets["hvoice_amount"]
	a.SeedsEscrow = daoObj.Assets["seeds_escrow_amount"]
	a.StartPeriod = periods[daoObj.Ints["start_period"]]
	a.EndPeriod = periods[daoObj.Ints["end_period"]]
	a.DeferredPay = float64(daoObj.Ints["deferred_perc_x100"]) / 100
	a.InstantHusdPerc = float64(daoObj.Ints["instant_husd_perc_x100"]) / 100
	a.CreatedDate = daoObj.CreatedDate
	return a
}

// ProposedPayouts provides the active payout proposals
func ProposedPayouts(ctx context.Context, api *eos.API, periods []Period) []Payout {
	objects := LoadObjects(ctx, api, "proposal")
	var propPayouts []Payout
	for index := range objects {
		daoObject := ToDAOObject(objects[index])
		if daoObject.Names["type"] == "payout" {
			payout := NewPayout(daoObject, periods)
			payout.Approved = false
			propPayouts = append(propPayouts, payout)
		}
	}
	return propPayouts
}

// Payouts provides the set of active approved payouts
func Payouts(ctx context.Context, api *eos.API, periods []Period) []Payout {
	objects := LoadObjects(ctx, api, "payout")
	var payouts []Payout
	for index := range objects {
		payout := NewPayout(ToDAOObject(objects[index]), periods)
		payout.Approved = true
		payouts = append(payouts, payout)
	}
	return payouts
}
