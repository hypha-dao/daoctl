package models

import (
	"fmt"

	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/document-graph/docgraph"
)

// Assignment represents a person assigned to a role for a specific period of time
type Assignment struct {
	// Approved            bool
	Hash            eos.Checksum256
	Title           string
	Description     string
	Owner           eos.Name
	Assigned        eos.Name
	BallotName      eos.Name
	HusdPerPhase    eos.Asset
	HyphaPerPhase   eos.Asset
	HvoicePerPhase  eos.Asset
	DeferredPay     float64
	InstantHusdPerc float64
	TimeShare       float64
	StartPeriod     Period
	PeriodCount     int64
	Document        docgraph.Document
}

func setOrDefault(values map[string]eos.Asset, key string, defaultValue *eos.Asset) *eos.Asset {
	if val, ok := values[key]; ok {
		return &val
	}
	return defaultValue
}

// NewAssignment converts a generic DAO Object to a typed Assignment
func NewAssignment(doc docgraph.Document) (Assignment, error) {

	a := Assignment{}
	a.Document = doc

	titleFv, err := doc.GetContentFromGroup("details", "title")
	if err != nil {
		return Assignment{}, fmt.Errorf("get content failed: %v", err)
	}
	a.Title = titleFv.String()

	descFv, err := doc.GetContentFromGroup("details", "description")
	if err != nil {
		return Assignment{}, fmt.Errorf("get content failed: %v", err)
	}
	a.Description = descFv.String()

	ballotName, err := doc.GetContentFromGroup("system", "ballot_id")
	if err != nil {
		return Assignment{}, fmt.Errorf("get content failed: %v", err)
	}
	a.BallotName, err = ballotName.Name()
	if err != nil {
		return Assignment{}, fmt.Errorf("value downcasting failed: %v", err)
	}

	owner, err := doc.GetContentFromGroup("details", "owner")
	if err != nil {
		return Assignment{}, fmt.Errorf("get content failed: %v", err)
	}
	a.Owner, err = owner.Name()
	if err != nil {
		return Assignment{}, fmt.Errorf("value downcasting failed: %v", err)
	}

	husd, err := doc.GetContentFromGroup("details", "husd_salary_per_phase")
	if err != nil {
		return Assignment{}, fmt.Errorf("get content failed: %v", err)
	}
	a.HusdPerPhase, err = husd.Asset()
	if err != nil {
		return Assignment{}, fmt.Errorf("get content failed: %v", err)
	}

	hypha, err := doc.GetContentFromGroup("details", "hypha_salary_per_phase")
	if err != nil {
		return Assignment{}, fmt.Errorf("get content failed: %v", err)
	}
	a.HyphaPerPhase, err = hypha.Asset()
	if err != nil {
		return Assignment{}, fmt.Errorf("get content failed: %v", err)
	}

	hvoice, err := doc.GetContentFromGroup("details", "hvoice_salary_per_phase")
	if err != nil {
		return Assignment{}, fmt.Errorf("get content failed: %v", err)
	}
	a.HvoicePerPhase, err = hvoice.Asset()
	if err != nil {
		return Assignment{}, fmt.Errorf("get content failed: %v", err)
	}

	periodCount, err := doc.GetContentFromGroup("details", "period_count")
	if err != nil {
		return Assignment{}, fmt.Errorf("missing period_count, cannot continue: %v", err)
	}
	a.PeriodCount, err = periodCount.Int64()
	if err != nil {
		return Assignment{}, fmt.Errorf("get content failed: %v", err)
	}

	timeShare, err := doc.GetContentFromGroup("details", "time_share_x100")
	if err != nil {
		return Assignment{}, fmt.Errorf("missing time_share, cannot continue: %v", err)
	}
	timeShareInt, err := timeShare.Int64()
	if err != nil {
		return Assignment{}, fmt.Errorf("get content failed: %v", err)
	}
	a.TimeShare = float64(timeShareInt) / 100

	deferredPay, err := doc.GetContentFromGroup("details", "deferred_pay_x100")
	if err != nil {
		return Assignment{}, fmt.Errorf("missing deferred_pay, cannot continue: %v", err)
	}
	deferredPayInt, err := deferredPay.Int64()
	if err != nil {
		return Assignment{}, fmt.Errorf("get content failed: %v", err)
	}
	a.DeferredPay = float64(deferredPayInt) / 100

	return a, nil
}
