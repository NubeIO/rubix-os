package utils

import (
	"encoding/json"
	"errors"
	"github.com/NubeIO/flow-framework/config"
	"github.com/NubeIO/flow-framework/interfaces"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var ExcludedServices = []string{
	"nubeio-rubix-edge-bios.service",
	"nubeio-rubix-edge.service",
	"nubeio-rubix-assist.service",
}

var RestartJobFile = "restart_job.json"
var RebootJobFile = "reboot_job.json"

func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func DeleteDir(source, parentDirectory string, depth int) error {
	dir, _ := os.Open(source)
	obs, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	var errs []error

	for _, obj := range obs {
		fSource := path.Join(source, obj.Name())
		if obj.IsDir() {
			if parentDirectory == "rubix-service/apps/install" &&
				!Contains([]string{"rubix-edge", "rubix-assist"}, obj.Name()) {
				_ = os.RemoveAll(path.Join(parentDirectory, obj.Name()))
			}
			err = DeleteDir(fSource, path.Join(parentDirectory, obj.Name()), depth+1)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	var errString string
	for _, err := range errs {
		errString += err.Error() + "\n"
	}
	if errString != "" {
		return errors.New(errString)
	}
	return nil
}

func CopyDir(source, dest, parentDirectory string, depth int) error {
	srcInfo, err := os.Stat(source)
	if err != nil {
		return err
	}
	err = os.MkdirAll(dest, srcInfo.Mode())
	if err != nil {
		return err
	}
	dir, _ := os.Open(source)
	obs, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	var errs []error
	for _, obj := range obs {
		fSource := path.Join(source, obj.Name())
		fDest := path.Join(dest, obj.Name())
		if obj.IsDir() {
			excludesData := []string{
				"rubix-edge",
				"rubix-assist",
				"tmp",
				"store",
				"backup",
				"socat",
			}
			excludesApps := []string{
				"rubix-service/apps/install/rubix-edge",
				"rubix-service/apps/install/rubix-assist",
			}
			if !((Contains(excludesData, obj.Name()) && depth == 0) ||
				(Contains(excludesApps, path.Join(parentDirectory, obj.Name())) && depth == 3)) {
				err = CopyDir(fSource, fDest, path.Join(parentDirectory, obj.Name()), depth+1)
				if err != nil {
					log.Error(err)
					errs = append(errs, err)
				}
			}
		} else {
			err = fileutils.CopyFile(fSource, fDest)
			if err != nil {
				log.Error(err)
				errs = append(errs, err)
			}
		}
	}
	var errString string
	for _, err := range errs {
		errString += err.Error() + "\n"
	}
	if errString != "" {
		return errors.New(errString)
	}
	return nil
}

func CopyFiles(srcFiles []string, dest string) {
	var wg sync.WaitGroup
	for _, srcFile := range srcFiles {
		wg.Add(1)
		go func(srcFile string) {
			defer wg.Done()
			if !Contains(ExcludedServices, srcFile) {
				err := fileutils.CopyFile(srcFile, path.Join(dest, filepath.Base(srcFile)))
				if err != nil {
					log.Errorf("failed to copy file %s to %s", srcFile, path.Join(dest, filepath.Base(srcFile)))
				}
			}
		}(srcFile)
	}
	wg.Wait()
}

func DeleteFiles(srcFiles []string, dest string) {
	var wg sync.WaitGroup
	for _, srcFile := range srcFiles {
		wg.Add(1)
		go func(srcFile string) {
			defer wg.Done()
			err := os.RemoveAll(path.Join(dest, filepath.Base(srcFile)))
			if err != nil {
				log.Errorf("failed to remove file %s", os.RemoveAll(path.Join(dest, filepath.Base(srcFile))))
			}
		}(srcFile)
	}
	wg.Wait()
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func ValidateCornExpression(exp string) error {
	expFields := strings.Fields(exp)
	if len(expFields) != 5 {
		return errors.New("invalid expression")
	}
	_, err := cron.ParseStandard(exp)
	if err != nil {
		return errors.New("invalid expression")
	}
	if (expFields[0] == "*" || strings.Contains(expFields[0], "*/")) && expFields[1] == "*" &&
		expFields[2] == "*" && expFields[3] == "*" && expFields[4] == "*" {
		return errors.New("you cannot schedule under an hour")
	}
	return nil
}

func GetRestartJobs() []*interfaces.RestartJob {
	var restartJobs []*interfaces.RestartJob
	data, err := ioutil.ReadFile(path.Join(config.Get().GetAbsDataDir(), RestartJobFile))
	if err != nil {
		return []*interfaces.RestartJob{}
	}
	err = json.Unmarshal(data, &restartJobs)
	if err != nil {
		return []*interfaces.RestartJob{}
	}
	return restartJobs
}

func SaveRestartJobs(restartJobs []*interfaces.RestartJob, fileMode int) error {
	mRestartJobs, err := json.Marshal(restartJobs)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(config.Get().GetAbsDataDir(), RestartJobFile), mRestartJobs,
		fs.FileMode(fileMode))
}

func GetRebootJob() *interfaces.RebootJob {
	var rebootJob *interfaces.RebootJob
	data, err := ioutil.ReadFile(path.Join(config.Get().GetAbsDataDir(), RebootJobFile))
	if err != nil {
		return rebootJob
	}
	_ = json.Unmarshal(data, &rebootJob)
	return rebootJob
}

func SaveRebootJob(rebootJob *interfaces.RebootJob, fileMode int) error {
	mRebootJob, err := json.Marshal(rebootJob)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(config.Get().GetAbsDataDir(), RebootJobFile), mRebootJob,
		fs.FileMode(fileMode))
}
