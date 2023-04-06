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
            "groups",
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
