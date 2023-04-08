package main

import "flag"

var (
    opts = &options{}
    cli  = flag.NewFlagSet("alfred-zenoss-search", flag.ContinueOnError)
)

type options struct {
    // Arguments
    Query string

    // Commands
    Update   bool
    Events   bool
    Auth     bool
}

func init() {
    cli.BoolVar(&opts.Update, "update", false, "check for updates")
    cli.BoolVar(&opts.Events, "events", false, "get events")
    cli.BoolVar(&opts.Auth, "auth", false, "authenticate")
}
