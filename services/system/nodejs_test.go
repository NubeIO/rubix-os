package system

import (
	"testing"
)

func TestSystem_NodeGetVersion(t *testing.T) {
	sys := New(&System{})

	sys.NodeGetVersion()
}
