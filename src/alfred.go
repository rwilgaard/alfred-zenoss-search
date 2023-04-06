package main

import (
    "fmt"

    aw "github.com/deanishe/awgo"
    "github.com/ncruces/zenity"
    "github.com/rwilgaard/go-zenoss"
)

type magicAuth struct {
    wf *aw.Workflow
}

func (a magicAuth) Keyword() string     { return "clearauth" }
func (a magicAuth) Description() string { return "Clear credentials for Confluence." }
func (a magicAuth) RunText() string     { return "Credentials cleared!" }
func (a magicAuth) Run() error          { return clearAuth() }

func runAuth() {
    _, pwd, err := zenity.Password(
        zenity.Title(fmt.Sprintf("Enter password for %s", cfg.Username)),
    )
    if err != nil {
        wf.FatalError(err)
    }
    if err := wf.Keychain.Set(keychainAccount, pwd); err != nil {
        wf.FatalError(err)
    }
}

func runSearch(api *zenoss.API) {
    devices, err := getDevices(api, opts.Query)
    if err != nil {
        wf.FatalError(err)
    }

    for _, d := range devices {
        cCount := d.Events.Critical.Count
        eCount := d.Events.Error.Count
        wCount := d.Events.Warning.Count
        eventCount := cCount + eCount + wCount
        wf.NewItem(d.Name).
            Arg("open").
            Subtitle(fmt.Sprintf("%s  •  Events: %d", prodStates[d.ProductionState], eventCount)).
            Var("device_url", cfg.URL+d.UID).
            Valid(true)
    }
}

func clearAuth() error {
    if err := wf.Keychain.Delete(keychainAccount); err != nil {
        return err
    }
    return nil
}
