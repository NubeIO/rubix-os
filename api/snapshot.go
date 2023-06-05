package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/lib-systemctl-go/systemctl"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/config"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/src/cli/bioscli"
	"github.com/NubeIO/rubix-os/src/cli/constants"
	"github.com/NubeIO/rubix-os/utils"
	"github.com/NubeIO/rubix-os/utils/namings"
	"github.com/NubeIO/rubix-registry-go/rubixregistry"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var dataFolder = "data"
var systemFolder = "system"

var createStatus = interfaces.CreateNotAvailable
var restoreStatus = interfaces.RestoreNotAvailable

type SnapshotAPI struct {
	RubixRegistry *rubixregistry.RubixRegistry
	FileMode      int
	SystemCtl     *systemctl.SystemCtl
}

func (a *SnapshotAPI) CreateSnapshot(c *gin.Context) {
	log.Info("creating snapshot...")
	if createStatus == interfaces.Creating {
		err := errors.New("snapshot creation process is in progress")
		log.Error(err)
		ResponseHandler(nil, err, c)
		return
	}
	createStatus = interfaces.Creating
	deviceInfo, err := a.RubixRegistry.GetDeviceInfo()
	if err != nil {
		log.Error(err)
		createStatus = interfaces.CreateFailed
		ResponseHandler(nil, err, c)
		return
	}
	clientName := strings.Replace(deviceInfo.ClientName, "/", "", -1)
	siteName := strings.Replace(deviceInfo.SiteName, "/", "", -1)
	deviceName := strings.Replace(deviceInfo.DeviceName, "/", "", -1)
	if clientName == "" || clientName == "-" {
		clientName = "na"
	}
	if siteName == "" || siteName == "-" {
		siteName = "na"
	}
	if deviceName == "" || deviceName == "-" {
		deviceName = "na"
	}
	filePrefix := strings.ReplaceAll(fmt.Sprintf("%s-%s-%s", clientName, siteName, deviceName), " ", "")
	previousFiles, _ := filepath.Glob(path.Join(config.Get().GetAbsTempDir(), fmt.Sprintf("%s*", filePrefix)))
	utils.DeleteFiles(previousFiles, config.Get().GetAbsTempDir())

	biosClient := bioscli.NewLocalBiosClient()
	arch, err := biosClient.GetArch()
	if err != nil {
		log.Error(err)
		createStatus = interfaces.CreateFailed
		ResponseHandler(nil, err, c)
		return
	}

	destinationPath := fmt.Sprintf("%s/%s_%s_%s", config.Get().GetAbsTempDir(), filePrefix,
		time.Now().UTC().Format("20060102T150405"), arch.Arch)
	absDataFolder := path.Join(destinationPath, dataFolder)
	err = os.MkdirAll(absDataFolder, os.FileMode(a.FileMode)) // create empty folder even we don't have content
	if err != nil {
		log.Error(err)
		createStatus = interfaces.CreateFailed
		ResponseHandler(nil, err, c)
		return
	}
	_ = utils.CopyDir(config.Get().GetSnapshotDir(), absDataFolder, "", 0)

	systemFiles, err := filepath.Glob(path.Join(constants.ServiceDirSoftLink, "nubeio-*"))
	if err != nil {
		log.Error(err)
		createStatus = interfaces.CreateFailed
		ResponseHandler(nil, err, c)
		return
	}
	absSystemFolder := path.Join(destinationPath, systemFolder)
	err = os.MkdirAll(absSystemFolder, os.FileMode(a.FileMode)) // create empty folder even we don't have content
	if err != nil {
		log.Error(err)
		createStatus = interfaces.CreateFailed
		ResponseHandler(nil, err, c)
		return
	}
	utils.CopyFiles(systemFiles, absSystemFolder)

	zipDestinationPath := destinationPath + ".zip"
	log.Infof("zipping snapshot: %s...", zipDestinationPath)
	err = fileutils.RecursiveZip(destinationPath, zipDestinationPath)
	if err != nil {
		log.Error(err)
		createStatus = interfaces.CreateFailed
		ResponseHandler(nil, err, c)
		return
	}
	_ = os.RemoveAll(destinationPath)
	createStatus = interfaces.Created
	log.Info("sending snapshot data...")
	c.FileAttachment(zipDestinationPath, filepath.Base(zipDestinationPath))
}

