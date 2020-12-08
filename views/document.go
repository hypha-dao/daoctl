package views

import (
	"strconv"

	"github.com/alexeyco/simpletable"
	"github.com/hypha-dao/document-graph/docgraph"
)

func docHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "Hash"},
			{Align: simpletable.AlignCenter, Text: "Creator"},
			{Align: simpletable.AlignCenter, Text: "Created Date"},
			{Align: simpletable.AlignCenter, Text: "Content"},
		},
	}
}

func docString(d *docgraph.Document) string {

	var documentString string

	for contentGroupIndex, contentGroup := range d.ContentGroups {
		if contentGroupIndex > 0 {
			documentString += ","
		}
		documentString += "["
		for _, content := range contentGroup {
			documentString += "[label="
			documentString += content.Label
			documentString += ","
			documentString += content.Value.String()
			documentString += "]"
		}
		documentString += "]"
	}

	return documentString
}

// DocTable is a simpleTable.Table object with documents
func DocTable(docs []docgraph.Document) *simpletable.Table {

	table := simpletable.New()
	table.Header = docHeader()

	for _, doc := range docs {

		documentString := docString(&doc)
		if len(documentString) > 45 {
			documentString = documentString[:40] + "<snip>"
		}

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: strconv.Itoa(int(doc.ID))},
			{Align: simpletable.AlignRight, Text: doc.Hash.String()},
			{Align: simpletable.AlignRight, Text: string(doc.Creator)},
			{Align: simpletable.AlignRight, Text: doc.CreatedDate.Time.Format("2006 Jan 02")},
			{Align: simpletable.AlignRight, Text: documentString},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	return table
}
