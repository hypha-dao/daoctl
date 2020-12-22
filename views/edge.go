package views

import (
	"strconv"

	"github.com/alexeyco/simpletable"
	"github.com/hypha-dao/document-graph/docgraph"
)

func edgeHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			// {Align: simpletable.AlignCenter, Text: "Creator"},
			{Align: simpletable.AlignCenter, Text: "Created"},
			{Align: simpletable.AlignCenter, Text: "From"},
			{Align: simpletable.AlignCenter, Text: "Edge"},
			{Align: simpletable.AlignCenter, Text: "To"},
		},
	}
}

// EdgeTable is a simpleTable.Table object with documents
func EdgeTable(edges []docgraph.Edge, overLabelSelf, fromSelf bool) *simpletable.Table {

	table := simpletable.New()
	table.Header = edgeHeader()

	for _, edge := range edges {

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: strconv.Itoa(int(edge.ID))},
			{Align: simpletable.AlignRight, Text: edge.CreatedDate.Time.Format("02 Jan 2006 15:04:05")}}

		if overLabelSelf {
			if fromSelf {
				r = append(r, &simpletable.Cell{Align: simpletable.AlignRight, Text: "[ self ]"},
					&simpletable.Cell{Align: simpletable.AlignRight, Text: string(edge.EdgeName)},
					&simpletable.Cell{Align: simpletable.AlignRight, Text: edge.ToNode.String()})
			} else {
				r = append(r, &simpletable.Cell{Align: simpletable.AlignRight, Text: edge.FromNode.String()},
					&simpletable.Cell{Align: simpletable.AlignRight, Text: string(edge.EdgeName)},
					&simpletable.Cell{Align: simpletable.AlignRight, Text: "[ self ]"})
			}
		} else {
			r = append(r, &simpletable.Cell{Align: simpletable.AlignRight, Text: edge.FromNode.String()},
				&simpletable.Cell{Align: simpletable.AlignRight, Text: string(edge.EdgeName)},
				&simpletable.Cell{Align: simpletable.AlignRight, Text: edge.ToNode.String()})
		}

		table.Body.Cells = append(table.Body.Cells, r)
	}

	return table
}
