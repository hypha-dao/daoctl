package models

import (
	"context"
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

// NewPayout converts a generic DAO Object to a typed Payout
func NewPayout(daoObj DAOObject, periods []Period) Payout {
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
			payout := NewPayout(daoObject, periods)
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
		payout := NewPayout(ToDAOObject(objects[index]), periods)
		payout.Approved = true
		payouts = append(payouts, payout)
	}
	return payouts
}

// PayoutTable is a simpleTable.Table object with payouts
func PayoutTable(payouts []Payout) *simpletable.Table {

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

		AssetAsFloats := true
		var husd, hypha, hvoice, seedsEscrow, seedsLiquid string
		if AssetAsFloats {
			husd = strconv.FormatFloat(float64(payouts[index].Husd.Amount/100), 'f', 2, 64)
			hypha = strconv.FormatFloat(float64(payouts[index].Hypha.Amount/100), 'f', 2, 64)
			hvoice = strconv.FormatFloat(float64(payouts[index].Hvoice.Amount/100), 'f', 2, 64)
			seedsEscrow = strconv.FormatFloat(float64(payouts[index].SeedsEscrow.Amount/10000), 'f', 2, 64)
			seedsLiquid = strconv.FormatFloat(float64(payouts[index].SeedsLiquid.Amount/10000), 'f', 2, 64)
		} else {
			husd = payouts[index].Husd.String()
			hypha = payouts[index].Hypha.String()
			hvoice = payouts[index].Hvoice.String()
			seedsEscrow = payouts[index].SeedsEscrow.String()
			seedsLiquid = payouts[index].SeedsLiquid.String()
		}

		r := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: strconv.Itoa(int(payouts[index].ID))},
			{Align: simpletable.AlignRight, Text: string(payouts[index].Receiver)},
			{Align: simpletable.AlignLeft, Text: payouts[index].Title},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(payouts[index].DeferredPay*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: husd},
			{Align: simpletable.AlignRight, Text: hypha},
			{Align: simpletable.AlignRight, Text: hvoice},
			{Align: simpletable.AlignRight, Text: seedsEscrow},
			{Align: simpletable.AlignRight, Text: seedsLiquid},

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

	return table
}
