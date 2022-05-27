package min

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

var (
	minioClient = &minio.Client{}
	ctx         = context.Background()
)

// MinioClient minio client interface
type MinioClient interface {
	// CreateBucket args bucket, region string
	CreateBucket(o ObjectArgs) error
	// DeleteBucket args bucket string
	DeleteBucket(o ObjectArgs) error
	// ListBuckets args nil
	ListBuckets() ([]minio.BucketInfo, error)
	// BucketExists args bucket string
	BucketExists(o ObjectArgs) (bool, error)
	// GetBucketLifecycle args bucket string
	GetBucketLifecycle(o ObjectArgs) (*lifecycle.Configuration, error)
	// UploadObject o ObjectArgs, makeBucket bool
	UploadObject(o ObjectArgs, makeBucket bool) error
	// RemoveObject removes an object from a bucket. args BucketName string, ObjectName string, RemoveObjectOptions
	RemoveObject(o ObjectArgs) error
	// StatsObject verifies if object exists, and you have permission to access.
	StatsObject(o ObjectArgs) (info minio.ObjectInfo, err error)
	// DownloadObject downloads the file at path to the specified local path args o.BucketName, o.ObjectName, o.FilePath, o.GetObjectOptions
	DownloadObject(o ObjectArgs) error
}

// Configuration config minio for new connection
type Configuration struct {
	Host            string
	AccessKeyID     string
	SecretAccessKey string
	UseSSl          bool
}

// NewConnection new ftp connection
func NewConnection(config Configuration) (err error) {
	minioClient, err = minio.New(config.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSl,
	})
	if err != nil {
		return err
	}
	return nil
}

type client struct {
	client *minio.Client
}

// GetClient get minio client
func GetClient() MinioClient {
	return &client{
		client: minioClient,
	}
}
