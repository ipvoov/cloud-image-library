package middleware

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Response 统一响应中间件
func Response(r *ghttp.Request) {
	r.Middleware.Next()

	// 如果已经有响应内容，直接返回
	if r.Response.BufferLength() > 0 {
		return
	}

	var (
		msg  = "操作成功"
		err  = r.GetError()
		res  = r.GetHandlerResponse()
		code = 0 // 成功用0，失败用非0
	)

	if err != nil {
		code = -1
		msg = err.Error()
		res = nil
	}

	r.Response.WriteJson(g.Map{
		"code":    code,
		"message": msg,
		"data":    res,
	})
}
