package main

import (
	"context"

	eos "github.com/eoscanada/eos-go"
)

func main() {
	api := eos.New("https://api.telos.kitchen")
	// api := eos.New("https://")

	// infoResp, _ := api.GetInfo(ctx)
	// infoRespStr, _ := json.MarshalIndent(infoResp, "", "  ")
	periods := LoadPeriods(api)
	PrintAssignments(context.Background(), api, periods, true)
	PrintRoles(context.Background(), api, periods, true)
	PrintPayouts(context.Background(), api, periods, true)
}
