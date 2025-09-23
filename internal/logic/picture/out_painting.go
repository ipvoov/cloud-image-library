package picture

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/consts"
	"cloud/internal/dao"
	"cloud/internal/model/entity"
	"cloud/internal/service"
	"context"
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

// CreateOutPainting 创建扩图
func (s *sPicture) CreateOutPainting(ctx context.Context, req *v1.CreatePictureOutPaintingReq) (res *v1.CreatePictureOutPaintingRes, err error) {
	// 1. 验证用户权限
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	// 2. 验证图片是否存在
	var picture *entity.Picture
	pic := dao.Picture.Columns()
	err = dao.Picture.Ctx(ctx).Where(pic.Id, req.PictureId).
		Where(pic.IsDelete, 0).Scan(&picture)
	if err != nil || picture == nil {
		return nil, gerror.New("图片不存在")
	}

	// 3. 验证用户权限：只有图片创建者或管理员可以扩图
	if picture.UserId != user.Id && user.UserRole != "admin" {
		return nil, gerror.New("无权限对此图片进行扩图")
	}

	// 4. 调用火山引擎扩图API
	client := arkruntime.NewClientWithApiKey(g.Cfg().MustGet(ctx, consts.AiKey).String())
	if client == nil {
		return nil, gerror.New("AI服务初始化失败，请检查配置")
	}

	// 构建扩图请求
	responseFormat := model.GenerateImagesResponseFormatURL
	size := model.GenerateImagesSizeAdaptive
	generateReq := model.GenerateImagesRequest{
		Model:          "doubao-seededit-3-0-i2i-250628",
		Prompt:         req.Prompt,
		Image:          picture.Url,
		ResponseFormat: &responseFormat,
		Seed:           volcengine.Int64(123),
		GuidanceScale:  volcengine.Float64(5.5),
		Size:           &size,
		Watermark:      volcengine.Bool(true),
	}

	g.Log().Infof(ctx, "开始扩图，用户ID: %d, 图片ID: %d, 描述: %s", user.Id, req.PictureId, req.Prompt)

	// 调用扩图API
	imagesResponse, err := client.GenerateImages(ctx, generateReq)
	if err != nil {
		g.Log().Errorf(ctx, "扩图失败: %v", err)
		return &v1.CreatePictureOutPaintingRes{
			Success: false,
			Message: fmt.Sprintf("扩图失败: %v", err),
		}, nil
	}

	if len(imagesResponse.Data) == 0 {
		return &v1.CreatePictureOutPaintingRes{
			Success: false,
			Message: "扩图返回结果为空",
		}, nil
	}

	outputImageUrl := *imagesResponse.Data[0].Url
	g.Log().Infof(ctx, "扩图成功，用户ID: %d, 图片ID: %d, 结果URL: %s", user.Id, req.PictureId, outputImageUrl)

	return &v1.CreatePictureOutPaintingRes{
		OutputImageUrl: outputImageUrl,
		Success:        true,
		Message:        "扩图成功",
	}, nil
}
