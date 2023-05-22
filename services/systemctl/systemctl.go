package systemctl

import (
	"errors"
	"fmt"
	"github.com/NubeIO/flow-framework/installer"
	"github.com/NubeIO/flow-framework/utils/namings"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/sergeymakinen/go-systemdconf/v2"
	"github.com/sergeymakinen/go-systemdconf/v2/unit"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
)

type ServiceFile struct {
	Name                        string   `json:"name"`
	Version                     string   `json:"version"`
	ServiceDescription          string   `json:"service_description"`
	RunAsUser                   string   `json:"run_as_user"`
	ServiceWorkingDirectory     string   `json:"service_working_directory"`        // /data/rubix-service/apps/install/flow-framework/v0.6.1/
	ExecStart                   string   `json:"exec_start"`                       // app -p 1660 -g <data_dir> -d data -prod
	AttachWorkingDirOnExecStart bool     `json:"attach_working_dir_on_exec_start"` // true, false
	EnvironmentVars             []string `json:"environment_vars"`                 // Environment="g=/data/bacnet-server-c"
}

func GenerateServiceFile(app *ServiceFile, installer *installer.Installer) (tmpDir, absoluteServiceFileName string, err error) {
	tmpFilePath, err := installer.MakeTmpDirUpload()
	if err != nil {
		return "", "", err
	}
	if app.Name == "" {
		return "", "", errors.New("app name can not be empty, try flow-framework")
	}
	if app.Version == "" {
		return "", "", errors.New("app version can not be empty, try v0.6.0")
	}
	if err = checkVersion(app.Version); err != nil {
		return "", "", err
	}
	workingDirectory := app.ServiceWorkingDirectory
	if workingDirectory == "" {
		workingDirectory = installer.GetAppInstallPathWithVersion(app.Name, app.Version)
	}
	log.Infof("generate service working dir: %s", workingDirectory)
	user := app.RunAsUser
	if user == "" {
		user = "root"
	}
	execCmd := app.ExecStart
	if app.AttachWorkingDirOnExecStart {
		workingDir := installer.GetAppInstallPathWithVersion(app.Name, app.Version)
		execCmd = path.Join(workingDir, execCmd)
		if !strings.Contains(execCmd, installer.GetAppInstallPathWithVersion(app.Name, app.Version)) {
			return "", "", errors.New(fmt.Sprintf(
				"ExecStart command is not matching app_name: %s & app_version: %s", app.Name, app.Version))
		}
	}
	if strings.Contains(execCmd, "<root_dir>") {
		execCmd = strings.ReplaceAll(execCmd, "<root_dir>", installer.RootDir)
	}
	if strings.Contains(execCmd, "<data_dir>") {
		execCmd = strings.ReplaceAll(execCmd, "<data_dir>", installer.GetAppDataPath(app.Name))
	}
	if strings.Contains(execCmd, "<data_dir_name>") {
		execCmd = strings.ReplaceAll(execCmd, "<data_dir_name>", namings.GetDataDirNameFromAppName(app.Name))
	}
	log.Infof("generate service exec_cmd: %s", execCmd)
	description := app.ServiceDescription
	if description == "" {
		description = fmt.Sprintf("NubeIO %s", app.Name)
	}
	var env systemdconf.Value
	for _, s := range app.EnvironmentVars {
		env = append(env, s)
	}
	service := unit.ServiceFile{
		Unit: unit.UnitSection{ // [Unit]
			Description: systemdconf.Value{description},
			After:       systemdconf.Value{"network.target"},
		},
		Service: unit.ServiceSection{ // [Service]
			ExecStartPre: nil,
			Type:         systemdconf.Value{"simple"},
			ExecOptions: unit.ExecOptions{
				User:             systemdconf.Value{user},
				WorkingDirectory: systemdconf.Value{workingDirectory},
				Environment:      env,
				StandardOutput:   systemdconf.Value{"syslog"},
				StandardError:    systemdconf.Value{"syslog"},
				SyslogIdentifier: systemdconf.Value{app.Name},
			},
			ExecStart: systemdconf.Value{
				execCmd,
			},
			Restart:    systemdconf.Value{"always"},
			RestartSec: systemdconf.Value{"10"},
		},
		Install: unit.InstallSection{ // [Install]
			WantedBy: systemdconf.Value{"multi-user.target"},
		},
	}
	b, _ := systemdconf.Marshal(service)
	serviceName := namings.GetServiceNameFromAppName(app.Name)
	absoluteServiceFileName = fmt.Sprintf("%s/%s", tmpFilePath, serviceName)
	err = fileutils.WriteFile(absoluteServiceFileName, string(b), os.FileMode(installer.FileMode))
	if err != nil {
		log.Errorf("write service file error %s", err.Error())
	}
	log.Infof("generate service file name: %s", serviceName)
	return tmpFilePath, absoluteServiceFileName, nil
}

func checkVersion(version string) error {
	if version[0:1] != "v" { // make sure have a v at the start v0.1.1
		return errors.New(fmt.Sprintf("incorrect provided: %s version number try: v1.2.3", version))
	}
	p := strings.Split(version, ".")
	if len(p) != 3 {
		return errors.New(fmt.Sprintf("incorrect length provided: %s version number try: v1.2.3", version))
	}
	return nil
}
