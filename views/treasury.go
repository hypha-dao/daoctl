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
			{Align: simpletable.AlignCenter, Text: "Requested Redemptions"},
		},
	}
}

// TreasuryTable returns a string representing an output table for a Treasury array
func TreasuryTable(members map[eos.Name]models.Balance) (*simpletable.Table, eos.Asset) {

	table := simpletable.New()
	table.Header = treasuryHeader()

	balanceTotal, _ := eos.NewAssetFromString("0.00 HUSD")
	redemptionsTotal, _ := eos.NewAssetFromString("0.00 HUSD")

	for member, treasuryBalance := range members {

		//if members[index].Balance.Amount > 0 {

		balanceTotal = balanceTotal.Add(treasuryBalance.Balance)
		redemptionsTotal = redemptionsTotal.Add(treasuryBalance.RequestedRedemptions)

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: string(member)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&treasuryBalance.Balance, 2)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&treasuryBalance.RequestedRedemptions, 2)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
		//}
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "Total"},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&balanceTotal, 2)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&redemptionsTotal, 2)},
		},
	}

	return table, balanceTotal
}
