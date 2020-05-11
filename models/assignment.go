package models

import (
	"context"
	"strconv"

	"github.com/alexeyco/simpletable"
	eos "github.com/eoscanada/eos-go"
)

// Assignment represents a person assigned to a role for a specific period of time
type Assignment struct {
	ID                  uint64
	Approved            bool
	Owner               eos.Name
	Assigned            eos.Name
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

func assignmentHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Assigned"},
			{Align: simpletable.AlignCenter, Text: "Role"},
			{Align: simpletable.AlignCenter, Text: "Role Annually"},
			{Align: simpletable.AlignCenter, Text: "Time %"},
			{Align: simpletable.AlignCenter, Text: "Deferred %"},
			{Align: simpletable.AlignCenter, Text: "HUSD %"},
			{Align: simpletable.AlignCenter, Text: "HUSD"},
			{Align: simpletable.AlignCenter, Text: "HYPHA"},
			{Align: simpletable.AlignCenter, Text: "HVOICE"},
			{Align: simpletable.AlignCenter, Text: "Escrow SEEDS"},
			{Align: simpletable.AlignCenter, Text: "Liquid SEEDS"},
			{Align: simpletable.AlignCenter, Text: "Start Date"},
			{Align: simpletable.AlignCenter, Text: "End Date"},
		},
	}
}

// NewAssignment converts a generic DAO Object to a typed Assignment
func NewAssignment(daoObj DAOObject, roles []Role, periods []Period) Assignment {
	var a Assignment
	a.ID = daoObj.ID
	a.Owner = daoObj.Names["owner"]
	a.Assigned = daoObj.Names["assigned_account"]
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

// ProposedAssignments provides the active assignment proposals
func ProposedAssignments(ctx context.Context, api *eos.API, roles []Role, periods []Period) []Assignment {
	objects := LoadObjects(ctx, api, "proposal")
	var propAssignments []Assignment
	for index := range objects {
		daoObject := ToDAOObject(objects[index])
		if daoObject.Names["type"] == "assignment" {
			assignment := NewAssignment(daoObject, roles, periods)
			assignment.Approved = false
			propAssignments = append(propAssignments, assignment)
		}
	}
	return propAssignments
}

// Assignments provides the set of active approved assignments
func Assignments(ctx context.Context, api *eos.API, roles []Role, periods []Period) []Assignment {
	objects := LoadObjects(ctx, api, "assignment")
	var assignments []Assignment
	for index := range objects {
		assignment := NewAssignment(ToDAOObject(objects[index]), roles, periods)
		assignment.Approved = true
		assignments = append(assignments, assignment)
	}
	return assignments
}

// AssignmentTable returns a string representing a table of the assignnments
func AssignmentTable(assignments []Assignment) *simpletable.Table {

	table := simpletable.New()
	table.Header = assignmentHeader()

	husdTotal, _ := eos.NewAssetFromString("0.00 HUSD")
	hvoiceTotal, _ := eos.NewAssetFromString("0.00 HVOICE")
	hyphaTotal, _ := eos.NewAssetFromString("0.00 HYPHA")
	seedsLiquidTotal, _ := eos.NewAssetFromString("0.0000 SEEDS")
	seedsEscrowTotal, _ := eos.NewAssetFromString("0.0000 SEEDS")

	for index := range assignments {

		husdTotal = husdTotal.Add(assignments[index].HusdPerPhase)
		hyphaTotal = hyphaTotal.Add(assignments[index].HyphaPerPhase)
		hvoiceTotal = hvoiceTotal.Add(assignments[index].HvoicePerPhase)
		seedsLiquidTotal = seedsLiquidTotal.Add(assignments[index].SeedsLiquidPerPhase)
		seedsEscrowTotal = seedsEscrowTotal.Add(assignments[index].SeedsEscrowPerPhase)

		AssetAsFloats := true
		var annualUsdSalary, husdPerPhase, hyphaPerPhase, hvoicePerPhase, seedsEscrowPerPhase, seedsLiquidPerPhase string
		if AssetAsFloats {
			annualUsdSalary = strconv.FormatFloat(float64(assignments[index].Role.AnnualUSDSalary.Amount/100), 'f', 2, 64)
			husdPerPhase = strconv.FormatFloat(float64(assignments[index].HusdPerPhase.Amount/100), 'f', 2, 64)
			hyphaPerPhase = strconv.FormatFloat(float64(assignments[index].HyphaPerPhase.Amount/100), 'f', 2, 64)
			hvoicePerPhase = strconv.FormatFloat(float64(assignments[index].HvoicePerPhase.Amount/100), 'f', 2, 64)
			seedsEscrowPerPhase = strconv.FormatFloat(float64(assignments[index].SeedsEscrowPerPhase.Amount/10000), 'f', 2, 64)
			seedsLiquidPerPhase = strconv.FormatFloat(float64(assignments[index].SeedsLiquidPerPhase.Amount/10000), 'f', 2, 64)
		} else {
			annualUsdSalary = assignments[index].Role.AnnualUSDSalary.String()
			husdPerPhase = assignments[index].HusdPerPhase.String()
			hyphaPerPhase = assignments[index].HyphaPerPhase.String()
			hvoicePerPhase = assignments[index].HvoicePerPhase.String()
			seedsEscrowPerPhase = assignments[index].SeedsEscrowPerPhase.String()
			seedsLiquidPerPhase = assignments[index].SeedsLiquidPerPhase.String()
		}

		r := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: strconv.Itoa(int(assignments[index].ID))},
			{Align: simpletable.AlignRight, Text: string(assignments[index].Assigned)},
			{Align: simpletable.AlignLeft, Text: string(assignments[index].Role.Title)},
			{Align: simpletable.AlignLeft, Text: annualUsdSalary},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(assignments[index].TimeShare*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(assignments[index].DeferredPay*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(assignments[index].InstantHusdPerc*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: husdPerPhase},
			{Align: simpletable.AlignRight, Text: hyphaPerPhase},
			{Align: simpletable.AlignRight, Text: hvoicePerPhase},
			{Align: simpletable.AlignRight, Text: seedsEscrowPerPhase},
			{Align: simpletable.AlignRight, Text: seedsLiquidPerPhase},
			{Align: simpletable.AlignRight, Text: assignments[index].StartPeriod.StartTime.Time.Format("2006 Jan 02")},
			{Align: simpletable.AlignRight, Text: assignments[index].EndPeriod.EndTime.Time.Format("2006 Jan 02")},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{},
			{},
			{}, {}, {}, {},
			{Align: simpletable.AlignRight, Text: "Subtotal"},
			{Align: simpletable.AlignRight, Text: husdTotal.String()},
			{Align: simpletable.AlignRight, Text: hyphaTotal.String()},
			{Align: simpletable.AlignRight, Text: hvoiceTotal.String()},
			{Align: simpletable.AlignRight, Text: seedsEscrowTotal.String()},
			{Align: simpletable.AlignRight, Text: seedsLiquidTotal.String()},
			{}, {},
		},
	}

	return table
}
