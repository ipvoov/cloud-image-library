package bucket

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/consts"
	"cloud/internal/service"
	"context"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func init() {
	service.RegisterBucket(New())
}

type sBucket struct {
	Cli *cos.Client
}

func New() *sBucket {
	ctx := context.Background()
	// 1. 检查文件是否存在
	u, _ := url.Parse(consts.BucketURL)
	b := &cos.BaseURL{BucketURL: u}
	cli := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  g.Cfg().MustGet(ctx, consts.SecretId).String(),
			SecretKey: g.Cfg().MustGet(ctx, consts.SecretKey).String(),
		},
	})
	return &sBucket{
		Cli: cli,
	}
}

// Upload 上传
func (s *sBucket) Upload(ctx context.Context, in *v1.BucketUploadReq) (res *v1.BucketUploadRes, err error) {
	fileName := in.File.Filename
	f, err := in.File.Open()
	if err != nil {
		g.Log().Errorf(ctx, "打开上传文件失败: %v", err)
		return nil, gerror.New("上传文件失败")
	}
	// 根据 spaceId 决定前缀：有 spaceId -> space/<id>/，否则 Public/
	var prefix string
	if in.SpaceId > 0 {
		prefix = "space/" + strings.TrimSpace(gconv.String(in.SpaceId)) + "/"
	} else {
		prefix = "Public/"
	}
	objectKey := prefix + fileName
	if _, err = s.Cli.Object.Put(ctx, objectKey, f, nil); err != nil {
		g.Log().Errorf(ctx, "上传到COS失败: %v", err)
		return nil, gerror.New("上传文件失败")
	}
	return &v1.BucketUploadRes{FileAddress: objectKey}, nil
}

// UploadByUrl URL上传
func (s *sBucket) UploadByUrl(ctx context.Context, in *v1.BucketUploadByUrlReq) (res *v1.BucketUploadByUrlRes, err error) {
	g.Log().Debugf(ctx, "URL上传: %s", in.FileUrl)

	// 从URL下载文件
	resp, err := http.Get(in.FileUrl)
	if err != nil {
		g.Log().Errorf(ctx, "下载文件失败: %v", err)
		return nil, gerror.New("下载文件失败")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		g.Log().Errorf(ctx, "下载文件失败，状态码: %d", resp.StatusCode)
		return nil, gerror.New("下载文件失败")
	}

	// 获取文件大小
	fileSize := resp.ContentLength
	if fileSize < 0 {
		// 如果Content-Length头不存在或为负数，设为0
		fileSize = 0
		g.Log().Warningf(ctx, "无法从Content-Length获取文件大小")
	}

	var fileName string
	if in.FileName == "" {
		// 生成文件名：日期 + UUID + 文件扩展名
		fileName = s.generateFileName(in.FileUrl)
	} else {
		fileName = in.FileName
	}

	// 根据 spaceId 决定前缀
	var prefix string
	if in.SpaceId > 0 {
		prefix = "space/" + strings.TrimSpace(gconv.String(in.SpaceId)) + "/"
	} else {
		prefix = "Public/"
	}
	// 上传到COS
	_, err = s.Cli.Object.Put(ctx, prefix+fileName, resp.Body, nil)
	if err != nil {
		g.Log().Errorf(ctx, "上传到COS失败: %v", err)
		return nil, gerror.New("上传文件失败")
	}

	return &v1.BucketUploadByUrlRes{
		FileAddress: prefix + fileName,
		FileName:    fileName,
		FileSize:    fileSize,
	}, nil
}

// Delete 删除
func (s *sBucket) Delete(ctx context.Context, in *v1.BucketDeleteReq) (out *v1.BucketDeleteRes, err error) {
	// 删除接口保持向后兼容：按传入的完整路径删除
	_, err = s.Cli.Object.Delete(ctx, in.FileName, nil)
	if err != nil {
		return nil, gerror.New("删除文件失败")
	}
	return &v1.BucketDeleteRes{}, nil
}

// generateFileName 生成文件名：日期 + UUID + 文件扩展名
func (s *sBucket) generateFileName(fileUrl string) string {
	// 获取当前日期，格式：20240320time.Now().Format("20060102")
	dateStr := gtime.Now().Format("Ymd")

	// 生成UUID（去掉连字符）
	uuid := strings.ReplaceAll(guid.S(), "-", "")
	uuid = uuid[:8]

	// 从URL中提取文件扩展名
	ext := s.GetFileExtFromUrl(fileUrl)

	// 组合文件名：日期_UUID.扩展名
	return dateStr + "_" + uuid + ext
}

// getFileExtFromUrl 从URL中提取文件扩展名
func (s *sBucket) GetFileExtFromUrl(fileUrl string) string {
	// 解析URL
	parsedUrl, err := url.Parse(fileUrl)
	if err != nil {
		return ".jpg" // 默认扩展名
	}

	// 获取路径部分的扩展名
	ext := path.Ext(parsedUrl.Path)

	// 如果没有扩展名，尝试从查询参数中获取
	if ext == "" {
		// 检查常见的图片格式参数
		query := parsedUrl.Query()
		if format := query.Get("format"); format != "" {
			ext = "." + format
		} else if format := query.Get("f"); format != "" {
			ext = "." + format
		}
	}

	// 如果仍然没有扩展名，使用默认值
	if ext == "" {
		ext = ".jpg"
	}

	// 确保扩展名是小写
	ext = strings.ToLower(ext)

	// 验证是否为支持的图片格式
	supportedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg"}
	for _, supportedExt := range supportedExts {
		if ext == supportedExt {
			return ext
		}
	}

	// 如果不是支持的格式，使用默认扩展名
	return ".jpg"
}