func (a *SnapshotAPI) RestoreSnapshot(c *gin.Context) {
	log.Info("restoring snapshot...")
	if restoreStatus == interfaces.Restoring {
		err := errors.New("snapshot restoring process is in progress")
		log.Error(err)
		ResponseHandler(nil, err, c)
		return
	}
	log.Info("receiving file data...")
	restoreStatus = interfaces.Restoring
	file, err := c.FormFile("file")
	if err != nil {
		log.Error(err)
		restoreStatus = interfaces.RestoreFailed
		ResponseHandler(nil, err, c)
		return
	}

	biosClient := bioscli.NewLocalBiosClient()
	arch, err := biosClient.GetArch()
	if err != nil {
		log.Error(err)
		restoreStatus = interfaces.RestoreFailed
		ResponseHandler(nil, err, c)
		return
	}

	fileParts := strings.Split(file.Filename, "_")
	archParts := fileParts[len(fileParts)-1]
	archFromSnapshot := strings.Split(archParts, ".")[0]
	if archFromSnapshot != arch.Arch {
		restoreStatus = interfaces.RestoreFailed
		err = errors.New(
			fmt.Sprintf("arch mismatch: snapshot arch is %s & device arch is %s", archFromSnapshot, arch.Arch))
		log.Error(err)
		ResponseHandler(nil, err, c)
		return
	}

	log.Info("saving received file data...")
	destinationFilePath := path.Join(config.Get().GetAbsTempDir(), file.Filename)
	err = c.SaveUploadedFile(file, destinationFilePath)
	if err != nil {
		log.Error(err)
		restoreStatus = interfaces.RestoreFailed
		ResponseHandler(nil, err, c)
		return
	}
	log.Info("unzipping file...")
	_, err = fileutils.Unzip(destinationFilePath, config.Get().GetAbsTempDir(), os.FileMode(a.FileMode))
	if err != nil {
		log.Error(err)
		restoreStatus = interfaces.RestoreFailed
		ResponseHandler(nil, err, c)
		return
	}
	_ = os.RemoveAll(destinationFilePath)

	unzippedFolderPath := path.Join(config.Get().GetAbsTempDir(), utils.FileNameWithoutExtension(file.Filename))

	copySystemFiles := true // for example in macOS, we don't have systemd file & so to prevent that failure
	services := make([]string, 0)
	if _, err := os.Stat(constants.ServiceDir); errors.Is(err, os.ErrNotExist) {
		copySystemFiles = false
	}
	if copySystemFiles {
		services, _ = fileutils.ListFiles(path.Join(unzippedFolderPath, systemFolder))
		a.stopServices(services)
		err = utils.CopyDir(path.Join(unzippedFolderPath, systemFolder), constants.ServiceDir, "", 0)
		if err != nil {
			log.Error(err)
			restoreStatus = interfaces.RestoreFailed
			ResponseHandler(nil, err, c)
			return
		}
		err = a.SystemCtl.DaemonReload()
		if err != nil {
			log.Error(err)
			restoreStatus = interfaces.RestoreFailed
			ResponseHandler(nil, err, c)
			return
		}
	}
	rubixRegistryFile := path.Join(unzippedFolderPath, a.RubixRegistry.RubixRegistryDeviceInfoFile)
	rubixRegistryFileExist := false
	if _, err = os.Stat(rubixRegistryFile); !errors.Is(err, os.ErrNotExist) {
		rubixRegistryFileExist = true
	}
	if rubixRegistryFileExist {
		deviceInfo, err := a.RubixRegistry.GetDeviceInfo()
		if err != nil {
			log.Error(err)
			restoreStatus = interfaces.RestoreFailed
			ResponseHandler(nil, err, c)
			return
		}
		err = a.retainGlobalUUID(deviceInfo.GlobalUUID, rubixRegistryFile)
		if err != nil {
			log.Error(err)
			restoreStatus = interfaces.RestoreFailed
			ResponseHandler(nil, err, c)
			return
		}
	}

	err = utils.DeleteDir(path.Join(unzippedFolderPath, dataFolder), "", 0)
	if err != nil {
		log.Error(err)
		restoreStatus = interfaces.RestoreFailed
		ResponseHandler(nil, err, c)
		return
	}
	err = utils.CopyDir(path.Join(unzippedFolderPath, dataFolder), config.Get().GetSnapshotDir(), "", 0)
	if err != nil {
		restoreStatus = interfaces.RestoreFailed
		ResponseHandler(nil, err, c)
		return
	}
	err = os.RemoveAll(unzippedFolderPath)
	if err != nil {
		log.Errorf("failed to remove file %s", unzippedFolderPath)
	}
	if copySystemFiles {
		// Put nubeio-rubix-os.service file at last in order, self restart could also can happen
		index := 0
		exist := false
		serviceFile := namings.GetServiceNameFromAppName(constants.RubixOs)
		for i, str := range services {
			if str == serviceFile {
				exist = true
				index = i
				break
			}
		}
		if exist {
			services = append(services[:index], services[index+1:]...)
			services = append(services, serviceFile)
		}

		a.enableAndRestartServices(services)
	}
	log.Info("snapshot is restored")
	message := model.Message{Message: "snapshot is restored successfully"}
	restoreStatus = interfaces.Restored
	ResponseHandler(message, err, c)
}

