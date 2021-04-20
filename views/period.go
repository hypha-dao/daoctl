package views

import (
	"strconv"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/hypha-dao/daoctl/models"
)

func periodHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "Hash"},
			{Align: simpletable.AlignCenter, Text: "Label"},
			{Align: simpletable.AlignCenter, Text: "Start Time"},
			{Align: simpletable.AlignCenter, Text: "Duration Days"},
			{Align: simpletable.AlignCenter, Text: "Duration"},
			{Align: simpletable.AlignCenter, Text: "Next"},
		},
	}
}

// PeriodTable returns a string representing an output table for a Role array
func PeriodTable(start models.Period) *simpletable.Table {

	table := simpletable.New()
	table.Header = periodHeader()
	period := start

	for {

		var duration time.Duration
		var durationDaysStr, durationStr, nextStr string

		if period.Next != nil {
			duration = period.Next.StartTime.Sub(period.StartTime)
			durationStr = duration.String()
			durationDaysStr = strconv.FormatFloat(duration.Hours()/24, 'f', 2, 64)
			nextStr = period.Next.Document.Hash.String()[:5]
		} else {
			durationStr = "n/a"
			durationDaysStr = "n/a"
			nextStr = "n/a"
		}

		r := []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: strconv.Itoa(int(period.Document.ID))},
			{Align: simpletable.AlignCenter, Text: period.Document.Hash.String()[:5]},
			{Align: simpletable.AlignLeft, Text: string(period.Label)},
			{Align: simpletable.AlignRight, Text: period.StartTime.Format("2006 Jan 02 15:04:05")},
			{Align: simpletable.AlignRight, Text: durationDaysStr},
			{Align: simpletable.AlignRight, Text: durationStr},
			{Align: simpletable.AlignCenter, Text: nextStr},
		}

		table.Body.Cells = append(table.Body.Cells, r)

		if period.Next == nil {
			return table
		}
		period = *period.Next
	}
}
