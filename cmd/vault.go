// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"github.com/spf13/cobra"
)

// vaultCmd represents the vault command
var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "The vault is a secure key store (wallet). Your key is stored encrypted by the passphrase.",
	Long:  "The vault is a secure key store (wallet). Your key is stored encrypted by the passphrase.",
}

func init() {
	RootCmd.AddCommand(vaultCmd)
}
