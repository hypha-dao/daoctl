package views

import (
	"math/big"
	"strconv"

	"github.com/hypha-dao/daoctl/util"

	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
)

func roleHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Hash"},
			{Align: simpletable.AlignCenter, Text: "Title"},
			{Align: simpletable.AlignCenter, Text: "Owner"},
			{Align: simpletable.AlignCenter, Text: "Min Time %"},
			{Align: simpletable.AlignCenter, Text: "Min Def %"},
			{Align: simpletable.AlignCenter, Text: "FTE Cap"},
			{Align: simpletable.AlignCenter, Text: "Annual USD"},
			{Align: simpletable.AlignCenter, Text: "Extended"},
			{Align: simpletable.AlignCenter, Text: "Ballot"},
		},
	}
}

// RoleTable returns a string representing an output table for a Role array
func RoleTable(roles []models.Role) *simpletable.Table {

	table := simpletable.New()
	table.Header = roleHeader()

	usdFteTotal, _ := eos.NewAssetFromString("0.00 USD")

	for index := range roles {

		usdFte := util.AssetMult(roles[index].AnnualUSDSalary, big.NewFloat(roles[index].FullTimeCapacity))
		usdFteTotal = usdFteTotal.Add(usdFte)

		r := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: roles[index].Hash.String()[:5]},
			{Align: simpletable.AlignLeft, Text: string(roles[index].Title)},
			{Align: simpletable.AlignRight, Text: string(roles[index].Owner)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(roles[index].MinTime*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(roles[index].MinDeferred*100, 'f', -1, 64)},
			{Align: simpletable.AlignRight, Text: strconv.FormatFloat(roles[index].FullTimeCapacity, 'f', 1, 64)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&roles[index].AnnualUSDSalary, 0)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&usdFte, 0)},
			{Align: simpletable.AlignRight, Text: string(roles[index].BallotName)[10:]},
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{}, {}, {}, {}, {}, {},
			{Align: simpletable.AlignRight, Text: "Subtotal"},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&usdFteTotal, 0)},
			{},
		},
	}

	return table
}
