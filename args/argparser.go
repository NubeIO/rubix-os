package args

import (
	"encoding/json"
	"strconv"
	"strings"
)

func (a Args) SerializeArgs(args Args) (string, error) {
	argsData, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	argsString := string(argsData)
	return argsString, nil
}

func DeserializeArgs(args string) (*Args, error) {
	deserializedArgs := Args{}
	if len(args) == 0 {
		return &deserializedArgs, nil
	}
	err := json.Unmarshal([]byte(args), &deserializedArgs)
	if err != nil {
		return nil, err
	}
	return &deserializedArgs, nil
}

func ParseArgs(args string) Args {
	apiArgs := Args{}
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
