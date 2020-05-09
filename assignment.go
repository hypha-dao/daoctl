package main

import "github.com/eoscanada/eos-go"

type Assignment struct {
	ID                  uint64
	Owner               eos.Name
	Assigned            eos.Name
	HusdPerPhase        eos.Asset
	HyphaPerPhase       eos.Asset
	HvoicePerPhase      eos.Asset
	SeedsEscrowPerPhase eos.Asset
	SeedsLiquidPerPhase eos.Asset
	DeferredPay         float32
	InstantHusdPerc     float32
	TimeShare           float32
	RoleID              uint64
	StartPeriodID       uint64
	EndPeriodID         uint64
	CreatedDate         eos.BlockTimestamp
	UpdatedDate         eos.BlockTimestamp
}

func ToAssignment(daoObj DAOObject) Assignment {
	var a Assignment
	a.Owner = daoObj.Names["owner"]
	a.Assigned = daoObj.Names["assigned_account"]
	a.HusdPerPhase = daoObj.Assets["husd_salary_per_phase"]
	a.HyphaPerPhase = daoObj.Assets["hypha_salary_per_phase"]
	a.HvoicePerPhase = daoObj.Assets["hvoice_salary_per_phase"]
	a.SeedsEscrowPerPhase = daoObj.Assets["seeds_escrow_salary_per_phase"]
	a.SeedsLiquidPerPhase = daoObj.Assets["seeds_instant_salary_per_phase"]
	a.RoleID = daoObj.Ints["role_id"]
	a.StartPeriodID = daoObj.Ints["start_period"]
	a.EndPeriodID = daoObj.Ints["end_period"]
	a.TimeShare = float32(daoObj.Ints["time_share_x100"]) / 100
	a.DeferredPay = float32(daoObj.Ints["deferred_perc_x100"]) / 100
	a.InstantHusdPerc = float32(daoObj.Ints["instant_husd_perc_x100"]) / 100
	a.CreatedDate = daoObj.CreatedDate
	a.UpdatedDate = daoObj.UpdatedDate
	return a
}
