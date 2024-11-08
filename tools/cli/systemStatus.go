package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
)

func doSystemStatus(cCtx *cli.Context) error {
	dev, err := getDevice(cCtx)
	if err != nil {
		return err
	}
	systemStatus, err := dev.GetSystemStatus()
	if err != nil {
		return err
	}
	jsonBytes, err := json.MarshalIndent(systemStatus, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonBytes))
	return nil
}
