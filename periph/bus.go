package periph

import (
	"fmt"
	"github.com/bearsh/hid"
	"github.com/ngyewch/ft260"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"time"
)

// Bus interface for the FT260 device.
type Bus struct {
	name string
	dev  *ft260.Dev
}

// NewBus constructs a new Bus instance.
func NewBus(name string, dev *ft260.Dev) *Bus {
	return &Bus{
		name: name,
		dev:  dev,
	}
}

// Close closes the bus.
func (bus *Bus) Close() error {
	return bus.dev.Close()
}

// String returns the bus name.
func (bus *Bus) String() string {
	return bus.name
}

// Tx does a single transaction.
func (bus *Bus) Tx(addr uint16, w, r []byte) error {
	if addr >= 0x80 {
		return fmt.Errorf("address out of range")
	}
	if w != nil {
		_, err := bus.dev.I2CWriteRequest(uint8(addr), w, ft260.I2CConditionStartAndStop)
		if err != nil {
			return err
		}
	}
	if (w != nil) && (r != nil) {
		time.Sleep(1 * time.Millisecond)
	}
	if r != nil {
		err := bus.dev.I2CReadRequest(uint8(addr), ft260.I2CConditionStartAndStop, uint16(len(r)))
		if err != nil {
			return err
		}
		readBytes, err := bus.dev.I2CInputReport()
		if err != nil {
			return err
		}
		for i, b := range readBytes {
			if i >= len(r) {
				break
			}
			r[i] = b
		}
	}
	return nil
}

// SetSpeed changes the bus speed, if supported.
func (bus *Bus) SetSpeed(f physic.Frequency) error {
	speedInKHz := int64(f) / 1_000_000_000
	if (speedInKHz < 0) || (speedInKHz > 65535) {
		return fmt.Errorf("invalid speed")
	}
	return bus.dev.SetI2CClockSpeed(uint16(speedInKHz))
}

// Register registers all enumerated FT260 I²C buses.
func Register() error {
	deviceInfoList := hid.Enumerate(ft260.VendorID, ft260.ProductID)

	for i, deviceInfo := range deviceInfoList {
		name := fmt.Sprintf("ft260-%d", i)
		err := i2creg.Register(name, nil, -1, newOpener(name, deviceInfo))
		if err != nil {
			return err
		}
	}

	return nil
}

var (
	hidDeviceMap = make(map[string]*hid.Device)
)

func newOpener(name string, deviceInfo hid.DeviceInfo) i2creg.Opener {
	return func() (i2c.BusCloser, error) {
		hidDevice, ok := hidDeviceMap[name]
		if !ok {
			hidDevice1, err := deviceInfo.Open()
			if err != nil {
				return nil, err
			}
			hidDevice = hidDevice1
			hidDeviceMap[name] = hidDevice
		}
		dev := ft260.New(hidDevice)
		bus := NewBus(name, dev)
		return bus, nil
	}
}
