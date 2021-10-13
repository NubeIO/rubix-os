package utils

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
	"os/user"
	"path/filepath"
)

func DirRemoveContents(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func DirIsWritable(path string) bool {
	return unix.Access(path, unix.W_OK) == nil
}

//DirChangePermissions
/*
@param path /etc
@param permissions 0700
*/
func DirChangePermissions(path string, permissions uint32) (ok bool, err error) {
	err = os.Chmod(path, os.FileMode(permissions))
	if err != nil {
		log.Error(err)
		return false, err
	}
	return true, err
}

func DirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func MakeDirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModeDir|0755)
	}
	return nil
}

func GetUserHomeDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.HomeDir, err
}
