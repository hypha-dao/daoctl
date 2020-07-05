package views

import (
	"strconv"

	"github.com/alexeyco/simpletable"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/daoctl/util"
)

func paymentHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "Creator"},
			{Align: simpletable.AlignCenter, Text: "Request ID"},
			{Align: simpletable.AlignCenter, Text: "Paid"},
			{Align: simpletable.AlignCenter, Text: "Created Date"},
			{Align: simpletable.AlignCenter, Text: "Confirmed Date"},
			{Align: simpletable.AlignCenter, Text: "Attestations"},
			{Align: simpletable.AlignCenter, Text: "Notes"},
		},
	}
}

// PaymentTable is a simpleTable.Table object with payouts
func PaymentTable(payments []models.Payment) *simpletable.Table {

	table := simpletable.New()
	table.Header = paymentHeader()

	for _, payment := range payments {

		attestationStr := ""

		if len(payment.Attestations) > 0 {
			for index, att := range payment.Attestations {
				if index == len(payment.Attestations)-1 {
					attestationStr += string(att.Key)
				} else {
					attestationStr += string(att.Key) + ", "
				}
			}
		}

		confirmedDateStr := payment.ConfirmedDate.Time.Format("2006 Jan 02")
		if payment.ConfirmedDate.Time.Before(payment.CreatedDate.Time) {
			confirmedDateStr = "unconfirmed"
		}

		r := []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: strconv.Itoa(int(payment.ID))},
			{Align: simpletable.AlignRight, Text: string(payment.Creator)},
			{Align: simpletable.AlignRight, Text: strconv.Itoa(int(payment.RequestID))},
			{Align: simpletable.AlignRight, Text: util.FormatAsset(&payment.Amount, 2)},
			{Align: simpletable.AlignRight, Text: payment.CreatedDate.Time.Format("2006 Jan 02")},
			{Align: simpletable.AlignRight, Text: confirmedDateStr},
			{Align: simpletable.AlignRight, Text: attestationStr},
			{Align: simpletable.AlignRight, Text: util.Snip(payment.NotesMap)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}

	return table
}
