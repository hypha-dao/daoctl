package views

import (
	"strconv"

	"github.com/alexeyco/simpletable"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/daoctl/util"
)

func docHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "Owner"},
			{Align: simpletable.AlignCenter, Text: "Type"},
			{Align: simpletable.AlignCenter, Text: "Created Date"},
			{Align: simpletable.AlignCenter, Text: "Updated Date"},
			{Align: simpletable.AlignCenter, Text: "Notes"},
		},
	}
}

// DocTable is a simpleTable.Table object with documents
func DocTable(docs []models.Document) *simpletable.Table {

	table := simpletable.New()
	table.Header = docHeader()

	for _, doc := range docs {

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: strconv.Itoa(int(doc.ID))},
			{Align: simpletable.AlignRight, Text: string(doc.Names["owner"])},
			{Align: simpletable.AlignRight, Text: string(doc.Names["type"])},
			{Align: simpletable.AlignRight, Text: doc.CreatedDate.Time.Format("2006 Jan 02")},
			{Align: simpletable.AlignRight, Text: doc.UpdatedDate.Time.Format("2006 Jan 02")},
			{Align: simpletable.AlignRight, Text: util.Snip(&doc.Strings)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	return table
}
