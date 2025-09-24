package controller

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/service"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

var WebSocket = cWebSocket{}

type cWebSocket struct{}

func (c *cWebSocket) WebSocketPictureEdit(ctx context.Context, req *v1.WebSocketPictureEditReq) (res *v1.WebSocketPictureEditRes, err error) {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return nil, gerror.New("无法获取请求对象")
	}
	service.WebSocket().PictureEdit(ctx, r, req)
	return &v1.WebSocketPictureEditRes{}, nil
}
