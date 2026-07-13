package cmd

import (
	"fmt"
	"strings"
	"sync"

	aw "github.com/deanishe/awgo"
	"github.com/rwilgaard/alfred-zenoss-search/src/internal/util"
	"github.com/rwilgaard/alfred-zenoss-search/src/internal/zenoss"
	"github.com/rwilgaard/go-alfredutils/alfredutils"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "search Zenoss devices",
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if ok := alfredutils.HandleAuthentication(wf, keychainAccount); !ok {
			return
		}

		var query string
		if len(args) > 0 {
			query = args[0]
		}

		urlList := strings.Split(cfg.URL, ",")
		errCh := make(chan error, len(urlList))
		var mu sync.Mutex
		var wg sync.WaitGroup

		for _, url := range urlList {
			url = strings.TrimSpace(url)
			client, err := setupAPI(url)
			if err != nil {
				wf.FatalError(err)
				return
			}
			wg.Add(1)
			go func(client *zenoss.Client, url string) {
				defer wg.Done()
				devices, err := client.GetDevices(cmd.Context(), query)
				if err != nil {
					errCh <- err
					return
				}
				mu.Lock()
				defer mu.Unlock()
				for _, d := range devices {
					eventCount := d.Events.Critical.Count + d.Events.Error.Count + d.Events.Warning.Count
					i := wf.NewItem(d.Name).
						Arg("open").
						Subtitle(fmt.Sprintf("%s  •  Events: %d", zenoss.ProdStates[d.ProductionState], eventCount)).
						Var("uid", d.UID).
						Var("url", url+d.UID).
						Var("endpoint", url).
						Var("query", query).
						Icon(util.GetIcon(d.OsManufacturer["name"])).
						Valid(true)

					i.NewModifier(aw.ModCmd).
						Arg("events").
						Subtitle("Show events")
				}
			}(client, url)
		}
		wg.Wait()
		close(errCh)

		if err := <-errCh; err != nil {
			wf.FatalError(err)
			return
		}

		alfredutils.HandleFeedback(wf)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
