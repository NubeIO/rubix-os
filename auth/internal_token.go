package auth

import (
	"fmt"
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/utils/security"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func GetInternalToken(withPrefix bool) string {
	conf := config.Get()
	absInternalTokenFile := conf.GetAbsInternalTokenFile()
	file, err := os.Open(absInternalTokenFile)
	if err != nil {
		log.Error(err)
		return ""
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Error(err)
		}
	}()
	internalToken, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err)
	}
	if withPrefix {
		return fmt.Sprintf("Internal %s", string(internalToken))
	} else {
		return string(internalToken)
	}
}

func CreateInternalToken() {
	conf := config.Get()
	if err := os.MkdirAll(conf.Location.TokenFolder, 0755); err != nil {
		panic(err)
	}
	absInternalTokenFile := conf.GetAbsInternalTokenFile()
	file, err := os.OpenFile(absInternalTokenFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Error(err)
	}
	token := security.GenerateToken()
	_, err = file.Write([]byte(token))
	if err != nil {
		log.Error(err)
	}
}
