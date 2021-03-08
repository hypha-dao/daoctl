package models

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/hypha-dao/daoctl/util"
	"github.com/hypha-dao/document-graph/docgraph"
	"github.com/ryanuber/columnize"

	"github.com/eoscanada/eos-go"
)

// Role is an approved or proposed role for the DAO
type Role struct {
	ID               uint64
	Hash             eos.Checksum256
	Creator          eos.AccountName
	PriorID          uint64
	Approved         bool
	Owner            eos.Name
	BallotName       eos.Name
	Title            string
	Description      string
	URL              string
	AnnualUSDSalary  eos.Asset
	MinTime          float64
	MinDeferred      float64
	FullTimeCapacity float64
	CreatedDate      eos.BlockTimestamp
}

func (r *Role) String() string {
	fteCapCost := util.AssetMult(r.AnnualUSDSalary, big.NewFloat(r.FullTimeCapacity))
	output := []string{
		fmt.Sprintf("Role ID|%v", strconv.Itoa(int(r.ID))),
		fmt.Sprintf("Owner|%v", string(r.Owner)),
		fmt.Sprintf("Title|%v", string(r.Title)),
		fmt.Sprintf("Annual USD Salary|%v", util.FormatAsset(&r.AnnualUSDSalary, 2)),
		fmt.Sprintf("Minimum Time Commitment|%v", strconv.FormatFloat(r.MinTime*100, 'f', -1, 64)),
		fmt.Sprintf("Minimum Deferred Pay|%v", strconv.FormatFloat(r.MinDeferred*100, 'f', -1, 64)),
		fmt.Sprintf("Full Time Capacity|%v", strconv.FormatFloat(r.FullTimeCapacity, 'f', 1, 64)),
		fmt.Sprintf("FTE Cap Cost|%v", util.FormatAsset(&fteCapCost, 2)),
		//fmt.Sprintf("Start Period|%v", r.StartPeriod.StartTime.Time.Format("2006 Jan 02 15:04:05")),
		fmt.Sprintf("Created Date|%v", r.CreatedDate.Time.Format("2006 Jan 02 15:04:05")),
		fmt.Sprintf("Ballot ID|%v", string(r.BallotName)[10:]),
		fmt.Sprintf("Description|%v", r.Description),
	}
	return columnize.SimpleFormat(output)
}

// NewRole creates a new Role instance based on the DAOObject
func NewRole(roleDoc docgraph.Document) (Role, error) {

	var r Role
	r.ID = roleDoc.ID
	r.Hash = roleDoc.Hash
	r.Creator = roleDoc.Creator
	r.CreatedDate = roleDoc.CreatedDate

	titleFv, err := roleDoc.GetContentFromGroup("details", "title")
	if err != nil {
		return Role{}, fmt.Errorf("get content failed: %v", err)
	}
	r.Title = titleFv.String()

	descFv, err := roleDoc.GetContentFromGroup("details", "description")
	if err != nil {
		return Role{}, fmt.Errorf("get content failed: %v", err)
	}
	r.Description = descFv.String()

	ballotName, err := roleDoc.GetContentFromGroup("system", "ballot_id")
	if err != nil {
		return Role{}, fmt.Errorf("get content failed: %v", err)
	}
	r.BallotName, err = ballotName.Name()
	if err != nil {
		return Role{}, fmt.Errorf("get content failed: %v", err)
	}

	owner, err := roleDoc.GetContentFromGroup("details", "owner")
	if err != nil {
		return Role{}, fmt.Errorf("get content failed: %v", err)
	}
	r.Owner, err = owner.Name()
	if err != nil {
		return Role{}, fmt.Errorf("get content failed: %v", err)
	}

	annualUsdSalary, err := roleDoc.GetContentFromGroup("details", "annual_usd_salary")
	if err != nil {
		return Role{}, fmt.Errorf("get content failed: %v", err)
	}
	r.AnnualUSDSalary, err = annualUsdSalary.Asset()
	if err != nil {
		return Role{}, fmt.Errorf("get content failed: %v", err)
	}

	minTime, err := roleDoc.GetContentFromGroup("details", "min_time_share_x100")
	if err != nil {
		r.MinTime = 0
	} else {
		minTimeInt, err := minTime.Int64()
		if err != nil {
			return Role{}, fmt.Errorf("get content failed: %v", err)
		}
		r.MinTime = float64(minTimeInt) / 100
	}

	fullTimeCap, err := roleDoc.GetContentFromGroup("details", "fulltime_capacity_x100")
	if err != nil {
		r.FullTimeCapacity = 1
	} else {
		ftc, err := fullTimeCap.Int64()
		if err != nil {
			return Role{}, fmt.Errorf("get content failed: %v", err)
		}
		r.FullTimeCapacity = float64(ftc) / 100
	}

	minDeferred, err := roleDoc.GetContentFromGroup("details", "min_deferred_x100")
	if err != nil {
		r.MinDeferred = 0
	} else {
		minDef, err := minDeferred.Int64()
		if err != nil {
			return Role{}, fmt.Errorf("get content failed: %v", err)
		}
		r.MinDeferred = float64(minDef) / 100
	}

	return r, nil
}
