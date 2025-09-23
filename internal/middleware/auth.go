package middleware

import (
	"cloud/internal/consts"
	"cloud/internal/model/entity"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

// Auth 鉴权中间件
func Auth(r *ghttp.Request) {
	session := r.Session
	userObj, _ := session.Get(consts.LoginState)

	if userObj == nil {
		r.SetError(gerror.New("用户未登录"))
		return
	}

	r.Middleware.Next()
}

// AdminAuth 管理员鉴权中间件
func AdminAuth(r *ghttp.Request) {
	session := r.Session
	userObj, _ := session.Get(consts.LoginState)

	if userObj == nil {
		r.SetError(gerror.New("用户未登录"))
		return
	}
	var user *entity.User
	err := gconv.Struct(userObj.Map(), &user)
	if err != nil {
		r.SetError(gerror.New("权限校验不足"))
		return
	}

	if user.UserRole != "admin" {
		r.SetError(gerror.New("权限不足"))
		return
	}

	r.SetCtxVar("user", userObj)
	r.Middleware.Next()
}
