package models

import (
	"context"
	"strconv"

	"github.com/alexeyco/simpletable"
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

func roleHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Title"},
			{Align: simpletable.AlignCenter, Text: "Owner"},
			{Align: simpletable.AlignCenter, Text: "Min Time %"},
			{Align: simpletable.AlignCenter, Text: "Min Def %"},
			{Align: simpletable.AlignCenter, Text: "FTE Cap"},
			{Align: simpletable.AlignCenter, Text: "Annual USD"},
			{Align: simpletable.AlignCenter, Text: "Start Date"},
			{Align: simpletable.AlignCenter, Text: "End Date"},
		},
	}
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

// RoleTable returns a string representing an output table for a Role array
func RoleTable(roles []Role) *simpletable.Table {

	table := simpletable.New()
	table.Header = roleHeader()

	usdTotal, _ := eos.NewAssetFromString("0.00 USD")

	for index := range roles {

		usdTotal = usdTotal.Add(roles[index].AnnualUSDSalary)
		var annualUsdSalary string
		AssetsAsFloats := true
		if AssetsAsFloats {
			annualUsdSalary = strconv.FormatFloat(float64(roles[index].AnnualUSDSalary.Amount/100), 'f', 2, 64)
		} else {
			annualUsdSalary = roles[index].AnnualUSDSalary.String()

		}
		r := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: strconv.Itoa(int(roles[index].ID))},
			{Align: simpletable.AlignLeft, Text: string(roles[index].Title)},
			{Align: simpletable.AlignRight, Text: string(roles[index].Owner)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(roles[index].MinTime*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(roles[index].MinDeferred*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(roles[index].FullTimeCapacity*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: annualUsdSalary},
			{Align: simpletable.AlignRight, Text: roles[index].StartPeriod.StartTime.Time.Format("2006 Jan 02")},
			{Align: simpletable.AlignRight, Text: roles[index].EndPeriod.EndTime.Time.Format("2006 Jan 02")},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{}, {}, {}, {}, {},
			{Align: simpletable.AlignRight, Text: "Subtotal"},
			{Align: simpletable.AlignRight, Text: usdTotal.String()},
			{}, {},
		},
	}

	return table
}
