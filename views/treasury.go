package views

import (
	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
  "github.com/hypha-dao/daoctl/util"
)

func treasuryHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Token Holder"},
			{Align: simpletable.AlignCenter, Text: "Balance"},
		},
	}
}

// TreasuryTable returns a string representing an output table for a Treasury array
func TreasuryTable(treasurys []models.TreasuryHolder) (*simpletable.Table, eos.Asset) {

	table := simpletable.New()
	table.Header = treasuryHeader()

	balanceTotal, _ := eos.NewAssetFromString("0.00 HUSD")

	for index := range treasurys {

		if treasurys[index].Balance.Amount > 0 {

			balanceTotal = balanceTotal.Add(treasurys[index].Balance)

			r := []*simpletable.Cell{
				{Align: simpletable.AlignRight, Text: string(treasurys[index].TokenHolder)},
				{Align: simpletable.AlignRight, Text: util.FormatAsset(&treasurys[index].Balance, 2)},
			}
			table.Body.Cells = append(table.Body.Cells, r)
		}
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "Total"},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&balanceTotal, 2)},
		},
	}

	return table, balanceTotal
}
