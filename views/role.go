package views

import (
	"math/big"
	"strconv"

	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
)

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
			{Align: simpletable.AlignCenter, Text: "Extended"},
			{Align: simpletable.AlignCenter, Text: "Start Date"},
			{Align: simpletable.AlignCenter, Text: "End Date"},
		},
	}
}

func assetMult(a eos.Asset, coeff *big.Float) eos.Asset {
	//var amount big.Int
	var f big.Float
	f.SetInt(big.NewInt(int64(a.Amount)))
	amount, _ := f.Mul(&f, coeff).Int64() // big.NewInt(int64(a.Amount)).Mul(coeff)
	newAmount := eos.Int64(amount)
	// intObject := big.NewInt(amount)
	return eos.Asset{Amount: newAmount, Symbol: a.Symbol}
}

// RoleTable returns a string representing an output table for a Role array
func RoleTable(roles []models.Role) *simpletable.Table {

	table := simpletable.New()
	table.Header = roleHeader()

	usdFteTotal, _ := eos.NewAssetFromString("0.00 USD")

	for index := range roles {

		usdFte := assetMult(roles[index].AnnualUSDSalary, big.NewFloat(roles[index].FullTimeCapacity))
		usdFteTotal = usdFteTotal.Add(usdFte)

		r := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: strconv.Itoa(int(roles[index].ID))},
			{Align: simpletable.AlignLeft, Text: string(roles[index].Title)},
			{Align: simpletable.AlignRight, Text: string(roles[index].Owner)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(roles[index].MinTime*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(roles[index].MinDeferred*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(roles[index].FullTimeCapacity, 'f', 1, 64)},
			{Align: simpletable.AlignRight, Text: FormatAsset(&roles[index].AnnualUSDSalary)},
			{Align: simpletable.AlignRight, Text: FormatAsset(&usdFte)},
			{Align: simpletable.AlignRight, Text: roles[index].StartPeriod.StartTime.Time.Format("2006 Jan 02")},
			{Align: simpletable.AlignRight, Text: roles[index].EndPeriod.EndTime.Time.Format("2006 Jan 02")},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{}, {}, {}, {}, {}, {},
			{Align: simpletable.AlignRight, Text: "Subtotal"},
			{Align: simpletable.AlignRight, Text: FormatAsset(&usdFteTotal)},
			{}, {},
		},
	}

	return table
}
