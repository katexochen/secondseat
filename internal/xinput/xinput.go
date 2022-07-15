package xinput

import (
	"fmt"
)

const requiredXVersion = "1.20"

type Xinput struct {
	Name     string
	ID       int
	Type     DeviceType
	Role     DeviceRole
	MasterID int
}

type DeviceType string

const (
	Keyboard DeviceType = "keyboard"
	Pointer  DeviceType = "pointer"
)

func parseDeviceType(s string) (DeviceType, error) {
	switch {
	case s == "keyboard":
		return Keyboard, nil
	case s == "pointer":
		return Pointer, nil
	default:
		return "", fmt.Errorf("unknown device type %s", s)
	}
}

type DeviceRole string

const (
	Master DeviceRole = "master"
	Slave  DeviceRole = "slave"
)

func parseDeviceRole(s string) (DeviceRole, error) {
	switch {
	case s == "master":
		return Master, nil
	case s == "slave":
		return Slave, nil
	default:
		return "", fmt.Errorf("unknown device role %s", s)
	}
}
