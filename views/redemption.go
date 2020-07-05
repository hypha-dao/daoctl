package views

import (
	"strconv"

	"github.com/alexeyco/simpletable"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/daoctl/util"
)

func requestHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "Requestor"},
			{Align: simpletable.AlignCenter, Text: "Requested"},
			{Align: simpletable.AlignCenter, Text: "Paid"},
			{Align: simpletable.AlignCenter, Text: "Requested Date"},
			{Align: simpletable.AlignCenter, Text: "Updated Date"},
			{Align: simpletable.AlignCenter, Text: "Notes"},
		},
	}
}

// RequestTable is a simpleTable.Table object with payouts
func RequestTable(requests []models.RedemptionRequest) *simpletable.Table {

	table := simpletable.New()
	table.Header = requestHeader()

	for _, request := range requests {

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: strconv.Itoa(int(request.ID))},
			{Align: simpletable.AlignRight, Text: string(request.Requestor)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&request.Requested, 2)},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&request.Paid, 2)},
			{Align: simpletable.AlignRight, Text: request.RequestedDate.Time.Format("2006 Jan 02")},
			{Align: simpletable.AlignRight, Text: request.UpdatedDate.Time.Format("2006 Jan 02")},
			{Align: simpletable.AlignRight, Text: util.Snip(request.NotesMap)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	return table
}
