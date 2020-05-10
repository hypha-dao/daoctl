package models

import (
	"context"
	"fmt"
	"strconv"

	"github.com/alexeyco/simpletable"
	eos "github.com/eoscanada/eos-go"
)

// Payout represents a person assigned to a role for a specific period of time
type Payout struct {
	ID              uint64
	Approved        bool
	Receiver        eos.Name
	Title           string
	Description     string
	Husd            eos.Asset
	Hypha           eos.Asset
	Hvoice          eos.Asset
	SeedsEscrow     eos.Asset
	SeedsLiquid     eos.Asset
	DeferredPay     float64
	InstantHusdPerc float64
	StartPeriod     Period
	EndPeriod       Period
	CreatedDate     eos.BlockTimestamp
}

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
		},
	}
}

// ToPayout converts a generic DAO Object to a typed Payout
func toPayout(daoObj DAOObject, periods []Period) Payout {
	var a Payout
	a.ID = daoObj.ID
	a.Receiver = daoObj.Names["recipient"]
	a.Title = daoObj.Strings["title"]
	a.Husd = daoObj.Assets["husd_amount"]

	if daoObj.Assets["husd_amount"].Amount == 0 {
		a.Husd, _ = eos.NewAssetFromString("0.00 HUSD")
	} else {
		a.Husd = daoObj.Assets["husd_amount"]
	}
	if daoObj.Assets["seeds_instant_amount"].Amount == 0 {
		a.SeedsLiquid, _ = eos.NewAssetFromString("0.0000 SEEDS")
	} else {
		a.SeedsLiquid = daoObj.Assets["seeds_instant_amount"]
	}
	a.Hypha = daoObj.Assets["hypha_amount"]
	a.Hvoice = daoObj.Assets["hvoice_amount"]
	a.SeedsEscrow = daoObj.Assets["seeds_escrow_amount"]
	a.StartPeriod = periods[daoObj.Ints["start_period"]]
	a.EndPeriod = periods[daoObj.Ints["end_period"]]
	a.DeferredPay = float64(daoObj.Ints["deferred_perc_x100"]) / 100
	a.InstantHusdPerc = float64(daoObj.Ints["instant_husd_perc_x100"]) / 100
	a.CreatedDate = daoObj.CreatedDate
	return a
}

// ProposedPayouts provides the active payout proposals
func ProposedPayouts(ctx context.Context, api *eos.API, periods []Period) []Payout {
	objects := LoadObjects(ctx, api, "proposal")
	var propPayouts []Payout
	for index := range objects {
		daoObject := ToDAOObject(objects[index])
		if daoObject.Names["type"] == "payout" {
			payout := toPayout(daoObject, periods)
			payout.Approved = false
			propPayouts = append(propPayouts, payout)
		}
	}
	return propPayouts
}

// Payouts provides the set of active approved payouts
func Payouts(ctx context.Context, api *eos.API, periods []Period) []Payout {
	objects := LoadObjects(ctx, api, "payout")
	var payouts []Payout
	for index := range objects {
		payout := toPayout(ToDAOObject(objects[index]), periods)
		payout.Approved = true
		payouts = append(payouts, payout)
	}
	return payouts
}

func payoutTable(payouts []Payout) string {

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
			{Align: simpletable.AlignRight, Text: payouts[index].Husd.String()},
			{Align: simpletable.AlignRight, Text: payouts[index].Hypha.String()},
			{Align: simpletable.AlignRight, Text: payouts[index].Hvoice.String()},
			{Align: simpletable.AlignRight, Text: payouts[index].SeedsEscrow.String()},
			{Align: simpletable.AlignRight, Text: payouts[index].SeedsLiquid.String()},
			{Align: simpletable.AlignRight, Text: payouts[index].CreatedDate.Time.Format("2006 Jan 02")},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{},
			{},
			{},
			{Align: simpletable.AlignRight, Text: "Subtotal"},
			{Align: simpletable.AlignRight, Text: husdTotal.String()},
			{Align: simpletable.AlignRight, Text: hyphaTotal.String()},
			{Align: simpletable.AlignRight, Text: hvoiceTotal.String()},
			{Align: simpletable.AlignRight, Text: seedsEscrowTotal.String()},
			{Align: simpletable.AlignRight, Text: seedsLiquidTotal.String()},
			{},
		},
	}

	table.SetStyle(simpletable.StyleCompactLite)
	return table.String()
}

// PrintPayouts prints a table with all active payouts
func PrintPayouts(ctx context.Context, api *eos.API, periods []Period, includeProposals bool) {

	payouts := Payouts(ctx, api, periods)
	fmt.Println("\n\n" + payoutTable(payouts) + "\n\n")

	if includeProposals {
		propPayouts := ProposedPayouts(ctx, api, periods)
		fmt.Println("\n\n" + payoutTable(propPayouts) + "\n\n")
	}
}
