// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	v1 "cloud/api/user/v1"
	"context"
)

type (
	IBucket interface {
		// Upload 上传
		Upload(ctx context.Context, in *v1.BucketUploadReq) (res *v1.BucketUploadRes, err error)
		// UploadByUrl URL上传
		UploadByUrl(ctx context.Context, in *v1.BucketUploadByUrlReq) (res *v1.BucketUploadByUrlRes, err error)
		// Delete 删除
		Delete(ctx context.Context, in *v1.BucketDeleteReq) (out *v1.BucketDeleteRes, err error)
		// getFileExtFromUrl 从URL中提取文件扩展名
		GetFileExtFromUrl(fileUrl string) string
	}
)

var (
	localBucket IBucket
)

func Bucket() IBucket {
	if localBucket == nil {
		panic("implement not found for interface IBucket, forgot register?")
	}
	return localBucket
}

func RegisterBucket(i IBucket) {
	localBucket = i
}
