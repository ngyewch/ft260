package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"
)

func doSystemStatus(ctx context.Context, cmd *cli.Command) error {
	dev, err := getDevice(ctx, cmd)
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
