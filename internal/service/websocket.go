// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	v1 "cloud/api/user/v1"
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	IWebSocket interface {
		// PictureEdit 添加用户
		PictureEdit(ctx context.Context, r *ghttp.Request, req *v1.WebSocketPictureEditReq)
	}
)

var (
	localWebSocket IWebSocket
)

func WebSocket() IWebSocket {
	if localWebSocket == nil {
		panic("implement not found for interface IWebSocket, forgot register?")
	}
	return localWebSocket
}

func RegisterWebSocket(i IWebSocket) {
	localWebSocket = i
}
