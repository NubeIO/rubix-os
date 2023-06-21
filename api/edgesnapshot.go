package api

import (
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/NubeIO/rubix-os/config"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/NubeIO/rubix-os/src/cli/cligetter"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type EdgeSnapshotDatabase interface {
	GetSnapshotLog() ([]*model.SnapshotLog, error)
	CreateSnapshotLog(body *model.SnapshotLog) (*model.SnapshotLog, error)
	UpdateSnapshotLog(file string, body *model.SnapshotLog) (*model.SnapshotLog, error)
	DeleteSnapshotLog(file string) (*interfaces.Message, error)
	DeleteSnapshotLogs(files []string) (*interfaces.Message, error)

	CreateSnapshotCreateLog(body *model.SnapshotCreateLog) (*model.SnapshotCreateLog, error)
	UpdateSnapshotCreateLog(uuid string, body *model.SnapshotCreateLog) (*model.SnapshotCreateLog, error)

	CreateSnapshotRestoreLog(body *model.SnapshotRestoreLog) (*model.SnapshotRestoreLog, error)
	UpdateSnapshotRestoreLog(uuid string, body *model.SnapshotRestoreLog) (*model.SnapshotRestoreLog, error)

	ResolveHost(uuid string, name string) (*model.Host, error)
	GetLocationGroupHostNamesByHostUUID(hostUUID string) (*interfaces.LocationGroupHostName, error)
}
type EdgeSnapshotApi struct {
	DB       EdgeSnapshotDatabase
	FileMode int
}

func (a *EdgeSnapshotApi) GetSnapshots(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	cli := cligetter.GetEdgeBiosClient(host)
	arch, err := cli.GetArch()
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	snapshots, err := a.getSnapshots(arch.Arch)
	ResponseHandler(snapshots, err, ctx)
}

func (a *EdgeSnapshotApi) UpdateSnapshot(ctx *gin.Context) {
	body, err := getBodySnapshotLog(ctx)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	updateLog, err := a.DB.UpdateSnapshotLog(ctx.Params.ByName("file"), body)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	ResponseHandler(updateLog, err, ctx)
}

func (a *EdgeSnapshotApi) DeleteSnapshot(c *gin.Context) {
	file := c.Query("file")
	if file == "" {
		ResponseHandler(nil, errors.New("file can not be empty"), c)
		return
	}
	err := os.Remove(path.Join(config.Get().GetAbsSnapShotDir(), file))
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	_, err = a.DB.DeleteSnapshotLog(file)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	ResponseHandler(interfaces.Message{Message: fmt.Sprintf("deleted file: %s", file)}, err, c)
}

func (a *EdgeSnapshotApi) CreateSnapshot(ctx *gin.Context) {
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	body, err := getBodyCreateSnapshot(ctx)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	names, _ := a.DB.GetLocationGroupHostNamesByHostUUID(host.UUID)
	createLog, err := a.DB.CreateSnapshotCreateLog(&model.SnapshotCreateLog{UUID: "", HostUUID: host.UUID, Msg: "",
		Status: model.Creating, Description: body.Description, CreatedAt: time.Now()})
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	go func() {
		cli := cligetter.GetEdgeClient(host)
		snapshot, filename, err := cli.CreateSnapshot(names)
		if err == nil {
			err = os.WriteFile(path.Join(config.Get().GetAbsSnapShotDir(), filename), snapshot,
				os.FileMode(a.FileMode))
		}
		createLog.Status = model.Created
		createLog.Msg = filename
		if err != nil {
			createLog.Status = model.CreateFailed
			createLog.Msg = err.Error()
		}
		_, err = a.DB.UpdateSnapshotCreateLog(createLog.UUID, createLog)
		if err != nil {
			log.Error(err)
		}
		snapshotLog := model.SnapshotLog{
			File:        filename,
			Description: body.Description,
		}
		_, err = a.DB.CreateSnapshotLog(&snapshotLog)
		if err != nil {
			log.Error(err)
		}
		files, err := a.getSnapshotsFiles()
		if err != nil {
			log.Error(err)
			return
		}
		_, err = a.DB.DeleteSnapshotLogs(files)
		if err != nil {
			log.Error(err)
		}
	}()
	ResponseHandler(interfaces.Message{Message: "create snapshot process has submitted"}, nil, ctx)
}

func (a *EdgeSnapshotApi) RestoreSnapshot(ctx *gin.Context) {
	body, err := getBodyRestoreSnapshot(ctx)
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	if body.File == "" {
		ResponseHandler(nil, errors.New("file can not be empty"), ctx)
		return
	}
	host, err := a.DB.ResolveHost(matchHostUUIDName(ctx))
	if err != nil {
		ResponseHandler(nil, err, ctx)
		return
	}
	restoreLog, err := a.DB.CreateSnapshotRestoreLog(&model.SnapshotRestoreLog{UUID: "", HostUUID: host.UUID,
		Msg: "", Status: model.Restoring, Description: body.Description, CreatedAt: time.Now()})
	go func() {
		cli := cligetter.GetEdgeClient(host)
		reader, err := os.Open(path.Join(config.Get().GetAbsSnapShotDir(), body.File))
		if err == nil {
			err = cli.RestoreSnapshot(body.File, reader)
		}
		restoreLog.Status = model.Restored
		restoreLog.Msg = body.File
		if err != nil {
			restoreLog.Status = model.RestoreFailed
			restoreLog.Msg = err.Error()
		}
		_, _ = a.DB.UpdateSnapshotRestoreLog(restoreLog.UUID, restoreLog)
	}()
	ResponseHandler(interfaces.Message{Message: "restore snapshot process has submitted"}, nil, ctx)
}

func (a *EdgeSnapshotApi) DownloadSnapshot(c *gin.Context) {
	file := c.Query("file")
	c.FileAttachment(path.Join(config.Get().GetAbsSnapShotDir(), file), file)
}

func (a *EdgeSnapshotApi) UploadSnapshot(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil || file == nil {
		ResponseHandler(nil, err, c)
		return
	}
	description := c.Query("description")
	fileName := strings.ReplaceAll(file.Filename, " ", "")
	if path.Ext(fileName) != ".zip" {
		ResponseHandler(nil, errors.New("file is not a valid zip file"), c)
		return
	}
	toFileLocation := path.Join(config.Get().GetAbsSnapShotDir(), filepath.Base(fileName))
	if err := c.SaveUploadedFile(file, toFileLocation); err != nil {
		ResponseHandler(nil, err, c)
		return
	}

	snapshotLog := model.SnapshotLog{
		File:        file.Filename,
		Description: description,
	}
	_, err = a.DB.UpdateSnapshotLog(file.Filename, &snapshotLog)
	if err != nil {
		log.Error(err)
	}
	files, err := a.getSnapshotsFiles()
	if err != nil {
		log.Error(err)
		return
	}
	_, err = a.DB.DeleteSnapshotLogs(files)
	if err != nil {
		log.Error(err)
	}

	ResponseHandler(interfaces.Message{Message: "snapshot uploaded successfully"}, nil, c)
}

func (a *EdgeSnapshotApi) getSnapshots(arch string) ([]interfaces.Snapshots, error) {
	_path := config.Get().GetAbsSnapShotDir()
	fileInfo, err := os.Stat(_path)
	if err != nil {
		return nil, err
	}
	dirContent := make([]interfaces.Snapshots, 0)
	if fileInfo.IsDir() {
		snapshotLogs, err := a.DB.GetSnapshotLog()
		files, err := ioutil.ReadDir(_path)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			fileParts := strings.Split(file.Name(), "_")
			archParts := fileParts[len(fileParts)-1]
			archFromSnapshot := strings.Split(archParts, ".")[0]
			if archFromSnapshot == arch {
				description := ""
				for _, snapshotLog := range snapshotLogs {
					if snapshotLog.File == file.Name() {
						description = snapshotLog.Description
						break
					}
				}
				dirContent = append(dirContent, interfaces.Snapshots{
					Name:        file.Name(),
					Size:        file.Size(),
					CreatedAt:   file.ModTime(),
					Description: description,
				})
			}
		}
	} else {
		return nil, errors.New("it needs to be a directory, found a file")
	}
	return dirContent, nil
}

func (a *EdgeSnapshotApi) getSnapshotsFiles() ([]string, error) {
	_path := config.Get().GetAbsSnapShotDir()
	fileInfo, err := os.Stat(_path)
	if err != nil {
		return nil, err
	}
	outputFiles := make([]string, 0)
	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(_path)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			outputFiles = append(outputFiles, file.Name())
		}
	} else {
		return nil, errors.New("it needs to be a directory, found a file")
	}
	return outputFiles, nil
}