func (a *SnapshotAPI) SnapshotStatus(c *gin.Context) {
	ResponseHandler(interfaces.SnapshotStatus{CreateStatus: createStatus, RestoreStatus: restoreStatus}, nil, c)
}

func (a *SnapshotAPI) stopServices(services []string) {
	log.Info("stopping services...")
	var wg sync.WaitGroup
	for _, service := range services {
		wg.Add(1)
		go func(service string) {
			defer wg.Done()
			if !utils.Contains(utils.ExcludedServices, service) {
				err := a.SystemCtl.Stop(service)
				if err != nil {
					log.Errorf("failed to stop service %s", service)
				}
			}
		}(service)
	}
	wg.Wait()
}

func (a *SnapshotAPI) enableAndRestartServices(services []string) {
	log.Info("enabling & restarting services")
	var wg sync.WaitGroup
	for _, service := range services {
		wg.Add(1)
		go func(service string) {
			defer wg.Done()
			if !utils.Contains(utils.ExcludedServices, service) {
				err := a.SystemCtl.Enable(service)
				if err != nil {
					log.Errorf("failed to enable service %s", service)
				}
				err = a.SystemCtl.Restart(service)
				if err != nil {
					log.Errorf("failed to restart service %s", service)
				}
			}
		}(service)
	}
	wg.Wait()
}

func (a *SnapshotAPI) retainGlobalUUID(globalUUID, rubixRegistryFile string) error {
	content, err := fileutils.ReadFile(rubixRegistryFile)
	if err != nil {
		return err
	}
	deviceInfoDefault := rubixregistry.DeviceInfoDefault{}
	err = json.Unmarshal([]byte(content), &deviceInfoDefault)
	if err != nil {
		return err
	}
	deviceInfoDefault.DeviceInfoFirstRecord.DeviceInfo.GlobalUUID = globalUUID
	deviceInfoDefaultRaw, err := json.Marshal(deviceInfoDefault)
	if err != nil {
		return err
	}
	err = os.WriteFile(rubixRegistryFile, deviceInfoDefaultRaw, os.FileMode(a.FileMode))
	if err != nil {
		return err
	}
	return nil
}
