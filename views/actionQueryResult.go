package views

import (
	"math"
	"strconv"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/hypha-dao/daoctl/models"
)

func qrActionHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Timestamp"},
			{Align: simpletable.AlignCenter, Text: "Time Since"},
			{Align: simpletable.AlignCenter, Text: "Contract"},
			{Align: simpletable.AlignCenter, Text: "Action"},
			{Align: simpletable.AlignCenter, Text: "Data"},
			{Align: simpletable.AlignCenter, Text: "Trx ID"},
		},
	}
}

// ActionQueryResultTable is a simpleTable.Table object with payouts
func ActionQueryResultTable(qrActions []models.QrAction) *simpletable.Table {

	table := simpletable.New()
	table.Header = qrActionHeader()

	for index := range qrActions {

		timeSince := time.Since(qrActions[index].Timestamp)
		timeSinceStr := timeSince.String()
		if timeSince.Hours() < 1 {
			timeSinceStr = strconv.Itoa(RoundTime(timeSince.Seconds()/60)) + " mins"
		} else if timeSince.Hours() < 72 {
			timeSinceStr = strconv.Itoa(RoundTime(timeSince.Seconds()/3600)) + " hours"
		} else {
			timeSinceStr = strconv.Itoa(RoundTime(timeSince.Seconds()/86400)) + " days"
		}

		r := []*simpletable.Cell{
			{Align: simpletable.AlignLeft, Text: qrActions[index].Timestamp.Format("2006 Jan 02 15:04:05")},
			{Align: simpletable.AlignLeft, Text: timeSinceStr},
			{Align: simpletable.AlignLeft, Text: qrActions[index].ActionContract},
			{Align: simpletable.AlignLeft, Text: qrActions[index].ActionName},
			{Align: simpletable.AlignLeft, Text: qrActions[index].Data},
			{Align: simpletable.AlignLeft, Text: qrActions[index].TrxID},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}
	return table
}

func RoundTime(input float64) int {
	var result float64

	if input < 0 {
		result = math.Ceil(input - 0.5)
	} else {
		result = math.Floor(input + 0.5)
	}

	// only interested in integer, ignore fractional
	i, _ := math.Modf(result)

	return int(i)
}
