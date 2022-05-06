package main

import (
	"fmt"
	"github.com/NubeIO/flow-framework/api"
	min "github.com/NubeIO/flow-framework/plugin/nube/utils/backup/minio"
	"github.com/NubeIO/flow-framework/utils/directory"
	"github.com/NubeIO/flow-framework/utils/file"
	"github.com/NubeIO/flow-framework/utils/nstring"
	"reflect"
	"strings"
	"time"
)

func (inst *Instance) connection() error {
	host := "127.0.0.1:9000"
	if inst.config.Host != "" {
		host = inst.config.Host
	}
	accessKeyID := "12345678"
	if inst.config.AccessKeyID != "" {
		accessKeyID = inst.config.AccessKeyID
	}
	secretAccessKey := "12345678"
	if inst.config.SecretAccessKey != "" {
		secretAccessKey = inst.config.SecretAccessKey
	}
	conf := min.Configuration{
		Host:            host,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}
	if err := min.NewConnection(conf); err != nil {
		return err
	}
	inst.minioClient = min.GetClient()
	return nil
}

const dir = "flow-framework-tmp"

func (inst *Instance) makeDir() (path string, err error) {
	homeDir, err := directory.GetUserHomeDir()
	if err != nil {
		return "", err
	}
	path = fmt.Sprintf("%s/%s", homeDir, dir)
	err = directory.MakeDirIfNotExists(path)
	return path, err
}

func (inst *Instance) bucketName() string {

	bucketName := "flowframework"
	if inst.config.BucketName != "" {
		bucketName = inst.config.BucketName
	}
	info, err := inst.db.GetDeviceInfo()
	if err != nil {
		return ""
	}
	d := fmt.Sprintf("%s%s%s", info.ClientName, info.SiteName, info.DeviceName)
	if len(d) >= 10 {
		bucketName = d
	}
	noDashes := strings.Replace(bucketName, "_", " ", -1)
	noSpaceString := strings.ReplaceAll(noDashes, " ", "")
	n := nstring.NewString(noSpaceString)
	nn := n.RemoveSpecialCharacter()
	nn = n.ToLower()
	return nn
}

func (inst *Instance) makeJson(data interface{}, objectName string) error {
	bucketName := inst.bucketName()
	dirName, err := inst.makeDir()
	if err != nil {
		return err
	}
	f := fmt.Sprintf("%s/tmp.json", dirName)
	var o min.ObjectArgs
	o.BucketName = bucketName
	o.ObjectName = objectName
	o.FilePath = f

	d := reflect.ValueOf(data).Interface().(interface{})
	_, err = file.WriteDataToFileAsJSON(d, f)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = inst.minioClient.UploadObject(o, true)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func (inst *Instance) backNetworks() error {

	var arg api.Args
	arg.WithDevices = true
	arg.WithPoints = true

	nets, err := inst.db.GetNetworks(arg)
	if err != nil {
		return err
	}
	t := time.Now()
	file := fmt.Sprintf("networks_%s.json", t.Format("20060102150405"))
	err = inst.makeJson(nets, file)
	if err != nil {
		return err
	}

	return nil
}
