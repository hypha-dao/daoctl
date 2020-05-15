package models

import (
	"context"

	eos "github.com/eoscanada/eos-go"
)

// Assignment represents a person assigned to a role for a specific period of time
type Assignment struct {
	ID                  uint64
	Approved            bool
	Owner               eos.Name
	Assigned            eos.Name
	BallotName          eos.Name
	HusdPerPhase        eos.Asset
	HyphaPerPhase       eos.Asset
	HvoicePerPhase      eos.Asset
	SeedsEscrowPerPhase eos.Asset
	SeedsLiquidPerPhase eos.Asset
	DeferredPay         float64
	InstantHusdPerc     float64
	TimeShare           float64
	Role                Role
	StartPeriod         Period
	EndPeriod           Period
	CreatedDate         eos.BlockTimestamp
}

func scopeApprovals(scope string) bool {
	if scope == "assignment" {
		return true
	}
	return false
}

// NewAssignment converts a generic DAO Object to a typed Assignment
func NewAssignment(daoObj DAOObject, roles []Role, periods []Period) Assignment {
	var a Assignment
	a.ID = daoObj.ID
	a.Owner = daoObj.Names["owner"]
	a.Assigned = daoObj.Names["assigned_account"]
	a.BallotName = daoObj.Names["ballot_id"]
	a.HusdPerPhase = daoObj.Assets["husd_salary_per_phase"]
	a.HyphaPerPhase = daoObj.Assets["hypha_salary_per_phase"]
	a.HvoicePerPhase = daoObj.Assets["hvoice_salary_per_phase"]
	a.SeedsEscrowPerPhase = daoObj.Assets["seeds_escrow_salary_per_phase"]
	a.SeedsLiquidPerPhase = daoObj.Assets["seeds_instant_salary_per_phase"]
	a.Role = roles[daoObj.Ints["role_id"]]
	a.StartPeriod = periods[daoObj.Ints["start_period"]]
	a.EndPeriod = periods[daoObj.Ints["end_period"]]
	a.TimeShare = float64(daoObj.Ints["time_share_x100"]) / 100
	a.DeferredPay = float64(daoObj.Ints["deferred_perc_x100"]) / 100
	a.InstantHusdPerc = float64(daoObj.Ints["instant_husd_perc_x100"]) / 100
	a.CreatedDate = daoObj.CreatedDate
	return a
}

// Assignments provides the set of active approved assignments
func Assignments(ctx context.Context, api *eos.API, roles []Role, periods []Period, scope string) []Assignment {
	objects := LoadObjects(ctx, api, scope)
	var assignments []Assignment
	for index := range objects {
		daoObject := ToDAOObject(objects[index])
		if daoObject.Names["type"] == "assignment" {
			assignment := NewAssignment(daoObject, roles, periods)
			assignment.Approved = scopeApprovals(scope)
			assignments = append(assignments, assignment)
		}
	}
	return assignments
}
