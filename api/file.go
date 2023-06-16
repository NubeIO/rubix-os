package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/lib-files/fileutils"
	"github.com/NubeIO/rubix-os/interfaces"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"time"
)

type FileAPI struct {
	FileMode int
}

func (a *FileAPI) FileExists(c *gin.Context) {
	file := c.Query("file")
	exists := fileutils.FileExists(file)
	fileExistence := interfaces.FileExistence{File: file, Exists: exists}
	ResponseHandler(fileExistence, nil, c)
}

func (a *FileAPI) WalkFile(c *gin.Context) {
	path_ := c.Query("path")
	files := make([]string, 0)
	err := filepath.WalkDir(path_, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		files = append(files, p)
		return nil
	})
	ResponseHandler(files, err, c)
}

func (a *FileAPI) ListFiles(c *gin.Context) {
	files, err := a.listFiles(c.Query("path"))
	ResponseHandler(files, err, c)
}

func (a *FileAPI) listFiles(_path string) ([]fileutils.FileDetails, error) {
	fileInfo, err := os.Stat(_path)
	dirContent := make([]fileutils.FileDetails, 0)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(_path)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			dirContent = append(dirContent, fileutils.FileDetails{Name: file.Name(), IsDir: file.IsDir()})
		}
	} else {
		return nil, errors.New("it needs to be a directory, found a file")
	}
	return dirContent, nil
}

func (a *FileAPI) CreateFile(c *gin.Context) {
	file := c.Query("file")
	if file == "" {
		ResponseHandler(nil, errors.New("file can not be empty"), c)
		return
	}
	_, err := fileutils.CreateFile(file, os.FileMode(a.FileMode))
	ResponseHandler(interfaces.Message{Message: fmt.Sprintf("created file: %s", file)}, err, c)
}

func (a *FileAPI) CopyFile(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	if from == "" || to == "" {
		ResponseHandler(nil, errors.New("from and to names can not be empty"), c)
		return
	}
	err := fileutils.Copy(from, to)
	ResponseHandler(interfaces.Message{Message: "copied successfully"}, err, c)
}

func (a *FileAPI) RenameFile(c *gin.Context) {
	oldPath := c.Query("old_path")
	newPath := c.Query("new_path")
	if oldPath == "" || newPath == "" {
		ResponseHandler(nil, errors.New("old_path & new_path names can not be empty"), c)
		return
	}
	err := os.Rename(oldPath, newPath)
	ResponseHandler(interfaces.Message{Message: "renamed successfully"}, err, c)
}

func (a *FileAPI) MoveFile(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	if from == "" || to == "" {
		ResponseHandler(nil, errors.New("from and to names can not be empty"), c)
		return
	}
	if from == to {
		ResponseHandler(nil, errors.New("from and to names are same"), c)
		return
	}
	err := os.Rename(from, to)
	ResponseHandler(interfaces.Message{Message: "moved successfully"}, err, c)
}

func (a *FileAPI) DownloadFile(c *gin.Context) {
	path_ := c.Query("path")
	fileName := c.Query("file")
	c.FileAttachment(fmt.Sprintf("%s/%s", path_, fileName), fileName)
}

// UploadFile
// curl -X POST http://localhost:1661/api/files/upload?destination=/data/ -F "file=@/home/user/Downloads/bios-master.zip" -H "Content-Type: multipart/form-data"
func (a *FileAPI) UploadFile(c *gin.Context) {
	now := time.Now()
	destination := c.Query("destination")
	file, err := c.FormFile("file")
	resp := &interfaces.FileUploadResponse{}
	if err != nil || file == nil {
		ResponseHandler(resp, err, c)
		return
	}
	if found := fileutils.DirExists(destination); !found {
		ResponseHandler(nil, errors.New(fmt.Sprintf("destination not found %s", destination)), c)
		return
	}
	toFileLocation := path.Join(destination, filepath.Base(file.Filename))
	if err := c.SaveUploadedFile(file, toFileLocation); err != nil {
		ResponseHandler(resp, err, c)
		return
	}
	if err := os.Chmod(toFileLocation, os.FileMode(a.FileMode)); err != nil {
		ResponseHandler(resp, err, c)
		return
	}
	size, err := fileutils.GetFileSize(toFileLocation)
	if err != nil {
		ResponseHandler(resp, err, c)
		return
	}
	resp = &interfaces.FileUploadResponse{
		Destination: toFileLocation,
		File:        file.Filename,
		Size:        size.String(),
		UploadTime:  TimeTrack(now),
	}
	ResponseHandler(resp, nil, c)
}

