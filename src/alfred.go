package main

import (
    "fmt"
    "os"
    "sync"
    "time"

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

func getIcon(query string) *aw.Icon {
    iconPath := fmt.Sprintf("icons/%s.png", query)
    icon := aw.IconWorkflow

    if _, err := os.Stat(iconPath); err == nil {
        icon = &aw.Icon{Value: iconPath}
    }

    return icon
}

func runSearch(api *zenoss.API, url string, wg *sync.WaitGroup) {
    defer wg.Done()
    devices, err := getDevices(api, opts.Query)
    if err != nil {
        wf.FatalError(err)
    }

    for _, d := range devices {
        cCount := d.Events.Critical.Count
        eCount := d.Events.Error.Count
        wCount := d.Events.Warning.Count
        eventCount := cCount + eCount + wCount
        i := wf.NewItem(d.Name).
            Arg("open").
            Subtitle(fmt.Sprintf("%s  •  Events: %d", prodStates[d.ProductionState], eventCount)).
            Var("uid", d.UID).
            Var("url", url+d.UID).
            Var("endpoint", url).
            Var("query", opts.Query).
            Icon(getIcon(d.OsManufacturer["name"])).
            Valid(true)

        i.NewModifier("cmd").
            Arg("events").
            Subtitle("Show events")
    }
}

func runEvents(api *zenoss.API, url string) {
    events, err := getEvents(api, opts.Query)
    if err != nil {
        wf.FatalError(err)
    }

    wf.NewItem("Go back").
        Arg("go_back").
        Icon(getIcon("go_back")).
        Valid(true)

    for _, e := range events {
        var lastSeen string
        switch e.LastTime.(type) {
        case string:
            lastSeen = e.LastTime.(string)
        case float64:
            ts := int64(e.LastTime.(float64))
            lastSeen = time.Unix(ts, 0).Format("02-01-2006 15:04")
        }

        i := wf.NewItem(e.Summary).
            Arg("open").
            Subtitle(fmt.Sprintf("Count: %d  •  Last Seen: %s", e.Count, lastSeen)).
            Var("url", fmt.Sprintf("%s/zport/dmd/Events/viewDetail?evid=%s", url, e.EvID)).
            Icon(getIcon(severityCodes[e.Severity])).
            Valid(true)

        i.NewModifier("opt").
            Subtitle(fmt.Sprintf("Component: %s", e.Component.Text))
    }
}

func clearAuth() error {
    if err := wf.Keychain.Delete(keychainAccount); err != nil {
        return err
    }
    return nil
}
