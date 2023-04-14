package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

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

    if opts.Events {
        url := os.Getenv("zenoss_url")
        api, err := zenoss.NewAPI(url, cfg.Username, cfg.Password)
        if err != nil {
            wf.FatalError(err)
        }
        runEvents(api, url)
        wf.SendFeedback()
        return
    }

    var wg sync.WaitGroup
    for _, url := range strings.Split(cfg.URL, ",") {
        api, err := zenoss.NewAPI(strings.TrimSpace(url), cfg.Username, cfg.Password)
        if err != nil {
            wf.FatalError(err)
        }

        wg.Add(1)
        go runSearch(api, url, &wg)
    }

    wg.Wait()

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
