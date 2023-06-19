package rubixregistry

import (
	"fmt"
	"testing"
)

func Test_GetDeviceInfo(*testing.T) {
	rr := New()
	deviceInfo, err := rr.GetDeviceInfo()
	fmt.Println("err", err)
	fmt.Println("deviceInfo", deviceInfo)
	if deviceInfo != nil {
		fmt.Println("deviceInfo.GlobalUUID", deviceInfo.GlobalUUID)
	}
}
