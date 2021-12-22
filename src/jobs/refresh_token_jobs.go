package jobs

import "fmt"

func (j *Jobs) RefreshTokenJobAdd() error {
	_, err := cron.Every(1).Hour().Tag("refreshToken").Do(j.refreshToken)
	if err != nil {
		return err
	}
	return nil
}

func (j *Jobs) refreshToken() {
	fmt.Println("REFRESH TOKEN RUN")
	_, err := j.db.RefreshLocalStorageFlowToken()
	_, err = j.db.RefreshFlowNetworksConnections()
	_, err = j.db.RefreshFlowNetworkClonesConnections()
	if err != nil {
		//TODO FIX ERROR
	}
}
