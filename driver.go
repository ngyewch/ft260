package ft260

import (
	"encoding/binary"
	"fmt"
	"github.com/bearsh/hid"
)

const (
	VendorID  = 0x0403 // VendorID is the USB vendor ID.
	ProductID = 0x6030 // ProductID is the USB product ID.
)

type SystemClockSpeed uint8

const (
	SystemClockSpeed12MHz SystemClockSpeed = 0
	SystemClockSpeed24MHz SystemClockSpeed = 1
	SystemClockSpeed48MHz SystemClockSpeed = 2
)

type UARTMode uint8

const (
	UARTModeOff           UARTMode = 0
	UARTModeCtsRts        UARTMode = 1
	UARTModeDtrDts        UARTMode = 2
	UARTModeXonXoff       UARTMode = 3
	UARTModeNoFlowControl UARTMode = 4
)

type I2CCondition uint8

const (
	I2CConditionNone          I2CCondition = 0
	I2CConditionStart         I2CCondition = 0x02
	I2CConditionRepeatedStart I2CCondition = 0x03
	I2CConditionStop          I2CCondition = 0x04
	I2CConditionStartAndStop  I2CCondition = 0x06
)

// Dev is the device handle.
type Dev struct {
	dev *hid.Device
}

// ChipVersion contains the results of a Chip Version operation.
type ChipVersion struct {
	ChipCode []byte
	Reserved []byte
}

// SystemStatus contains the results of a Get System Status operation.
type SystemStatus struct {
	ChipMode          uint8
	ClkCtl            SystemClockSpeed
	SuspendStatus     bool
	PowerEnableStatus bool
	I2CEnable         bool
	UARTMode          UARTMode
	HIDOverI2CEnable  bool
	GPIO2Function     uint8
	GPIOAFunction     uint8
	GPIOGFunction     uint8
	SuspendOutPol     uint8
	EnableWakeupInt   bool
	IntrCond          uint8
	EnablePowerSaving bool
	Reserved          []byte
}

// New creates a new device handle.
func New(dev *hid.Device) *Dev {
	return &Dev{
		dev: dev,
	}
}

// Close closes the device.
func (dev *Dev) Close() error {
	return dev.dev.Close()
}

// ChipVersion performs a Chip Version operation.
func (dev *Dev) ChipVersion() (*ChipVersion, error) {
	b := make([]byte, 13)
	b[0] = 0xa0
	_, err := dev.dev.GetFeatureReport(b)
	if err != nil {
		return nil, err
	}
	return &ChipVersion{
		ChipCode: b[1:5],
		Reserved: b[5:],
	}, nil
}

// GetSystemStatus performs a Get System Status operation.
func (dev *Dev) GetSystemStatus() (*SystemStatus, error) {
	b := make([]byte, 26)
	b[0] = 0xa1
	_, err := dev.dev.GetFeatureReport(b)
	if err != nil {
		return nil, err
	}
	return &SystemStatus{
		ChipMode:          b[1],
		ClkCtl:            SystemClockSpeed(b[2]),
		SuspendStatus:     b[3] == 1,
		PowerEnableStatus: b[4] == 1,
		I2CEnable:         b[5] == 1,
		UARTMode:          UARTMode(b[6]),
		HIDOverI2CEnable:  b[7] == 1,
		GPIO2Function:     b[8],
		GPIOAFunction:     b[9],
		GPIOGFunction:     b[10],
		SuspendOutPol:     b[11],
		EnableWakeupInt:   b[12] == 1,
		IntrCond:          b[13],
		EnablePowerSaving: b[14] == 1,
		Reserved:          b[15:],
	}, nil
}

// SetI2CMode performs a Set I2C Mode operation.
func (dev *Dev) SetI2CMode(enable bool) error {
	var v byte
	if enable {
		v = 1
	}
	b := []byte{0xa1, 0x02, v}
	_, err := dev.dev.SendFeatureReport(b)
	return err
}

// I2CReset performs a I²C Reset operation.
func (dev *Dev) I2CReset() error {
	b := []byte{0xa1, 0x20}
	_, err := dev.dev.SendFeatureReport(b)
	return err
}

// SetI2CClockSpeed performs a Set I²C Clock Speed operation.
func (dev *Dev) SetI2CClockSpeed(speedInKhz uint16) error {
	b := []byte{0xa1, 0x22}
	b = binary.LittleEndian.AppendUint16(b, speedInKhz)
	_, err := dev.dev.SendFeatureReport(b)
	return err
}

// I2CWriteRequest performs an I²C Write Request operation.
func (dev *Dev) I2CWriteRequest(slaveAddr uint8, data []byte, flag I2CCondition) (int, error) {
	if slaveAddr >= 0x80 {
		return 0, fmt.Errorf("invalid I2C slave address")
	}
	if len(data) == 0 {
		return 0, fmt.Errorf("data must not be empty")
	}
	if len(data) > 60 { // TODO support longer data
		return 0, fmt.Errorf("data must not exceed 60 bytes")
	}
	b := []byte{0xd0 + byte((len(data)-1)/4), slaveAddr, byte(flag), byte(len(data))}
	b = append(b, data...)
	return dev.dev.Write(b)
}

// I2CReadRequest performs an I²C Read Request operation.
func (dev *Dev) I2CReadRequest(slaveAddr uint8, flag I2CCondition, dataLength uint16) error {
	if slaveAddr >= 0x80 {
		return fmt.Errorf("invalid I2C slave address")
	}
	if dataLength == 0 {
		return fmt.Errorf("data length must be greater than 0")
	}
	if dataLength > 62 { // TODO support longer data
		return fmt.Errorf("data length must not exceed 62 bytes")
	}
	b := []byte{0xc2, slaveAddr, byte(flag)}
	b = binary.LittleEndian.AppendUint16(b, dataLength)
	_, err := dev.dev.Write(b)
	return err
}

// I2CInputReport performs an I²C Input Report operation.
func (dev *Dev) I2CInputReport() ([]byte, error) {
	b := make([]byte, 64)
	_, err := dev.dev.Read(b)
	if err != nil {
		return nil, err
	}
	dataLength := int(b[1])
	return b[2 : 2+dataLength], nil
}

// I2CInputReportWithTimeout performs an I²C Input Report operation with timeout (ms).
// If timeout is -1, a blocking read is performed.
func (dev *Dev) I2CInputReportWithTimeout(timeout int) ([]byte, error) {
	b := make([]byte, 64)
	_, err := dev.dev.ReadTimeout(b, timeout)
	if err != nil {
		return nil, err
	}
	dataLength := int(b[1])
	return b[2 : 2+dataLength], nil
}
