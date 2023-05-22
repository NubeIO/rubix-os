package system

import (
	"fmt"
	"github.com/NubeIO/flow-framework/utils/pprint"
	"testing"
)

func TestSystem_DiscUsage(t *testing.T) {
	sys := New(&System{})
	r, err := sys.GetSystem()
	if err != nil {
		fmt.Println(err)
	}
	pprint.PrintJSON(r)
}
