package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"
)

func doChipVersion(ctx context.Context, cmd *cli.Command) error {
	dev, err := getDevice(ctx, cmd)
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
