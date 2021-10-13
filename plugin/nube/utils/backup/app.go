package main

import (
	"fmt"
	"github.com/NubeDev/flow-framework/api"
	min "github.com/NubeDev/flow-framework/plugin/nube/utils/backup/minio"
	"github.com/NubeDev/flow-framework/utils"
	"reflect"
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

func (i *Instance) makeDir() {
	dir := ""
	err := utils.MakeDirIfNotExists(dir)
	if err != nil {
		//return err
	}
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
	bucketName := "images"
	objectName := "test.json"
	////objectName := "profile.jpeg"
	filePath := "/tmp/test.json"
	//filepath := "/home/aidan/code/go/tests/minio/go-minio/assets/images/profile.jpeg"
	var o min.ObjectArgs
	o.BucketName = bucketName
	o.ObjectName = objectName
	o.FilePath = filePath

	data := reflect.ValueOf(nets).Interface().(interface{})
	_, err = utils.WriteDataToFileAsJSON(data, filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = i.minioClient.UploadObject(o, false)
	if err != nil {
		fmt.Println(err, 222)
		return err
	}

	return nil
}
