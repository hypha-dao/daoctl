package models

import (
	"context"
	"fmt"
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
	StartPeriod      Period
	EndPeriod        Period
	CreatedDate      eos.BlockTimestamp
}

func (r *Role) String() string {
	// fteCapCost := util.AssetMult(r.AnnualUSDSalary, big.NewFloat(r.FullTimeCapacity))
	output := []string{
		fmt.Sprintf("Role ID|%v", strconv.Itoa(int(r.ID))),
		// fmt.Sprintf("Prior ID|%v", strconv.Itoa(int(r.PriorID))),
		// fmt.Sprintf("Owner|%v", string(r.Owner)),
		fmt.Sprintf("Title|%v", string(r.Title)),
		// fmt.Sprintf("URL|%v", string(r.URL)),
		fmt.Sprintf("Annual USD Salary|%v", util.FormatAsset(&r.AnnualUSDSalary, 2)),
		fmt.Sprintf("Minimum Time Commitment|%v", strconv.FormatFloat(r.MinTime*100, 'f', -1, 64)),
		fmt.Sprintf("Minimum Deferred Pay|%v", strconv.FormatFloat(r.MinDeferred*100, 'f', -1, 64)),
		// fmt.Sprintf("Full Time Capacity|%v", strconv.FormatFloat(r.FullTimeCapacity, 'f', 1, 64)),
		// fmt.Sprintf("FTE Cap Cost|%v", util.FormatAsset(&fteCapCost, 2)),
		fmt.Sprintf("Start Period|%v", r.StartPeriod.StartTime.Time.Format("2006 Jan 02 15:04:05")),
		fmt.Sprintf("End Period|%v", r.EndPeriod.EndTime.Time.Format("2006 Jan 02 15:04:05")),
		fmt.Sprintf("Created Date|%v", r.CreatedDate.Time.Format("2006 Jan 02 15:04:05")),
		fmt.Sprintf("Ballot ID|%v", string(r.BallotName)[10:]),
		fmt.Sprintf("Description|%v", r.Description),
	}
	return columnize.SimpleFormat(output)
}

// GetContentAsStringOrFail returns a string value of found content or it panics
func GetContentAsStringOrFail(d docgraph.Document, label string) string {
	fv, err := d.GetContent(label)
	if err != nil {
		panic("get content failed: label: %v")
	}
	return fv.String()
}

// GetContentAsName returns a string value of found content or it panics
func GetContentAsName(d docgraph.Document, label string) (eos.Name, error) {
	fv, err := d.GetContent(label)
	if err != nil {
		return eos.Name("error"), fmt.Errorf("get content as name failed: %v", err)
	}
	switch v := fv.Impl.(type) {
	case eos.Name:
		return v, nil
	case string:
		return eos.Name(v), nil
	default:
		return eos.Name("error"), fmt.Errorf("get content as name failed: %v", err)
	}
}

// GetContentAsInt returns an int64 value of found content or it panics
func GetContentAsInt(d docgraph.Document, label string) (int64, error) {
	fv, err := d.GetContent(label)
	if err != nil {
		return -1, fmt.Errorf("get content as int failed: %v", err)
	}
	switch v := fv.Impl.(type) {
	case int64:
		return v, nil
	default:
		return -1, fmt.Errorf("get content as int failed: %v", err)
	}
}

// GetContentAsAsset returns a string value of found content or it panics
func GetContentAsAsset(d docgraph.Document, label string) (eos.Asset, error) {
	fv, err := d.GetContent(label)
	if err != nil {
		return eos.Asset{}, fmt.Errorf("get content as asset failed: %v", err)
	}
	switch v := fv.Impl.(type) {
	case *eos.Asset:
		return *v, nil
	default:
		return eos.Asset{}, fmt.Errorf("get content as asset failed: %v", err)
	}
}

// NewRole creates a new Role instance based on the DAOObject
func NewRole(roleDoc docgraph.Document, periods []Period) (Role, error) {

	var r Role
	r.ID = roleDoc.ID
	r.Hash = roleDoc.Hash
	r.Title = GetContentAsStringOrFail(roleDoc, "title")
	r.Creator = roleDoc.Creator
	r.Description = GetContentAsStringOrFail(roleDoc, "description")
	r.CreatedDate = roleDoc.CreatedDate

	ballotName, err := GetContentAsName(roleDoc, "ballot_id")
	if err != nil {
		return Role{}, fmt.Errorf("get content failed: %v", err)
	}
	r.BallotName = ballotName

	annualUsdSalary, err := GetContentAsAsset(roleDoc, "annual_usd_salary")
	r.AnnualUSDSalary = annualUsdSalary

	minTime, err := GetContentAsInt(roleDoc, "min_time_share_x100")
	if err != nil {
		return Role{}, fmt.Errorf("get content failed: %v", err)
	}
	r.MinTime = float64(minTime) / 100

	minDeferred, err := GetContentAsInt(roleDoc, "min_deferred_x100")
	if err != nil {
		return Role{}, fmt.Errorf("get content failed: %v", err)
	}
	r.MinDeferred = float64(minDeferred) / 100

	startPeriod, err := GetContentAsInt(roleDoc, "start_period")
	if err != nil {
		return Role{}, fmt.Errorf("get content failed: %v", err)
	}
	r.StartPeriod = periods[startPeriod]

	endPeriod, err := GetContentAsInt(roleDoc, "end_period")
	if err != nil {
		return Role{}, fmt.Errorf("get content failed: %v", err)
	}
	r.EndPeriod = periods[endPeriod]

	return r, nil
}

// NewRoleByID loads a single role based on its ID number
func NewRoleByID(ctx context.Context, api *eos.API, periods []Period, ID uint64) Role {
	return Role{}
	// roleDoc := LoadDocument(ctx, api, "role", ID)
	// return NewRole(roleDoc, periods)
}

// Roles provides the set of active approved roles
func Roles(ctx context.Context, api *eos.API, periods []Period, scope string) []Role {
	return []Role{}
	// objects := LoadObjects(ctx, api, scope)
	// var roles []Role
	// for index := range objects {
	// 	roleDocect := ToDocument(objects[index])
	// 	if roleDocect.Names["type"] == "role" {
	// 		role := NewRole(roleDocect, periods)
	// 		role.Approved = scopeApprovals(scope)
	// 		roles = append(roles, role)
	// 	}
	// }
	// return roles
}
