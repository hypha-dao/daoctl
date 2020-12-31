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
			{Align: simpletable.AlignCenter, Text: "Node Label"},
			{Align: simpletable.AlignCenter, Text: "Type"},
			{Align: simpletable.AlignCenter, Text: "Created Date"},
			{Align: simpletable.AlignCenter, Text: "Creator"},
			{Align: simpletable.AlignCenter, Text: "Hash"},
			// {Align: simpletable.AlignCenter, Text: "Content"},
		},
	}
}

func isSkipped(label string) bool {
	// skipLabels := []string{"period"}

	// for _, skipLabel := range skipLabels {
	// 	if label == skipLabel {
	// 		return true
	// 	}
	// }
	return false
}

func docString(d *docgraph.Document) string {

	var documentString string

	for contentGroupIndex, contentGroup := range d.ContentGroups {
		if contentGroupIndex > 0 {
			documentString += ","
		}
		documentString += "["
		for _, content := range contentGroup {

			if !isSkipped(content.Label) {
				documentString += "["
				documentString += content.Label
				documentString += "="

				if len(content.Value.String()) > 45 {
					documentString += content.Value.String()[:40] + "<snip>"
				} else {
					documentString += content.Value.String()
				}

				documentString += "]"
			}
		}
	}
	documentString += "]"
	return documentString
}

// DocTable is a simpleTable.Table object with documents
func DocTable(docs []docgraph.Document) *simpletable.Table {

	table := simpletable.New()
	table.Header = docHeader()

	for _, doc := range docs {

		// documentString := docString(&doc)
		// if len(documentString) > 300 {
		// 	documentString = documentString[:295] + "<snip>"
		// }

		typeLabel := "Unknown"
		documentType, _ := doc.GetContent("type")
		if documentType != nil {
			typeLabel = documentType.String()
		}

		nodeLabel := "Unknown"
		documentLabel, _ := doc.GetContent("node_label")
		if documentLabel != nil {
			nodeLabel = documentLabel.String()
		}

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: strconv.Itoa(int(doc.ID))},
			{Align: simpletable.AlignRight, Text: nodeLabel},
			{Align: simpletable.AlignRight, Text: typeLabel},
			{Align: simpletable.AlignRight, Text: doc.CreatedDate.Time.Format("2006 Jan 02 15:04:05")},
			{Align: simpletable.AlignRight, Text: string(doc.Creator)},
			{Align: simpletable.AlignRight, Text: doc.Hash.String()},
			// {Align: simpletable.AlignRight, Text: documentString},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	return table
}
