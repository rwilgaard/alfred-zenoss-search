package cmd

import (
	"fmt"
	"time"

	"github.com/rwilgaard/alfred-zenoss-search/src/internal/util"
	"github.com/rwilgaard/alfred-zenoss-search/src/internal/zenoss"
	"github.com/rwilgaard/go-alfredutils/alfredutils"
	"github.com/spf13/cobra"
)

var eventsCmd = &cobra.Command{
	Use:   "events [uid]",
	Short: "show events for a device",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := alfredutils.HandleAuthentication(wf, keychainAccount); !ok {
			return nil
		}

		uid := args[0]
		url := cfg.Endpoint

		client, err := setupAPI(url)
		if err != nil {
			return err
		}

		events, err := client.GetEvents(cmd.Context(), uid)
		if err != nil {
			return err
		}

		wf.NewItem("Go back").
			Arg("go_back").
			Icon(util.GetIcon("go_back")).
			Valid(true)

		for _, e := range events {
			var lastSeen string
			switch v := e.LastTime.(type) {
			case string:
				lastSeen = v
			case float64:
				lastSeen = time.Unix(int64(v), 0).Format("02-01-2006 15:04")
			}

			i := wf.NewItem(e.Summary).
				Arg("open").
				Subtitle(fmt.Sprintf("Count: %d  •  Last Seen: %s", e.Count, lastSeen)).
				Var("url", fmt.Sprintf("%s/zport/dmd/Events/viewDetail?evid=%s", url, e.EvID)).
				Icon(util.GetIcon(zenoss.SeverityCodes[e.Severity])).
				Valid(true)

			i.NewModifier("opt").
				Subtitle(fmt.Sprintf("Component: %s", e.Component.Text))
		}

		alfredutils.HandleFeedback(wf)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(eventsCmd)
}
