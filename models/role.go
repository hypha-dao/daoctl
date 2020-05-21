package models

import (
  "context"
  "fmt"
  "github.com/hypha-dao/daoctl/util"
  "github.com/ryanuber/columnize"
  "math/big"
  "strconv"

  "github.com/eoscanada/eos-go"
)

// Role is an approved or proposed role for the DAO
type Role struct {
	ID               uint64
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
  fteCapCost := util.AssetMult(r.AnnualUSDSalary, big.NewFloat(r.FullTimeCapacity))
	output := []string{
		fmt.Sprintf("Role ID|%v", strconv.Itoa(int(r.ID))),
		fmt.Sprintf("Prior ID|%v", strconv.Itoa(int(r.PriorID))),
		fmt.Sprintf("Owner|%v", string(r.Owner)),
		fmt.Sprintf("Title|%v", string(r.Title)),
		fmt.Sprintf("URL|%v", string(r.URL)),
		fmt.Sprintf("Annual USD Salary|%v", util.FormatAsset(&r.AnnualUSDSalary)),
		fmt.Sprintf("Minimum Time Commitment|%v", strconv.FormatFloat(r.MinTime*100, 'f', -1, 64)),
		fmt.Sprintf("Minimum Deferred Pay|%v", strconv.FormatFloat(r.MinDeferred*100, 'f', -1, 64)),
		fmt.Sprintf("Full Time Capacity|%v", strconv.FormatFloat(r.FullTimeCapacity, 'f', 1, 64)),
		fmt.Sprintf("FTE Cap Cost|%v", util.FormatAsset(&fteCapCost)),
		fmt.Sprintf("Start Period|%v", r.StartPeriod.StartTime.Time.Format("2006 Jan 02 15:04:05")),
		fmt.Sprintf("End Period|%v", r.EndPeriod.EndTime.Time.Format("2006 Jan 02 15:04:05")),
		fmt.Sprintf("Created Date|%v", r.CreatedDate.Time.Format("2006 Jan 02 15:04:05")),
		fmt.Sprintf("Ballot ID|%v", string(r.BallotName)[11:]),
		fmt.Sprintf("Description|%v", r.Description),
	}
	return columnize.SimpleFormat(output)
}

// NewRole creates a new Role instance based on the DAOObject
func NewRole(daoObj DAOObject, periods []Period) Role {
	var r Role
	r.ID = daoObj.ID
	r.PriorID = daoObj.Ints["prior_id"]
	r.Title = daoObj.Strings["title"]
	r.Owner = daoObj.Names["owner"]
	r.BallotName = daoObj.Names["ballot_id"]
	r.Description = daoObj.Strings["description"]
	r.URL = daoObj.Strings["url"]
	r.AnnualUSDSalary = daoObj.Assets["annual_usd_salary"]
	r.MinTime = float64(daoObj.Ints["min_time_share_x100"]) / 100
	r.MinDeferred = float64(daoObj.Ints["min_deferred_x100"]) / 100
	r.FullTimeCapacity = float64(daoObj.Ints["fulltime_capacity_x100"]) / 100
	r.StartPeriod = periods[daoObj.Ints["start_period"]]
	r.EndPeriod = periods[daoObj.Ints["end_period"]]
	r.CreatedDate = daoObj.CreatedDate
	return r
}

// NewRoleByID loads a single role based on its ID number
func NewRoleByID(ctx context.Context, api *eos.API, periods []Period, ID uint64) Role {
	daoObj := LoadObject(ctx, api, "role", ID)
	return NewRole(daoObj, periods)
}

// ProposedRoles provides the set of active approved roles
func ProposedRoles(ctx context.Context, api *eos.API, periods []Period) []Role {
	objects := LoadObjects(ctx, api, "proposal")
	var roles []Role
	for index := range objects {
		daoObject := ToDAOObject(objects[index])
		if daoObject.Names["type"] == "role" {
			role := NewRole(ToDAOObject(objects[index]), periods)
			role.Approved = true
			roles = append(roles, role)
		}
	}
	return roles
}

// Roles provides the set of active approved roles
func Roles(ctx context.Context, api *eos.API, periods []Period) []Role {
	objects := LoadObjects(ctx, api, "role")
	var roles []Role
	for index := range objects {
		role := NewRole(ToDAOObject(objects[index]), periods)
		role.Approved = true
		roles = append(roles, role)
	}
	return roles
}
