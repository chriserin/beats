package devices

import (
	"fmt"
	"strings"

	"github.com/rakyll/portmidi"
)

//FindDeviceID func
func FindDeviceID(deviceName string) portmidi.DeviceID {
	fmt.Println("Looking for " + deviceName)
	for i := 0; i < portmidi.CountDevices(); i++ {
		info := portmidi.Info(portmidi.DeviceID(i))
		if strings.Contains(info.Name, deviceName) && info.IsOutputAvailable {
			return portmidi.DeviceID(i)
		}
	}
	return portmidi.DeviceID(2)
}
