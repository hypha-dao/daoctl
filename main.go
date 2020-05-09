package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/eoscanada/eos-go"
	"github.com/gcla/gowid"
	"github.com/gcla/gowid/examples"
	"github.com/gcla/gowid/widgets/fill"
	"github.com/gcla/gowid/widgets/table"
	"github.com/gdamore/tcell"

	// "github.com/ipfs/go-log"
	log "github.com/sirupsen/logrus"

	"golang.org/x/net/context"
)

type handler struct{}

func (h handler) UnhandledInput(app gowid.IApp, ev interface{}) bool {
	if evk, ok := ev.(*tcell.EventKey); ok {
		if evk.Key() == tcell.KeyCtrlC || evk.Rune() == 'q' || evk.Rune() == 'Q' {
			app.Quit()
			return true
		}
	}
	return false
}

func main() {
	api := eos.New("https://api.telos.kitchen")
	ctx := context.Background()

	infoResp, _ := api.GetInfo(ctx)
	infoRespStr, _ := json.MarshalIndent(infoResp, "", "  ")
	fmt.Println("Info Resp: ", string(infoRespStr))

	var req eos.GetTableRowsRequest
	req.Code = "dao.hypha"
	req.Scope = "dao.hypha"
	req.Table = "periods"
	req.JSON = true
	req.Limit = 100

	var out []*Period
	periods, _ := api.GetTableRows(ctx, req)
	periods.JSONToStructs(&out)
	fmt.Println("Info: ", out)

	var assignmentObjects []Object
	var assignmentRequest eos.GetTableRowsRequest
	assignmentRequest.Code = "dao.hypha"
	assignmentRequest.Scope = "assignment"
	assignmentRequest.Table = "objects"
	assignmentRequest.Limit = 100
	assignmentRequest.JSON = true
	assignmentsResponse, _ := api.GetTableRows(ctx, assignmentRequest)

	assignmentsResponse.JSONToStructs(&assignmentObjects)

	var assignments []Assignment
	for index := range assignmentObjects {
		assignments = append(assignments, ToAssignment(ToDAOObject(assignmentObjects[index])))
	}

	assString, _ := json.MarshalIndent(assignments, "", "  ")
	fmt.Println(string(assString))

	var data [][]string
	var headers []string
	headers = make([]string, 8)
	headers[0] = "ID"
	headers[1] = "Owner"
	headers[2] = "Assigned"
	headers[3] = "HUSD"
	headers[4] = "HYPHA"
	headers[5] = "HVOICE"
	headers[6] = "SEEDS Escrow"
	headers[7] = "SEEDS Liquid"

	data = make([][]string, len(assignments))
	for index := range assignments {
		data[index] = make([]string, 8)
		data[index][0] = strconv.Itoa(int(assignments[index].ID))
		data[index][1] = string(assignments[index].Owner)
		data[index][2] = string(assignments[index].Assigned)
		data[index][3] = assignments[index].HusdPerPhase.String()
		data[index][4] = assignments[index].HyphaPerPhase.String()
		data[index][5] = assignments[index].HvoicePerPhase.String()
		data[index][6] = assignments[index].SeedsEscrowPerPhase.String()
		data[index][7] = assignments[index].SeedsLiquidPerPhase.String()
	}

	model := table.NewSimpleModel(headers, data, table.SimpleOptions{
		Style: table.StyleOptions{
			VerticalSeparator:   fill.New('|'),
			HorizontalSeparator: fill.New('-'),
		},
	})
	table := table.New(model)

	palette := gowid.Palette{
		"green": gowid.MakePaletteEntry(gowid.ColorDarkGreen, gowid.ColorDefault),
		"red":   gowid.MakePaletteEntry(gowid.ColorRed, gowid.ColorDefault),
	}

	app, err := gowid.NewApp(gowid.AppArgs{
		View:    table,
		Palette: &palette,
		Log:     log.StandardLogger(),
	})
	examples.ExitOnErr(err)

	app.MainLoop(handler{})
}
