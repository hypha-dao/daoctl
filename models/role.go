package models

import (
	"context"

	eos "github.com/eoscanada/eos-go"
)

// Role is an approved or proposed role for the DAO
type Role struct {
	ID               uint64
	Approved         bool
	Owner            eos.Name
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

// NewRole creates a new Role instance based on the DAOObject
func NewRole(daoObj DAOObject, periods []Period) Role {
	var r Role
	r.ID = daoObj.ID
	r.Title = daoObj.Strings["title"]
	r.Owner = daoObj.Names["owner"]
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
