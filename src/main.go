package main

import (
    "log"
    "os"
    "os/exec"

    aw "github.com/deanishe/awgo"
    "github.com/deanishe/awgo/update"
    "github.com/rwilgaard/go-zenoss"
)

type workflowConfig struct {
    URL      string `env:"zenoss_url"`
    Username string `env:"username"`
    Password string
}

const (
    repo            = "rwilgaard/alfred-zenoss-search"
    updateJobName   = "checkForUpdates"
    keychainAccount = "alfred-zenoss-search"
)

var (
    wf  *aw.Workflow
    cfg *workflowConfig
)

func init() {
    wf = aw.New(
        update.GitHub(repo),
        aw.AddMagic(magicAuth{wf}),
    )
}

func run() {
    if err := cli.Parse(wf.Args()); err != nil {
        wf.FatalError(err)
    }
    opts.Query = cli.Arg(0)

    if opts.Update {
        wf.Configure(aw.TextErrors(true))
        log.Println("Checking for updates...")
        if err := wf.CheckForUpdate(); err != nil {
            wf.FatalError(err)
        }
        return
    }

    if wf.UpdateCheckDue() && !wf.IsRunning(updateJobName) {
        log.Println("Running update check in background...")
        cmd := exec.Command(os.Args[0], "-update")
        if err := wf.RunInBackground(updateJobName, cmd); err != nil {
            log.Printf("Error starting update check: %s", err)
        }
    }

    if wf.UpdateAvailable() {
        wf.Configure(aw.SuppressUIDs(true))
        wf.NewItem("Update Available!").
            Subtitle("Press ⏎ to install").
            Autocomplete("workflow:update").
            Valid(false).
            Icon(aw.IconInfo)
    }

    cfg = &workflowConfig{}
    if err := wf.Config.To(cfg); err != nil {
        panic(err)
    }

    if opts.Auth {
        runAuth()
    }

    password, err := wf.Keychain.Get(keychainAccount)
    if err != nil {
        wf.NewItem("You're not logged in.").
            Subtitle("Press ⏎ to authenticate").
            Icon(aw.IconInfo).
            Arg("auth").
            Valid(true)
        wf.SendFeedback()
        return
    }

    cfg.Password = password

    api, err := zenoss.NewAPI(cfg.URL, cfg.Username, cfg.Password)
    if err != nil {
        wf.FatalError(err)
    }

    runSearch(api)

    if wf.IsEmpty() {
        wf.NewItem("No results found...").
            Subtitle("Try a different query?").
            Icon(aw.IconInfo)
    }
    wf.SendFeedback()
}

func main() {
    wf.Run(run)
}
