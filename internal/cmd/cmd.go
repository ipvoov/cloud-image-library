package cmd

import (
	"cloud/internal/controller"
	"cloud/internal/middleware"
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			// 配置跨域
			s.Use(ghttp.MiddlewareCORS)
			// 统一响应中间件
			// s.Use(middleware.Response)
			s.Use(ghttp.MiddlewareHandlerResponse)

			// API路由组
			s.Group("/api", func(group *ghttp.RouterGroup) {
				// 用户相关路由
				group.Group("/user", func(group *ghttp.RouterGroup) {
					// 不需要登录的接口
					group.POST("/register", controller.User.Register)
					group.POST("/login", controller.User.Login)

					// 需要登录的接口
					group.Group("/", func(group *ghttp.RouterGroup) {
						group.Middleware(middleware.Auth)
						group.GET("/get/login", controller.User.GetLoginUser)
						group.POST("/logout", controller.User.Logout)
						group.GET("/get", controller.User.GetById)
						group.GET("/profile", controller.User.GetProfile)
						group.POST("/profile/update", controller.User.UpdateProfile)

					})

					// 需要管理员权限的接口
					group.Group("/", func(group *ghttp.RouterGroup) {
						group.Middleware(middleware.AdminAuth)
						group.POST("/add", controller.User.Add)
						group.POST("/update", controller.User.Update)
						group.POST("/delete", controller.User.Delete)
						group.POST("/list/page/vo", controller.User.ListByPage)
					})
				})

				group.Group("/picture", func(group *ghttp.RouterGroup) {
					// 上传相关
					group.Group("/upload", func(group *ghttp.RouterGroup) {
						group.Middleware(middleware.Auth)
						group.POST("/", controller.Picture.Upload)
						group.POST("/url", controller.Picture.UploadByUrl)
						group.POST("/batch", controller.Picture.UploadByBatch)
					})
					// 下载
					// group.Group("/download", func(group *ghttp.RouterGroup) {
					// 	group.POST("/", controller.File.Download)
					// })
					// 获取图片标签分类
					group.GET("/tag_category", controller.Picture.TagCategory)
					// 图片CRUD操作
					group.POST("/edit", controller.Picture.Edit)
					group.Group("/edit", func(group *ghttp.RouterGroup) {
						group.Middleware(middleware.Auth)
						group.POST("/batch", controller.Picture.EditByBatch)
					})
					group.POST("/update", controller.Picture.Update)
					group.POST("/delete", controller.Picture.Delete)
					group.GET("/get", controller.Picture.Get)
					group.GET("/get/vo", controller.Picture.GetVO)
					// 分页查询
					group.POST("/list/page", controller.Picture.ListByPage)
					group.POST("/list/page/vo", controller.Picture.ListVOByPage)
					// 缓存版本的分页查询（暂时指向同一个方法）
					group.POST("/list/page/vo/cache", controller.Picture.ListVOByPage)

					// 搜索功能
					group.Group("/search", func(group *ghttp.RouterGroup) {
						group.Middleware(middleware.Auth)
						group.POST("/picture", controller.Picture.SearchByPicture)
						group.POST("/color", controller.Picture.SearchByColor)
					})

					//AI扩图功能
					group.Group("/out_painting", func(group *ghttp.RouterGroup) {
						group.Middleware(middleware.Auth)
						group.POST("/create", controller.Picture.CreateOutPainting)
					})

					// 管理员功能
					group.Group("/", func(group *ghttp.RouterGroup) {
						group.Middleware(middleware.AdminAuth)
						// 图片审核
						group.POST("/review", controller.Picture.Review)
					})
				})

				// 空间相关路由
				group.Group("/space", func(group *ghttp.RouterGroup) {
					// 需要登录的接口
					group.Group("/", func(group *ghttp.RouterGroup) {
						group.Middleware(middleware.Auth)
						group.POST("/add", controller.Space.Add)
						group.POST("/edit", controller.Space.Edit)
						group.POST("/delete", controller.Space.Delete)
						group.GET("/get", controller.Space.Get)
						group.GET("/get/vo", controller.Space.GetVO)
						group.POST("/list/page", controller.Space.ListByPage)
						group.POST("/list/page/vo", controller.Space.ListVOByPage)
						group.GET("/list/level", controller.Space.ListLevel)
					})

					// 管理员功能
					group.Group("/", func(group *ghttp.RouterGroup) {
						group.Middleware(middleware.AdminAuth)
						group.POST("/update", controller.Space.Update)
					})

					// 空间分析相关路由
					group.Group("/analyze", func(group *ghttp.RouterGroup) {
						group.Middleware(middleware.Auth)
						group.POST("/category", controller.SpaceAnalyze.CategoryAnalyze)
						group.POST("/tag", controller.SpaceAnalyze.TagAnalyze)
						group.POST("/size", controller.SpaceAnalyze.SizeAnalyze)
						group.POST("/usage", controller.SpaceAnalyze.UsageAnalyze)
						group.POST("/user", controller.SpaceAnalyze.UserAnalyze)
						group.POST("/rank", controller.SpaceAnalyze.RankAnalyze)
					})
				})

				// 空间用户相关路由
				group.Group("/spaceUser", func(group *ghttp.RouterGroup) {
					// 需要登录的接口
					group.Group("/", func(group *ghttp.RouterGroup) {
						group.Middleware(middleware.Auth)
						group.POST("/add", controller.SpaceUser.Add)
						group.POST("/edit", controller.SpaceUser.Edit)
						group.POST("/delete", controller.SpaceUser.Delete)
						group.POST("/get", controller.SpaceUser.Get)
						group.POST("/list", controller.SpaceUser.List)
						group.POST("/list/my", controller.SpaceUser.ListMy)
					})
				})

				// WebSocket 相关路由
				group.Group("/ws", func(group *ghttp.RouterGroup) {
					// 图片协同编辑 WebSocket
					group.GET("/picture/edit", controller.WebSocket.PictureEdit)
					// WebSocket 测试端点
					group.GET("/test", controller.WebSocket.Test)
				})
			})
			s.Run()
			return nil
		},
	}
)
