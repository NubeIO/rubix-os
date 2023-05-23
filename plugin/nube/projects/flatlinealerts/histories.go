package main

import (
	"errors"
	"github.com/NubeIO/rubix-os/plugin/nube/projects/flatlinealerts/ffhistoryrest"
)

func (inst *Instance) GetFFToken(user, pass string) (*ffhistoryrest.FFToken, error) {
	inst.flatlinealertsDebugMsg("GetFFToken()")
	host := inst.config.Job.FFHost
	// host := "0.0.0.0"
	if host == "" {
		host = "0.0.0.0"
	}
	port := inst.config.Job.FFPort
	// port := 1616
	if port <= 0 {
		port = 1616
	}
	rest := ffhistoryrest.NewNoAuth(host, int(port))
	token, err := rest.GetFFToken(user, pass)
	if err != nil {
		inst.flatlinealertsErrorMsg(err)
	}
	if err != nil {
		return nil, errors.New("could not get ff token")
	}
	return token, nil
}

func (inst *Instance) GetFFHistories(FFToken ffhistoryrest.FFToken, queryParams string) ([]ffhistoryrest.FFHistory, error) {
	inst.flatlinealertsDebugMsg("GetFFHistories()")
	host := inst.config.Job.FFHost
	// host := "0.0.0.0"
	if host == "" {
		host = "0.0.0.0"
	}
	port := inst.config.Job.FFPort
	// port := 8080
	if port <= 0 {
		port = 8080
	}
	rest := ffhistoryrest.NewNoAuth(host, int(port))
	ffHistoryArray, err := rest.GetFFHistories(FFToken, queryParams)
	if err != nil || ffHistoryArray == nil {
		inst.flatlinealertsErrorMsg(err)
		return nil, errors.New("could not get ff histories")
	}
	return *ffHistoryArray, nil
}
