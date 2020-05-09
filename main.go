package main

import (
	"fmt"

	"github.com/eoscanada/eos-go"
	"golang.org/x/net/context"
)

type Period struct {
	PeriodID  uint64             `json:"period_id"`
	StartTime eos.BlockTimestamp `json:"start_date"`
	EndTime   eos.BlockTimestamp `json:"end_date"`
	Phase     string             `json:"phase"`
}

func main() {
	api := eos.New("https://api.telos.kitchen")
	ctx := context.Background()

	infoResp, _ := api.GetInfo(ctx)
	fmt.Println("Info Resp: ", infoResp)

	var req eos.GetTableRowsRequest
	req.Code = "dao.hypha"
	req.Scope = "dao.hypha"
	req.Table = "periods"
	req.JSON = true
	req.Limit = 100

	var out []*Period
	var out2 []*Period
	periods, _ := api.GetTableRows(ctx, req)
	periods.BinaryToStructs(&out)
	periods.JSONToStructs(&out2)

	var assignments []Object
	var assignmentRequest eos.GetTableRowsRequest
	assignmentRequest.Code = "dao.hypha"
	assignmentRequest.Scope = "assignment"
	assignmentRequest.Table = "objects"
	assignmentRequest.Limit = 1
	assignmentRequest.JSON = true
	assignmentsResponse, _ := api.GetTableRows(ctx, assignmentRequest)

	assignmentsResponse.JSONToStructs(&assignments)
	fmt.Println("Info: ", out)

	daoObject := ToDAOObject(assignments[0])
	fmt.Println(daoObject)

	ass := ToAssignment(daoObject)
	fmt.Println(ass)
	// var o Object
	// o.Assets = make(map[string]eos.Asset)
	// o.Strings = make(map[string]string)
	// o.Names = make(map[string]eos.Name)
	// o.Assets["eos_asset"] = eos.NewEOSAsset(10000)
	// o.Strings["test_string_key"] = "here is a new string"

	// j, _ := json.Marshal(o)
	// fmt.Println("j: ", string(j))
}
