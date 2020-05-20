package cmd

import (
	"context"
	"fmt"
	"strconv"

	eos "github.com/eoscanada/eos-go"
	"github.com/hypha-dao/daoctl/models"
	"github.com/hypha-dao/daoctl/views"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getTreasuryCmd = &cobra.Command{
	Use:   "treasury",
	Short: "retrieve multi-chain balance information for the treasury",
	// Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		api := eos.New(viper.GetString("EosioEndpoint"))

		addlBalance, err := eos.NewAssetFromString(viper.GetString("get-treasury-cmd-addl-balance"))
		if err != nil {
			fmt.Println("Unable to read addl-balance parameter, using 0.00 HUSD")
			addlBalance, _ = eos.NewAssetFromString("0.00 HUSD")
		}

		accountName := toAccount(viper.GetString("Treasury.Contract"), "treasury contract account name")
		account, err := api.GetAccount(context.Background(), accountName)
		printTreasurers(account)

		// config := models.LoadTreasConfig(context.Background(), api)
		// fmt.Println(config)

		treasury := models.LoadTreasury(api, viper.GetString("Treasury.TokenContract"), viper.GetString("Treasury.Symbol"))

		fmt.Println()
		treasuryConfig := []string{
			fmt.Sprintf("Redemption Symbol|%v", *treasury.Config.RedemptionSymbol),
			fmt.Sprintf("Redemption Token Contract|%v", *treasury.Config.RedemptionTokenContract),
			fmt.Sprintf("Approval Threshold|%v", *treasury.Config.Threshold),
			fmt.Sprintf("Last Updated|%v", treasury.Config.UpdatedDate.Time.Format("2006 Jan 02 15:04:05")),
		}
		fmt.Println(columnize.SimpleFormat(treasuryConfig))

		treasuryTable, circulatingBalance := views.TreasuryTable(treasury.TreasuryHolders)
		fmt.Println("\n" + treasuryTable.String() + "\n\n")

		totalAssets := treasury.EthUSDTBalance.Add(addlBalance)
		usd := float64(totalAssets.Amount)
		circulating := float64(circulatingBalance.Amount)
		coverage := float64(usd / circulating * 100)
		netTreasury := totalAssets.Sub(circulatingBalance)
		output := []string{
			fmt.Sprintf("Awaiting burning|%v", views.FormatAsset(&treasury.BankBalance)),
			fmt.Sprintf("Circulating|%v", views.FormatAsset(&circulatingBalance)),
			fmt.Sprintf("Eth USDT balance|%v", views.FormatAsset(&treasury.EthUSDTBalance)),
			fmt.Sprintf("Additional balances|%v", views.FormatAsset(&addlBalance)),
			fmt.Sprintf("Total balances|%v", views.FormatAsset(&totalAssets)),
			fmt.Sprintf("Net Treasury balance|%v", views.FormatAsset(&netTreasury)),
			string(fmt.Sprintf("Coverage|%v", strconv.FormatFloat(coverage, 'f', 3, 64)) + " %"),
		}

		fmt.Println(columnize.SimpleFormat(output))
	},
}

func init() {
	getCmd.AddCommand(getTreasuryCmd)
	getTreasuryCmd.Flags().StringP("addl-balance", "", "0.00 HUSD", "Value of other accounts (e.g. banking, BTC) to optionally add manually")
}

func printTreasurers(account *eos.AccountResp) {
	if account != nil {
		// dereference this so we can safely mutate it to accomodate uninitialized symbols
		act := *account
		cfg := &columnize.Config{
			NoTrim: true,
		}

		for _, s := range []string{
			formatPermissions(&act, cfg),
		} {
			fmt.Println(s)
			fmt.Println("")
		}
	}
}

func formatPermissions(account *eos.AccountResp, config *columnize.Config) string {
	output := formatNestedPermission([]string{"\nTreasurers:"}, account.Permissions, eos.PermissionName("owner"), "")
	return columnize.Format(output, config)
}

func formatNestedPermission(in []string, permissions []eos.Permission, showChildsOf eos.PermissionName, indent string) (out []string) {
	const indentPadding = "      "
	out = in
	for _, perm := range permissions {
		if perm.Parent != string(showChildsOf) {
			continue
		}
		permValues := []string{}
		for _, key := range perm.RequiredAuth.Keys {
			permValues = append(permValues, fmt.Sprintf("+%d %s", key.Weight, key.PublicKey))
		}
		for _, acct := range perm.RequiredAuth.Accounts {
			permValues = append(permValues, fmt.Sprintf("+%d %s@%s", acct.Weight, acct.Permission.Actor, acct.Permission.Permission))
		}
		for _, wait := range perm.RequiredAuth.Waits {
			permValues = append(permValues, fmt.Sprintf("+%d wait %d seconds", wait.Weight, wait.WaitSec))
		}
		for i, keyValue := range permValues {
			if i == 0 {
				out = append(out,
					fmt.Sprintf("     %s Required Points: %d|:|%s",
						indent,
						perm.RequiredAuth.Threshold,
						keyValue,
					),
				)
			} else {
				out = append(out,
					fmt.Sprintf("     ||%s",
						keyValue,
					),
				)
			}
		}
		out = formatNestedPermission(out, permissions, eos.PermissionName(perm.PermName), indent+indentPadding)
	}
	return out
}
