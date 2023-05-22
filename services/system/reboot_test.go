package system

import (
	"testing"
)

func TestSystem_RebootHost(t *testing.T) {
	sys := New(&System{})
	sys.RebootHost()

}
