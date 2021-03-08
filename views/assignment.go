package views

import (
	"strconv"

	"github.com/hypha-dao/daoctl/util"

	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
)

func assignmentHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Hash"},
			{Align: simpletable.AlignCenter, Text: "Title"},
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

	for index := range assignments {

		husdTotal = husdTotal.Add(assignments[index].HusdPerPhase)
		hyphaTotal = hyphaTotal.Add(assignments[index].HyphaPerPhase)
		hvoiceTotal = hvoiceTotal.Add(assignments[index].HvoicePerPhase)

		r := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: assignments[index].Hash.String()[:5]},
			{Align: simpletable.AlignRight, Text: string(assignments[index].Assigned)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(assignments[index].TimeShare*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(assignments[index].DeferredPay*100, 'f', 0, 64)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(assignments[index].InstantHusdPerc*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&assignments[index].HusdPerPhase, 0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&assignments[index].HyphaPerPhase, 0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&assignments[index].HvoicePerPhase, 0)},
			// {Align: simpletable.AlignRight, Text: assignments[index].StartPeriod.StartTime.Time.Format("2006 Jan 02")},
			{Align: simpletable.AlignRight, Text: string(assignments[index].BallotName)[10:]},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{},
			{},
			{}, {}, {}, {},
			{Align: simpletable.AlignRight, Text: "Subtotal"},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&husdTotal, 0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&hyphaTotal, 0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&hvoiceTotal, 0)},
			{}, {}, {},
		},
	}

	return table
}
