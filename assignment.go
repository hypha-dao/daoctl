package main

import (
	"context"
	"fmt"
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
	RoleID              uint64
	StartPeriod         Period
	EndPeriod           Period
	CreatedDate         eos.BlockTimestamp
}

func assignmentHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Assigned"},
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

// ToAssignment converts a generic DAO Object to a typed Assignment
func toAssignment(daoObj DAOObject, periods []Period) Assignment {
	var a Assignment
	a.ID = daoObj.ID
	a.Owner = daoObj.Names["owner"]
	a.Assigned = daoObj.Names["assigned_account"]
	a.HusdPerPhase = daoObj.Assets["husd_salary_per_phase"]
	a.HyphaPerPhase = daoObj.Assets["hypha_salary_per_phase"]
	a.HvoicePerPhase = daoObj.Assets["hvoice_salary_per_phase"]
	a.SeedsEscrowPerPhase = daoObj.Assets["seeds_escrow_salary_per_phase"]
	a.SeedsLiquidPerPhase = daoObj.Assets["seeds_instant_salary_per_phase"]
	a.RoleID = daoObj.Ints["role_id"]
	a.StartPeriod = periods[daoObj.Ints["start_period"]]
	a.EndPeriod = periods[daoObj.Ints["end_period"]]
	a.TimeShare = float64(daoObj.Ints["time_share_x100"]) / 100
	a.DeferredPay = float64(daoObj.Ints["deferred_perc_x100"]) / 100
	a.InstantHusdPerc = float64(daoObj.Ints["instant_husd_perc_x100"]) / 100
	a.CreatedDate = daoObj.CreatedDate
	return a
}

// ProposedAssignments provides the active assignment proposals
func ProposedAssignments(ctx context.Context, api *eos.API, periods []Period) []Assignment {
	objects := LoadObjects(ctx, api, "proposal")
	var propAssignments []Assignment
	for index := range objects {
		daoObject := ToDAOObject(objects[index])
		if daoObject.Names["type"] == "assignment" {
			assignment := toAssignment(daoObject, periods)
			assignment.Approved = false
			propAssignments = append(propAssignments, assignment)
		}
	}
	return propAssignments
}

// Assignments provides the set of active approved assignments
func Assignments(ctx context.Context, api *eos.API, periods []Period) []Assignment {
	objects := LoadObjects(ctx, api, "assignment")
	var assignments []Assignment
	for index := range objects {
		assignment := toAssignment(ToDAOObject(objects[index]), periods)
		assignment.Approved = true
		assignments = append(assignments, assignment)
	}
	return assignments
}

func assignmentTable(assignments []Assignment) string {

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

		r := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: strconv.Itoa(int(assignments[index].ID))},
			{Align: simpletable.AlignRight, Text: string(assignments[index].Assigned)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(assignments[index].TimeShare*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(assignments[index].DeferredPay*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(assignments[index].InstantHusdPerc*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: assignments[index].HusdPerPhase.String()},
			{Align: simpletable.AlignRight, Text: assignments[index].HyphaPerPhase.String()},
			{Align: simpletable.AlignRight, Text: assignments[index].HvoicePerPhase.String()},
			{Align: simpletable.AlignRight, Text: assignments[index].SeedsEscrowPerPhase.String()},
			{Align: simpletable.AlignRight, Text: assignments[index].SeedsLiquidPerPhase.String()},
			{Align: simpletable.AlignRight, Text: assignments[index].StartPeriod.StartTime.Time.Format("2006 Jan 02")},
			{Align: simpletable.AlignRight, Text: assignments[index].EndPeriod.EndTime.Time.Format("2006 Jan 02")},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{},
			{},
			{}, {},
			{Align: simpletable.AlignRight, Text: "Subtotal"},
			{Align: simpletable.AlignRight, Text: husdTotal.String()},
			{Align: simpletable.AlignRight, Text: hyphaTotal.String()},
			{Align: simpletable.AlignRight, Text: hvoiceTotal.String()},
			{Align: simpletable.AlignRight, Text: seedsEscrowTotal.String()},
			{Align: simpletable.AlignRight, Text: seedsLiquidTotal.String()},
			{}, {},
		},
	}

	table.SetStyle(simpletable.StyleCompactLite)
	return table.String()
}

// PrintAssignments prints a table with all active assignments
func PrintAssignments(ctx context.Context, api *eos.API, periods []Period, includeProposals bool) {

	assignments := Assignments(ctx, api, periods)
	fmt.Println("\n\n" + assignmentTable(assignments) + "\n\n")

	if includeProposals {
		propAssignments := ProposedAssignments(ctx, api, periods)
		fmt.Println("\n\n" + assignmentTable(propAssignments) + "\n\n")
	}
}
