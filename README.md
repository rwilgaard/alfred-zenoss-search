# alfred-zenoss-search

Search Zenoss devices and events from Alfred.

## Requirements

- [Alfred](https://www.alfredapp.com/) with Powerpack
- Zenoss instance with API access

## Installation

1. Download the latest `.alfredworkflow` from [Releases](https://github.com/rwilgaard/alfred-zenoss-search/releases)
2. Double-click to import into Alfred
3. If running macOS Catalina or later, add Alfred to security exceptions for unsigned software — see [this guide](https://github.com/deanishe/awgo/wiki/Catalina)
4. Set `zenoss_url` and `username` in workflow configuration
5. Trigger `zen` → press ⏎ on "You're not logged in" → enter your Zenoss password

## Usage

### `zen [query]` — search devices

| Action | Result |
|--------|--------|
| ⏎ | Open device in browser |
| ⌘⏎ | Show events for device |

### Events sub-menu

Shows open warning/error/critical events for the selected device. Press ⏎ on "Go back" to return to device list.

| Action | Result |
|--------|--------|
| ⏎ | Open event in browser |
| ⌥ | Show component name in subtitle |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `zenoss_url` | — | Zenoss instance URL(s), comma-separated for multiple instances |
| `username` | — | Zenoss username |

## Development

```sh
make build          # build arch-specific binaries
make package-alfred # build + zip into .alfredworkflow
make fmt            # format with gofumpt
make release VERSION=x.y.z  # bump version, package, tag, push
```

Requires Go 1.24+ and `golangci-lint`.