func (a *FileAPI) ReadFile(c *gin.Context) {
	file := c.Query("file")
	if file == "" {
		ResponseHandler(nil, errors.New("file can not be empty"), c)
		return
	}
	found := fileutils.FileExists(file)
	if !found {
		ResponseHandler(nil, errors.New(fmt.Sprintf("file not found: %s", file)), c)
		return
	}
	c.File(file)
}

func (a *FileAPI) WriteFile(c *gin.Context) {
	file := c.Query("file")
	if file == "" {
		ResponseHandler(nil, errors.New("file can not be empty"), c)
		return
	}
	var m *interfaces.WriteFileData
	err := c.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	err = fileutils.WriteFile(file, m.Data, fs.FileMode(a.FileMode))
	ResponseHandler(interfaces.Message{Message: fmt.Sprintf("wrote the file: %s", file)}, err, c)
}

func (a *FileAPI) DeleteFile(c *gin.Context) {
	file := c.Query("file")
	if !fileutils.FileExists(file) {
		ResponseHandler(nil, errors.New(fmt.Sprintf("file doesn't exist: %s", file)), c)
		return
	}
	err := fileutils.Rm(file)
	ResponseHandler(interfaces.Message{Message: fmt.Sprintf("deleted file: %s", file)}, err, c)
}

func (a *FileAPI) DeleteAllFiles(c *gin.Context) {
	filePath := c.Query("path")
	if !fileutils.FileOrDirExists(filePath) {
		ResponseHandler(nil, errors.New(fmt.Sprintf("doesn't exist: %s", filePath)), c)
		return
	}
	err := fileutils.RemoveAllFiles(filePath)
	ResponseHandler(interfaces.Message{Message: fmt.Sprintf("deleted path: %s", filePath)}, err, c)
}

func TimeTrack(start time.Time) (out string) {
	elapsed := time.Since(start)
	// Skip this function, and fetch the PC and file for its parent.
	pc, _, _, _ := runtime.Caller(1)
	// Retrieve a function object this functions parent.
	funcObj := runtime.FuncForPC(pc)
	// Regex to extract just the function name (and not the module path).
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")
	out = fmt.Sprintf("%s took %s", name, elapsed)
	return out
}

func (a *FileAPI) WriteStringFile(c *gin.Context) {
	var m *interfaces.WriteFormatFile
	err := c.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	if m.FilePath == "" {
		ResponseHandler(nil, errors.New("file path can not be empty"), c)
		return
	}
	err = fileutils.WriteFile(m.FilePath, m.BodyAsString, fs.FileMode(a.FileMode))
	ResponseHandler(interfaces.Message{Message: fmt.Sprintf("wrote the file: %s", m.FilePath)}, err, c)
}

func (a *FileAPI) WriteFileYml(c *gin.Context) {
	var m *interfaces.WriteFormatFile
	err := c.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	if m.FilePath == "" {
		ResponseHandler(nil, errors.New("file path can not be empty"), c)
		return
	}
	data, err := yaml.Marshal(m.Body)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	err = ioutil.WriteFile(m.FilePath, data, fs.FileMode(a.FileMode))
	ResponseHandler(interfaces.Message{Message: fmt.Sprintf("wrote file: %s ok", m.FilePath)}, err, c)
}

func (a *FileAPI) WriteFileJson(c *gin.Context) {
	var m *interfaces.WriteFormatFile
	err := c.ShouldBindJSON(&m)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	if m.FilePath == "" {
		ResponseHandler(nil, errors.New("file path can not be empty"), c)
		return
	}
	data, err := json.Marshal(m.Body)
	if err != nil {
		ResponseHandler(nil, err, c)
		return
	}
	err = ioutil.WriteFile(m.FilePath, data, fs.FileMode(a.FileMode))
	ResponseHandler(interfaces.Message{Message: fmt.Sprintf("wrote file:%s ok", m.FilePath)}, err, c)
}
