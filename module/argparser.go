package module

import (
	"encoding/json"
	"github.com/NubeIO/rubix-os/api"
	"strconv"
	"strings"
)

func SerializeArgs(args api.Args) (string, error) {
	argsData, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	argsString := string(argsData)
	return argsString, nil
}

func parseArgs(args string) api.Args {
	apiArgs := api.Args{}
	argsParts := strings.Split(args, "&&")
	for _, arg := range argsParts {
		argParts := strings.Split(arg, "=")
		if len(argParts) == 2 {
			if argParts[0] == "with_devices" {
				r, _ := strconv.ParseBool(argParts[1])
				apiArgs.WithDevices = r
			} else if argParts[0] == "with_points" {
				r, _ := strconv.ParseBool(argParts[1])
				apiArgs.WithPoints = r
			} else if argParts[0] == "with_priority" {
				r, _ := strconv.ParseBool(argParts[1])
				apiArgs.WithPriority = r
			} else if argParts[0] == "address_uuid" {
				apiArgs.AddressUUID = &argParts[1]
			} else if argParts[0] == "io_number" {
				apiArgs.IoNumber = &argParts[1]
			}
		}
	}
	return apiArgs
}
