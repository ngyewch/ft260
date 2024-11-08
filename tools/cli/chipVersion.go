package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
)

func doChipVersion(cCtx *cli.Context) error {
	dev, err := getDevice(cCtx)
	if err != nil {
		return err
	}
	chipVersion, err := dev.ChipVersion()
	if err != nil {
		return err
	}
	jsonBytes, err := json.MarshalIndent(chipVersion, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonBytes))
	return nil
}
