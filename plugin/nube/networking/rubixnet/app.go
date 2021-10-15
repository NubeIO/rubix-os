package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	min "github.com/NubeDev/flow-framework/plugin/nube/utils/backup/minio"
	"github.com/NubeDev/flow-framework/utils"
	"reflect"
	"strings"
	"time"
)

func (i *Instance) connection() error {
	host := "127.0.0.1:9000"
	if i.config.Host != "" {
		host = i.config.Host
	}
	accessKeyID := "12345678"
	if i.config.AccessKeyID != "" {
		accessKeyID = i.config.AccessKeyID
	}
	secretAccessKey := "12345678"
	if i.config.SecretAccessKey != "" {
		secretAccessKey = i.config.SecretAccessKey
	}
	conf := min.Configuration{
		Host:            host,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}
	if err := min.NewConnection(conf); err != nil {
		return err
	}
	i.minioClient = min.GetClient()
	return nil
}

const dir = "flow-framework-tmp"

func (i *Instance) makeDir() (path string, err error) {
	homeDir, err := utils.GetUserHomeDir()
	if err != nil {
		return "", err
	}
	path = fmt.Sprintf("%s/%s", homeDir, dir)
	err = utils.MakeDirIfNotExists(path)
	return path, err
}

func (i *Instance) bucketName() string {

	bucketName := "flowframework"
	if i.config.BucketName != "" {
		bucketName = i.config.BucketName
	}
	info, err := i.db.GetDeviceInfo()
	if err != nil {
		return ""
	}
	d := fmt.Sprintf("%s%s%s", info.ClientName, info.SiteName, info.DeviceName)
	if len(d) >= 10 {
		bucketName = d
	}
	noDashes := strings.Replace(bucketName, "_", " ", -1)
	noSpaceString := strings.ReplaceAll(noDashes, " ", "")
	n := utils.NewString(noSpaceString)
	nn := n.RemoveSpecialCharacter()
	nn = n.ToLower()
	return nn
}

func (i *Instance) makeJson(data interface{}, objectName string) error {
	bucketName := i.bucketName()
	dirName, err := i.makeDir()
	if err != nil {
		return err
	}
	file := fmt.Sprintf("%s/tmp.json", dirName)
	var o min.ObjectArgs
	o.BucketName = bucketName
	o.ObjectName = objectName
	o.FilePath = file

	d := reflect.ValueOf(data).Interface().(interface{})
	_, err = utils.WriteDataToFileAsJSON(d, file)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = i.minioClient.UploadObject(o, true)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func (i *Instance) backNetworks() error {

	var arg api.Args
	arg.WithDevices = true
	arg.WithPoints = true
	arg.WithSerialConnection = true
	arg.WithIpConnection = true

	nets, err := i.db.GetNetworks(arg)
	if err != nil {
		return err
	}
	t := time.Now()
	file := fmt.Sprintf("networks_%s.json", t.Format("20060102150405"))
	err = i.makeJson(nets, file)
	if err != nil {
		return err
	}

	return nil
}
