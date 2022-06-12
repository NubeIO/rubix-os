package auth

import (
	"fmt"
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/utils/file"
	"github.com/NubeIO/flow-framework/utils/security"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func GetInternalToken(withPrefix bool) string {
	conf := config.Get()
	absInternalTokenFile := conf.GetAbsInternalTokenFile()
	f, err := os.Open(absInternalTokenFile)
	if err != nil {
		log.Error(err)
		return ""
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Error(err)
		}
	}()
	internalToken, err := ioutil.ReadAll(f)
	if err != nil {
		log.Error(err)
	}
	if withPrefix {
		return fmt.Sprintf("Internal %s", string(internalToken))
	} else {
		return string(internalToken)
	}
}

func CreateInternalTokenIfDoesNotExist() {
	conf := config.Get()
	if err := os.MkdirAll(conf.Location.TokenFolder, 0755); err != nil {
		panic(err)
	}
	absInternalTokenFile := conf.GetAbsInternalTokenFile()
	_, err := file.ReadFile(absInternalTokenFile)
	if err == nil {
		return
	}
	f, err := os.OpenFile(absInternalTokenFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Error(err)
	}
	token := security.GenerateToken()
	_, err = f.Write([]byte(token))
	if err != nil {
		log.Error(err)
	}
}
