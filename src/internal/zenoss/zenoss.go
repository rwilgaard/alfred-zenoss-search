package zenoss

import (
	"fmt"

	goz "github.com/rwilgaard/go-zenoss"
)

type Client struct {
	c   *goz.API
	URL string
}

var ProdStates = map[int]string{
	-1:   "Decommissioned",
	100:  "Disabled",
	300:  "Maintenance",
	400:  "Test",
	500:  "Pre-Production",
	600:  "Migrated",
	900:  "Customer-Only",
	1000: "Production",
}

var SeverityCodes = map[int]string{
	0: "Clear",
	1: "Debug",
	2: "Info",
	3: "Warning",
	4: "Error",
	5: "Critical",
}

func NewClient(url, username, password string) (*Client, error) {
	api, err := goz.NewAPI(url, username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	return &Client{c: api, URL: url}, nil
}

func (cl *Client) TestAuthentication() error {
	q := goz.GetDevicesQuery{Limit: 1, Keys: []string{"name"}}
	_, res, err := cl.c.GetDevices(q)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("authentication failed: HTTP %d", res.StatusCode)
	}
	return nil
}

func (cl *Client) GetDevices(query string) ([]goz.Device, error) {
	q := goz.GetDevicesQuery{
		Limit: 10,
		Keys: []string{
			"name",
			"uid",
			"productionState",
			"osManufacturer",
			"events",
		},
		Params: map[string]any{
			"name": query,
		},
	}

	d, res, err := cl.c.GetDevices(q)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error fetching devices. statuscode: %d", res.StatusCode)
	}

	return d.Result.Devices, nil
}

func (cl *Client) GetEvents(uid string) ([]goz.Event, error) {
	q := goz.QueryEventsQuery{
		UID:   uid,
		Limit: 25,
		Keys: []string{
			"summary",
			"severity",
			"component",
			"count",
			"evid",
			"lastTime",
		},
		Params: map[string]any{
			"eventState": []int{0, 1},
			"severity":   []int{3, 4, 5},
		},
	}

	e, res, err := cl.c.QueryEvents(q)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("error querying events. statuscode: %d", res.StatusCode)
	}

	return e.Result.Events, nil
}
