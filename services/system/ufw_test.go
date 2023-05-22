package system

import (
	"testing"
)

func TestSystem_UWFStatus(t *testing.T) {
	sys := New(&System{})

	sys.UWFStatus()
}
