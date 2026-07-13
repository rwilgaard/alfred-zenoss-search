package cmd

import (
	"fmt"
	"strings"

	"github.com/ncruces/zenity"
	"github.com/rwilgaard/alfred-zenoss-search/src/internal/zenoss"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "authenticate with Zenoss",
	RunE: func(cmd *cobra.Command, _ []string) error {
		_, pwd, err := zenity.Password(
			zenity.Title(fmt.Sprintf("Enter password for %s", cfg.Username)),
		)
		if err != nil {
			return err
		}
		pwd = strings.TrimSpace(pwd)

		url := strings.TrimSpace(strings.SplitN(cfg.URL, ",", 2)[0])
		client, err := zenoss.NewClient(url, cfg.Username, pwd)
		if err != nil {
			return err
		}
		if err := client.TestAuthentication(cmd.Context()); err != nil {
			zerr := zenity.Error(
				fmt.Sprintf("Authentication failed: %s", err),
				zenity.ErrorIcon,
			)
			if zerr != nil {
				return err
			}
			return err
		}

		if err := wf.Keychain.Set(keychainAccount, pwd); err != nil {
			return err
		}
		fmt.Println("Successfully authenticated")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}
