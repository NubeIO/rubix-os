package min

import (
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

type ObjectArgs struct {
	BucketName          string
	ObjectName          string
	FilePath            string
	Region              string
	PutObjectOptions    minio.PutObjectOptions
	MakeBucketOptions   minio.MakeBucketOptions
	StatObjectOptions   minio.StatObjectOptions
	GetObjectOptions    minio.GetObjectOptions
	RemoveObjectOptions minio.RemoveObjectOptions
}

// UploadObject upload image
func (m *client) UploadObject(o ObjectArgs, makeBucket bool) error {
	if makeBucket {
		err := m.client.MakeBucket(ctx, o.BucketName, o.MakeBucketOptions)
		if err != nil {
			exists, errBucketExists := m.client.BucketExists(ctx, o.BucketName)
			if errBucketExists != nil {
				logrus.Errorf("[UploadObject] check bucket exists error: %s", err)
				return err
			}
			if !exists {
				logrus.Errorf("[UploadObject] make bucket error: %s", err)
				return err
			}
		}
	}
	_, err := m.client.FPutObject(ctx, o.BucketName, o.ObjectName, o.FilePath, o.PutObjectOptions)
	if err != nil {
		logrus.Errorf("[UploadObject] put object error: %s", err)
		return err
	}
	return nil
}

func (m *client) StatsObject(o ObjectArgs) (info minio.ObjectInfo, err error) {
	info, err = m.client.StatObject(ctx, o.BucketName, o.ObjectName, o.StatObjectOptions)
	return info, err
}

func (m *client) RemoveObject(o ObjectArgs) error {
	return m.client.RemoveObject(ctx, o.BucketName, o.ObjectName, o.RemoveObjectOptions)
}

func (m *client) DownloadObject(o ObjectArgs) error {
	return m.client.FGetObject(ctx, o.BucketName, o.ObjectName, o.FilePath, o.GetObjectOptions)
}
