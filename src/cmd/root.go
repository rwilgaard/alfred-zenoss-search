package cmd

import (
	"fmt"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
	"github.com/rwilgaard/alfred-zenoss-search/src/internal/zenoss"
	"github.com/rwilgaard/go-alfredutils/alfredutils"
	"github.com/spf13/cobra"
)

type workflowConfig struct {
	URL      string `env:"zenoss_url"`
	Username string `env:"username"`
	Endpoint string `env:"endpoint"`
}

const (
	repo            = "rwilgaard/alfred-zenoss-search"
	keychainAccount = "alfred-zenoss-search"
)

var (
	wf  *aw.Workflow
	cfg = &workflowConfig{}

	rootCmd = &cobra.Command{
		Use:           "alfred-zenoss-search",
		Short:         "Alfred workflow for searching Zenoss devices",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
)

func Execute() {
	wf.Run(run)
}

func run() {
	alfredutils.AddClearAuthMagic(wf, keychainAccount)

	if err := alfredutils.InitWorkflow(wf, cfg); err != nil {
		wf.FatalError(err)
	}

	if err := alfredutils.CheckForUpdates(wf); err != nil {
		wf.FatalError(err)
	}

	if err := rootCmd.Execute(); err != nil {
		wf.FatalError(err)
	}
}

func setupAPI(url string) (*zenoss.Client, error) {
	if url == "" {
		return nil, fmt.Errorf("zenoss_url is not configured in workflow settings")
	}
	if cfg.Username == "" {
		return nil, fmt.Errorf("username is not configured in workflow settings")
	}
	password, err := wf.Keychain.Get(keychainAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to get password from keychain: %w", err)
	}
	return zenoss.NewClient(url, cfg.Username, password)
}

func init() {
	wf = aw.New(update.GitHub(repo))
}
