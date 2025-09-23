package v1

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// BucketUploadReq 上传请求
type BucketUploadReq struct {
	File    *ghttp.UploadFile `json:"file" type:"file" dc:"文件"`
	SpaceId int64             `json:"spaceId" dc:"空间ID（>0 时上传到 space 目录，否则 Public）"`
}

// BucketUploadRes 上传响应
type BucketUploadRes struct {
	FileAddress string `json:"fileAddress"`
}

// BucketDownloadlReq 下载请求
type BucketDownloadlReq struct {
	FileUrl string `json:"fileUrl"`
}

// BucketDownloadlRes 下载响应
type BucketDownloadlRes struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// BucketUploadByUrlReq URL上传请求
type BucketUploadByUrlReq struct {
	FileUrl  string `json:"fileUrl" v:"required#文件URL不能为空"`
	FileName string `json:"fileName"`
	SpaceId  int64  `json:"spaceId"`
	Id       int64  `json:"id"`
}

// BucketUploadByUrlRes URL上传响应
type BucketUploadByUrlRes struct {
	FileAddress string `json:"fileAddress"`
	FileName    string `json:"fileName"`
	FileSize    int64  `json:"fileSize"` // 文件大小（字节）
}

// BucketDeleteReq 删除请求
type BucketDeleteReq struct {
	FileName string `json:"fileName"`
}

// BucketDeleteRes 删除响应
type BucketDeleteRes struct {
}
