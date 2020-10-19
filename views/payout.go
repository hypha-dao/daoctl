package views

import (
	"strconv"

	"github.com/hypha-dao/daoctl/util"

	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
)

func payoutHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Receiver"},
			{Align: simpletable.AlignCenter, Text: "Title"},
			{Align: simpletable.AlignCenter, Text: "Deferred %"},
			{Align: simpletable.AlignCenter, Text: "HUSD"},
			{Align: simpletable.AlignCenter, Text: "HYPHA"},
			{Align: simpletable.AlignCenter, Text: "HVOICE"},
			{Align: simpletable.AlignCenter, Text: "Escrow SEEDS"},
			{Align: simpletable.AlignCenter, Text: "Liquid SEEDS"},
			{Align: simpletable.AlignCenter, Text: "Created Date"},
			{Align: simpletable.AlignCenter, Text: "Ballot"},
		},
	}
}

// PayoutTable is a simpleTable.Table object with payouts
func PayoutTable(payouts []models.Payout) *simpletable.Table {

	table := simpletable.New()
	table.Header = payoutHeader()

	husdTotal, _ := eos.NewAssetFromString("0.00 HUSD")
	hvoiceTotal, _ := eos.NewAssetFromString("0.00 HVOICE")
	hyphaTotal, _ := eos.NewAssetFromString("0.00 HYPHA")
	seedsLiquidTotal, _ := eos.NewAssetFromString("0.0000 SEEDS")
	seedsEscrowTotal, _ := eos.NewAssetFromString("0.0000 SEEDS")

	for index := range payouts {

		husdTotal = husdTotal.Add(payouts[index].Husd)
		hyphaTotal = hyphaTotal.Add(payouts[index].Hypha)
		hvoiceTotal = hvoiceTotal.Add(payouts[index].Hvoice)
		seedsLiquidTotal = seedsLiquidTotal.Add(payouts[index].SeedsLiquid)
		seedsEscrowTotal = seedsEscrowTotal.Add(payouts[index].SeedsEscrow)

		r := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: strconv.Itoa(int(payouts[index].ID))},
			{Align: simpletable.AlignRight, Text: string(payouts[index].Receiver)},
			{Align: simpletable.AlignLeft, Text: payouts[index].Title},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(payouts[index].DeferredPay*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&payouts[index].Husd, 0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&payouts[index].Hypha, 0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&payouts[index].Hvoice, 0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&payouts[index].SeedsEscrow, 0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&payouts[index].SeedsLiquid, 0)},
			{Align: simpletable.AlignRight, Text: payouts[index].CreatedDate.Time.Format("2006 Jan 02")},
			{Align: simpletable.AlignRight, Text: string(payouts[index].BallotName)[11:]},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{},
			{},
			{},
			{Align: simpletable.AlignRight, Text: "Subtotal"},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&husdTotal, 0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&hyphaTotal, 0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&hvoiceTotal, 0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&seedsEscrowTotal, 0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&seedsLiquidTotal, 0)},
			{}, {},
		},
	}

	return table
}
