package auth

import (
	"fmt"
	"github.com/NubeDev/flow-framework/config"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
)

func GetRubixServiceInternalToken() string {
	conf := config.Get()
	authDataDir := conf.Location.AuthDataDir
	relativeAuthDataFile := conf.Location.RelativeAuthDataFile
	authDataFile := path.Join(authDataDir, relativeAuthDataFile)
	file, err := os.Open(authDataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	internalToken, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err)
	}
	return fmt.Sprintf("Internal %s", string(internalToken))
}
