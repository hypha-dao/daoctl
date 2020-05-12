package views

import (
  "fmt"
	"github.com/alexeyco/simpletable"
	"github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
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
func TreasuryTable(treasurys []models.Treasury) *simpletable.Table {

	table := simpletable.New()
	table.Header = treasuryHeader()

	balanceTotal, _ := eos.NewAssetFromString("0.00 HUSD")

	for index := range treasurys {

		fmt.Println("Balance	: ", treasurys[index].Balance.String())
		fmt.Println("Balance Total	: ", balanceTotal.String())
		balanceTotal = balanceTotal.Add(treasurys[index].Balance)

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: string(treasurys[index].TokenHolder)},
			{Align: simpletable.AlignRight, Text: FormatAsset(&treasurys[index].Balance)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "Total"},
			{Align: simpletable.AlignRight, Text: FormatAsset(&balanceTotal)},
		},
	}

	return table
}
