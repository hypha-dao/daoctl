package views

import (
  "github.com/hypha-dao/daoctl/util"
  "strconv"

	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
)

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
			{Align: simpletable.AlignCenter, Text: "Ballot"},
		},
	}
}

// AssignmentTable returns a string representing a table of the assignnments
func AssignmentTable(assignments []models.Assignment) *simpletable.Table {

	table := simpletable.New()
	table.Header = assignmentHeader()

	husdTotal, _ := eos.NewAssetFromString("0.00 HUSD")
	hvoiceTotal, _ := eos.NewAssetFromString("0.00 HVOICE")
	hyphaTotal, _ := eos.NewAssetFromString("0.00 HYPHA")
	seedsLiquidTotal, _ := eos.NewAssetFromString("0.0000 SEEDS")
	seedsEscrowTotal, _ := eos.NewAssetFromString("0.0000 SEEDS")

	for index := range assignments {

		if assignments[index].HusdPerPhase.Symbol.Symbol == "USD" {
			assignments[index].HusdPerPhase.Symbol.Symbol = "HUSD"
		}
		husdTotal = husdTotal.Add(assignments[index].HusdPerPhase)
		hyphaTotal = hyphaTotal.Add(assignments[index].HyphaPerPhase)
		hvoiceTotal = hvoiceTotal.Add(assignments[index].HvoicePerPhase)
		seedsLiquidTotal = seedsLiquidTotal.Add(assignments[index].SeedsLiquidPerPhase)
		seedsEscrowTotal = seedsEscrowTotal.Add(assignments[index].SeedsEscrowPerPhase)

		r := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: strconv.Itoa(int(assignments[index].ID))},
			{Align: simpletable.AlignRight, Text: string(assignments[index].Assigned)},
			{Align: simpletable.AlignLeft, Text: string(assignments[index].Role.Title)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&assignments[index].Role.AnnualUSDSalary, 0)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(assignments[index].TimeShare*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(assignments[index].DeferredPay*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(assignments[index].InstantHusdPerc*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&assignments[index].HusdPerPhase,0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&assignments[index].HyphaPerPhase,0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&assignments[index].HvoicePerPhase,0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&assignments[index].SeedsEscrowPerPhase,0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&assignments[index].SeedsLiquidPerPhase,0)},
			{Align: simpletable.AlignRight, Text: assignments[index].StartPeriod.StartTime.Time.Format("2006 Jan 02")},
			{Align: simpletable.AlignRight, Text: assignments[index].EndPeriod.EndTime.Time.Format("2006 Jan 02")},
			{Align: simpletable.AlignRight, Text: string(assignments[index].BallotName)[11:]},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{},
			{},
			{}, {}, {}, {},
			{Align: simpletable.AlignRight, Text: "Subtotal"},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&husdTotal,0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&hyphaTotal,0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&hvoiceTotal,0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&seedsEscrowTotal,0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&seedsLiquidTotal,0)},
			{}, {}, {},
		},
	}

	return table
}
