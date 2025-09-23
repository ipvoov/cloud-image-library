package picture

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/dao"
	"cloud/internal/model/entity"
	"cloud/internal/service"
	"context"
	"fmt"
	"os"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

// CreateAIEditingTask 创建AI编辑任务
func (s *sPicture) CreateAIEditingTask(ctx context.Context, req *v1.CreatePictureAIEditingTaskReq) (res *v1.CreatePictureAIEditingTaskRes, err error) {
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

	// 3. 验证用户权限：只有图片创建者或管理员可以编辑
	if picture.UserId != user.Id && user.UserRole != "admin" {
		return nil, gerror.New("无权限编辑此图片")
	}

	// 4. 调用火山引擎AI编辑API
	client := arkruntime.NewClientWithApiKey(os.Getenv("ARK_API_KEY"))
	if client == nil {
		return nil, gerror.New("AI服务初始化失败，请检查配置")
	}

	// 构建AI编辑请求
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

	g.Log().Infof(ctx, "开始AI编辑任务，用户ID: %d, 图片ID: %d, 描述: %s", user.Id, req.PictureId, req.Prompt)

	// 调用AI编辑API
	imagesResponse, err := client.GenerateImages(ctx, generateReq)
	if err != nil {
		g.Log().Errorf(ctx, "AI编辑任务失败: %v", err)
		return nil, gerror.Newf("AI编辑任务失败: %v", err)
	}

	if len(imagesResponse.Data) == 0 {
		return nil, gerror.New("AI编辑任务返回结果为空")
	}

	// 生成任务ID（这里简化处理，实际应该使用更复杂的任务管理系统）
	taskId := fmt.Sprintf("ai_edit_%d_%d", user.Id, req.PictureId)

	g.Log().Infof(ctx, "AI编辑任务创建成功，任务ID: %s, 结果URL: %s", taskId, *imagesResponse.Data[0].Url)

	return &v1.CreatePictureAIEditingTaskRes{
		Output: &v1.CreateAIEditingTaskResponse{
			TaskId: taskId,
		},
		RequestId: "req_" + taskId,
	}, nil
}

// GetAIEditingTask 获取AI编辑任务
func (s *sPicture) GetAIEditingTask(ctx context.Context, req *v1.GetPictureAIEditingTaskReq) (res *v1.GetPictureAIEditingTaskRes, err error) {
	// 这里简化处理，实际应该从任务管理系统或缓存中获取任务状态
	// 为了演示，我们直接返回成功状态

	g.Log().Infof(ctx, "获取AI编辑任务状态，任务ID: %s", req.TaskId)

	// 模拟任务状态查询
	// 实际实现中，这里应该：
	// 1. 从Redis或数据库中查询任务状态
	// 2. 如果任务完成，返回结果图片URL
	// 3. 如果任务失败，返回错误信息

	// 这里为了演示，我们假设任务总是成功的
	// 实际实现需要根据任务ID查询真实的状态

	return &v1.GetPictureAIEditingTaskRes{
		Output: &v1.GetAIEditingTaskResponse{
			TaskId:         req.TaskId,
			TaskStatus:     "SUCCEEDED",
			OutputImageUrl: "https://example.com/edited_image.jpg", // 这里应该是真实的编辑结果URL
			ErrorMessage:   "",
		},
		RequestId: "req_" + req.TaskId,
	}, nil
}
