package main

import (
	"fmt"

	"github.com/rwilgaard/go-zenoss"
)

var (
    prodStates = map[int]string{
        -1:   "Decommissioned",
        100:  "Disabled",
        300:  "Maintenance",
        400:  "Test",
        500:  "Pre-Production",
        600:  "Migrated",
        900:  "Customer-Only",
        1000: "Production",
    }

    severityCodes = map[int]string{
        0: "Clear",
        1: "Debug",
        2: "Info",
        3: "Warning",
        4: "Error",
        5: "Critical",
    }
)

func getDevices(api *zenoss.API, query string) ([]zenoss.Device, error) {
    q := zenoss.GetDevicesQuery{
        Limit: 10,
        Keys: []string{
            "name",
            "uid",
            "productionState",
            "osManufacturer",
            "events",
        },
        Params: map[string]interface{}{
            "name": query,
        },
    }

    d, res, err := api.GetDevices(q)
    if err != nil {
        return nil, err
    }

    if res.StatusCode != 200 {
        return nil, fmt.Errorf("Error fetching devices. StatusCode: %d", res.StatusCode)
    }

    return d.Result.Devices, nil
}

func getEvents(api *zenoss.API, uid string) ([]zenoss.Event, error) {
    q := zenoss.QueryEventsQuery{
        UID: uid,
        Limit: 25,
        Keys: []string{
            "summary",
            "severity",
            "component",
            "count",
            "evid",
            "lastTime",
        },
        Params: map[string]interface{}{
            "eventState": []int{0,1},
            "severity": []int{3,4,5},
        },
    }

    e, res, err := api.QueryEvents(q)
    if err != nil {
        return nil, err
    }

    if res.StatusCode != 200 {
        return nil, fmt.Errorf("Error querying events. StatusCode: %d", res.StatusCode)
    }

    return e.Result.Events, nil
}
