package system

import (
	"github.com/NubeIO/rubix-os/utils/pprint"
	"testing"
)

func TestSystem_GetHardwareClock(t *testing.T) {
	clock, err := New(&System{}).GetHardwareClock()
	if err != nil {
		return
	}
	pprint.PrintJSON(clock)
}
