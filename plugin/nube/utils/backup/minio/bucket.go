package min

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

func (m *client) CreateBucket(o ObjectArgs) error {
	ok, err := m.client.BucketExists(ctx, o.BucketName)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return m.client.MakeBucket(ctx, o.BucketName, minio.MakeBucketOptions{
		Region: o.Region,
	})
}

func (m *client) DeleteBucket(o ObjectArgs) error {
	return m.client.RemoveBucket(ctx, o.BucketName)
}

func (m *client) ListBuckets() ([]minio.BucketInfo, error) {
	return m.client.ListBuckets(ctx)
}

func (m *client) BucketExists(o ObjectArgs) (bool, error) {
	return m.client.BucketExists(ctx, o.BucketName)
}

func (m *client) GetBucketLifecycle(o ObjectArgs) (*lifecycle.Configuration, error) {
	return m.client.GetBucketLifecycle(ctx, o.BucketName)
}
