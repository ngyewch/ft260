package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/bearsh/hid"
	"github.com/ngyewch/ft260"
	"github.com/urfave/cli/v3"
)

var (
	deviceIndexFlag = &cli.UintFlag{
		Name:    "device-index",
		Sources: cli.EnvVars("DEVICE_INDEX"),
	}

	app = &cli.Command{
		Name:  "ft260",
		Usage: "FT260 CLI",
		Flags: []cli.Flag{
			deviceIndexFlag,
		},
		Commands: []*cli.Command{
			{
				Name:    "chipVersion",
				Aliases: []string{"chip-version"},
				Action:  doChipVersion,
			},
			{
				Name:    "systemStatus",
				Aliases: []string{"system-status"},
				Action:  doSystemStatus,
			},
		},
	}
)

func main() {
	buildInfo, _ := debug.ReadBuildInfo()
	if buildInfo != nil {
		app.Version = buildInfo.Main.Version
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getDevice(ctx context.Context, cmd *cli.Command) (*ft260.Dev, error) {
	deviceIndex := cmd.Uint(deviceIndexFlag.Name)
	deviceInfoList := hid.Enumerate(ft260.VendorID, ft260.ProductID)
	if int(deviceIndex) >= len(deviceInfoList) {
		return nil, fmt.Errorf("device index out of range")
	}
	dev, err := deviceInfoList[deviceIndex].Open()
	if err != nil {
		return nil, err
	}
	return ft260.New(dev), nil
}
