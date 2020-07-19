package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// used to construct the action data parameter
type attestActionParam struct {
	Treasurer    eos.AccountName   `json:"treasurer"`
	PaymentID    uint64            `json:"payment_id"`
	RedemptionID uint64            `json:"redemption_id"`
	Amount       eos.Asset         `json:"amount"`
	Notes        map[string]string `json:"notes"`
}

var treasuryAttestPaymentCmd = &cobra.Command{
	Use:   "attest [paymentID] [redemptionID] [amount]",
	Short: "treasurer only; attests to the validity/truth of a payment created by another treasurer",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		paymentID, err := strconv.ParseUint(args[0], 10, 64)
		if err != nil {
			panic("Cannot read payment ID to submit payment against; halting")
		}

		redemptionID, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			panic("Cannot read redemption ID to submit payment against; halting")
		}

		amount, err := eos.NewAssetFromString(args[2])
		if err != nil {
			panic("Unable to read the amount as an Asset object; halting")
		}

		notes := make(map[string]string)

		if len(viper.GetString("treasury-cmd-memo")) > 0 {
			notes["memo"] = viper.GetString("treasury-cmd-memo")
		}

		action := eos.Action{
			Account: eos.AN(viper.GetString("Treasury.Contract")),
			Name:    toActionName("attestpaymnt", "new payment action name"),
			Authorization: []eos.PermissionLevel{
				{Actor: eos.AN(viper.GetString("DAOUser")), Permission: eos.PN("active")},
			},
			ActionData: eos.NewActionData(attestActionParam{
				Treasurer:    eos.AN(viper.GetString("DAOUser")),
				PaymentID:    paymentID,
				RedemptionID: redemptionID,
				Amount:       amount,
				Notes:        notes,
			}),
		}
		act, _ := json.MarshalIndent(action, "", " ")
		fmt.Println(string(act))

		pushEOSCActions(ctx, getAPI(), &action)
	},
}

func init() {
	treasuryCmd.AddCommand(treasuryAttestPaymentCmd)
}
