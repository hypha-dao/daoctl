package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hypha-dao/daoctl/models"
	"github.com/spf13/cobra"
)

var treasuryGetRequestCmd = &cobra.Command{
	Use:   "request <redemption_id>",
	Short: "view the details of a specific redemption",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		requestID, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			fmt.Println("Parse error: Request ID must be a positive integer (uint64)")
			return
		}

		request, err := models.LoadRequestByID(ctx, getAPI(), requestID)
		if err != nil {
			fmt.Println("Request ID not found")
			return
		}

		jsonDoc, _ := json.MarshalIndent(request, "", "  ")

		fmt.Println("\nRequest Details")
		fmt.Println()
		fmt.Println(string(jsonDoc))
		fmt.Println()
	},
}

func init() {
	treasuryGetCmd.AddCommand(treasuryGetRequestCmd)
}
