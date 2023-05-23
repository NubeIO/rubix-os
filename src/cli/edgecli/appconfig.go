package edgecli

import (
	"errors"
	"github.com/NubeIO/rubix-os/global"
	"github.com/NubeIO/rubix-os/interfaces"
	log "github.com/sirupsen/logrus"
	"path"
)

func (inst *Client) EdgeReadConfig(appName, configName string) (*interfaces.EdgeConfigResponse, error) {
	appDataConfigPath := global.Installer.GetAppDataConfigPath(appName)
	absoluteAppDataConfigName := path.Join(appDataConfigPath, configName)
	file, err := inst.ReadFile(absoluteAppDataConfigName)
	if err != nil {
		return nil, err
	}
	return &interfaces.EdgeConfigResponse{
		Data:     file,
		FilePath: absoluteAppDataConfigName,
	}, nil
}

func (inst *Client) EdgeWriteConfig(body *interfaces.EdgeConfig) (*interfaces.Message, error) {
	if body.AppName == "" {
		return nil, errors.New("app_name can not be empty")
	}
	configName := body.ConfigName
	if configName == "" {
		configName = "config.yml"
	}
	appDataConfigPath := global.Installer.GetAppDataConfigPath(body.AppName)
	dirExistence, err := inst.DirExists(appDataConfigPath)
	if err != nil {
		return nil, err
	}
	if !dirExistence.Exists {
		dir, err := inst.CreateDir(appDataConfigPath)
		if err != nil {
			return nil, err
		}
		log.Infof("made config dir as was not existing: %s", dir.Message)
	}
	absoluteAppDataConfigName := path.Join(appDataConfigPath, configName)

	writeFile := &interfaces.WriteFile{
		FilePath:     absoluteAppDataConfigName,
		Body:         body.Body,
		BodyAsString: body.BodyAsString,
	}
	if configName == "config.yml" {
		return inst.WriteFileYml(writeFile)
	} else if configName == ".env" {
		return inst.WriteString(writeFile)
	} else if configName == "config.json" {
		return inst.WriteFileJson(writeFile)
	}

	return nil, errors.New("no valid config_name, config.yml or .env or config.json")
}
