package models

import (
	"github.com/alexeyco/simpletable"
)

// TableToData converts a simpletable.Table object to a 2-dimensional array, which can be used for exporting to CSV
func TableToData(table *simpletable.Table) [][]string {

	data := make([][]string, len(table.Header.Cells)+len(table.Body.Cells)+len(table.Footer.Cells))

	data[0] = make([]string, len(table.Header.Cells))
	for index, element := range table.Header.Cells {
		data[0][index] = element.Text
	}

	for rowIndex := range table.Body.Cells {
		data[rowIndex+1] = make([]string, len(table.Body.Cells))
		for columnIndex := range table.Body.Cells[rowIndex] {
			data[rowIndex+1][columnIndex] = table.Body.Cells[rowIndex][columnIndex].Text
		}
	}
	return data
}
