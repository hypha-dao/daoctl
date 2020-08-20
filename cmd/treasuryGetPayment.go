package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hypha-dao/daoctl/models"
	"github.com/spf13/cobra"
)

var treasuryGetPaymentCmd = &cobra.Command{
	Use:   "payment <payment_id>",
	Short: "view the details of a specific treasury payment",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		paymentID, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			fmt.Println("Parse error: Payment ID must be a positive integer (uint64)")
			return
		}

		payment, err := models.LoadPaymentByID(ctx, getAPI(), paymentID)
		if err != nil {
			fmt.Println("Payment ID not found")
			return
		}

		jsonDoc, _ := json.MarshalIndent(payment, "", "  ")

		fmt.Println("\nPayment Details")
		fmt.Println()
		fmt.Println(string(jsonDoc))
		fmt.Println()
	},
}

func init() {
	treasuryGetCmd.AddCommand(treasuryGetPaymentCmd)
}
