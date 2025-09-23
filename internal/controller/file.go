package controller

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/service"
	"context"
	"io"
	"net/http"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

var File = cFile{}

type cFile struct{}

func (f *cFile) Upload(ctx context.Context, req *v1.BucketUploadReq) (res *v1.BucketUploadRes, err error) {
	if req.File == nil {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "文件不能为空")
	}
	return service.Bucket().Upload(ctx, req)
}

func (f *cFile) UploadByUrl(ctx context.Context, req *v1.BucketUploadByUrlReq) (res *v1.BucketUploadByUrlRes, err error) {
	if req.FileUrl == "" {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "文件URL不能为空")
	}
	return service.Bucket().UploadByUrl(ctx, req)
}

func (f *cFile) Download(ctx context.Context, in *v1.BucketDownloadlReq) (out *v1.BucketDownloadlRes, err error) {
	g.Log().Debugf(ctx, in.FileUrl)
	r := ghttp.RequestFromCtx(ctx)

	// 获取前端传递的文件URL参数
	fileURL := in.FileUrl
	if fileURL == "" {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "URL参数不能为空")
	}

	// 创建HTTP客户端请求
	client := &http.Client{}
	req, err := http.NewRequest("GET", fileURL, nil)
	if err != nil {
		glog.Error(r.Context(), "创建请求失败:", err)
		return nil, gerror.NewCode(gcode.CodeInternalError, "创建请求失败")
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		glog.Error(r.Context(), "请求远程文件失败:", err)
		return nil, gerror.NewCode(gcode.CodeInternalError, "无法获取远程文件")
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		glog.Error(r.Context(), "远程服务器返回错误状态码:", resp.StatusCode)
		return nil, gerror.NewCode(gcode.CodeInternalError, "远程服务器返回错误")
	}

	// 设置响应头
	r.Response.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		r.Response.Header().Set("Content-Length", contentLength)
	}

	// 设置下载文件名（可选）
	// filename := "downloaded_file"
	// r.Response.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

	// 将远程文件内容流式传输到响应
	_, err = io.Copy(r.Response.Writer, resp.Body)
	if err != nil {
		glog.Error(r.Context(), "传输数据失败:", err)
		return nil, gerror.NewCode(gcode.CodeInternalError, "传输数据失败")
	}

	// 对于文件下载，返回空结构体表示成功
	return &v1.BucketDownloadlRes{}, nil
}
