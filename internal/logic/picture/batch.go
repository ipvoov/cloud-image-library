package picture

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/consts"
	"cloud/internal/dao"
	"cloud/internal/model/do"
	"cloud/internal/model/entity"
	"cloud/internal/service"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// UploadByBatch 批量上传图片
func (s *sPicture) UploadByBatch(ctx context.Context, req *v1.PictureUploadByBatchReq) (res *v1.PictureUploadByBatchRes, err error) {
	// 名称前缀默认等于搜索关键词
	namePrefix := req.NamePrefix
	if namePrefix == "" {
		namePrefix = req.SearchText
	}

	// 构建必应图片搜索URL
	fetchURL := fmt.Sprintf("https://cn.bing.com/images/async?q=%s&mmasync=1", url.QueryEscape(req.SearchText))

	g.Log().Infof(ctx, "开始批量抓取图片，搜索关键词: %s, 数量: %d", req.SearchText, req.Count)

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 发送HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "GET", fetchURL, nil)
	if err != nil {
		g.Log().Errorf(ctx, "创建HTTP请求失败: %v", err)
		return nil, gerror.New("创建请求失败")
	}

	// 设置请求头，模拟浏览器
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	httpReq.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	httpReq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := client.Do(httpReq)
	if err != nil {
		g.Log().Errorf(ctx, "获取页面失败: %v", err)
		return nil, gerror.New("获取页面失败")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		g.Log().Errorf(ctx, "HTTP状态码错误: %d", resp.StatusCode)
		return nil, gerror.New("获取页面失败")
	}

	// 解析HTML文档
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		g.Log().Errorf(ctx, "解析HTML失败: %v", err)
		return nil, gerror.New("解析页面失败")
	}

	// 查找图片容器
	dgControl := doc.Find(".dgControl").First()
	if dgControl.Length() == 0 {
		g.Log().Errorf(ctx, "未找到图片容器元素")
		return nil, gerror.New("获取元素失败")
	}

	// 查找所有图片元素
	imgElements := dgControl.Find("img.mimg")
	if imgElements.Length() == 0 {
		g.Log().Warningf(ctx, "未找到图片元素，尝试其他选择器")
		// 尝试其他可能的选择器
		imgElements = doc.Find("img[src*='bing.com']")
		if imgElements.Length() == 0 {
			return nil, gerror.New("未找到图片")
		}
	}

	g.Log().Infof(ctx, "找到 %d 个图片元素", imgElements.Length())

	uploadCount := 0
	imgElements.Each(func(i int, selection *goquery.Selection) {
		// 如果已经达到目标数量，停止处理
		if uploadCount >= req.Count {
			return // 停止整个Each循环
		}

		fileURL, exists := selection.Attr("src")
		if !exists || fileURL == "" {
			g.Log().Infof(ctx, "当前链接为空，已跳过: %s", fileURL)
			return // 跳过当前图片，继续下一张
		}

		// 处理图片地址，防止转义或者和对象存储冲突的问题
		// 移除URL参数部分
		if questionMarkIndex := strings.Index(fileURL, "?"); questionMarkIndex > -1 {
			fileURL = fileURL[:questionMarkIndex]
		}

		// 确保URL是完整的
		if strings.HasPrefix(fileURL, "//") {
			fileURL = "https:" + fileURL
		} else if strings.HasPrefix(fileURL, "/") {
			fileURL = "https://cn.bing.com" + fileURL
		}

		// 上传图片
		uploadReq := &v1.PictureUploadByUrlReq{
			FileUrl:  fileURL,
			FileName: namePrefix + gconv.String(uploadCount+1),
			SpaceId:  req.SpaceId,
		}

		_, uploadErr := s.UploadByUrl(ctx, uploadReq)
		if uploadErr != nil {
			g.Log().Errorf(ctx, "图片上传失败: %v, URL: %s", uploadErr, fileURL)
			return // 跳过当前图片，继续下一张
		}

		g.Log().Infof(ctx, "图片上传成功，当前已上传: %d 张", uploadCount+1)
		uploadCount++
	})

	g.Log().Infof(ctx, "批量上传完成，成功上传 %d 张图片", uploadCount)

	return &v1.PictureUploadByBatchRes{
		SuccessCount: uploadCount,
	}, nil
}

// EditByBatch 批量编辑图片
func (s *sPicture) EditByBatch(ctx context.Context, req *v1.PictureEditByBatchReq) (res *v1.PictureEditByBatchRes, err error) {
	// 验证用户权限
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	// 验证图片ID列表不能为空
	if len(req.PictureIdList) == 0 {
		return nil, gerror.New("图片ID列表不能为空")
	}

	g.Log().Infof(ctx, "开始批量编辑图片，用户ID: %d, 图片数量: %d", user.Id, len(req.PictureIdList))

	// 1. 验证图片ID列表的有效性和用户权限
	var pictures []*entity.Picture
	pic := dao.Picture.Columns()
	err = dao.Picture.Ctx(ctx).Where(pic.Id+" IN (?)", req.PictureIdList).
		Where(pic.IsDelete, 0).Scan(&pictures)
	if err != nil {
		g.Log().Errorf(ctx, "查询图片失败: %v", err)
		return nil, gerror.New("查询图片失败")
	}

	if len(pictures) == 0 {
		return nil, gerror.New("未找到有效的图片")
	}

	// 2. 验证用户权限：检查是否为管理员或空间创建者
	hasPermission := false
	if user.UserRole == consts.Admin {
		hasPermission = true
	} else {
		// 检查是否为空间创建者
		if req.SpaceId > 0 {
			hasPermission = s.hasSpacePermission(ctx, req.SpaceId, user)
		} else {
			// 如果没有指定空间ID，检查是否为所有图片的创建者
			allOwned := true
			for _, picture := range pictures {
				if picture.UserId != user.Id {
					allOwned = false
					break
				}
			}
			hasPermission = allOwned
		}
	}

	if !hasPermission {
		return nil, gerror.New("无权限批量编辑这些图片")
	}

	// 3. 使用事务进行批量更新
	successCount := 0
	err = dao.Picture.Ctx(ctx).Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for _, picture := range pictures {
			// 准备更新数据
			updateData := do.Picture{
				EditTime:     gtime.Now(),
				UpdateTime:   gtime.Now(),
				ReviewStatus: consts.DefRwStatus,
			}

			// 更新分类
			if req.Category != "" {
				updateData.Category = req.Category
			}

			// 更新标签
			if len(req.Tags) > 0 {
				tagsJson, _ := gjson.New(req.Tags).ToJson()
				updateData.Tags = string(tagsJson)
			}

			// 根据重命名规则更新名称
			if req.NameRule != "" {
				newName := fmt.Sprintf("%s_%d", req.NameRule, successCount+1)
				updateData.Name = newName
			}

			// 更新图片信息
			_, updateErr := dao.Picture.Ctx(ctx).Where(pic.Id, picture.Id).Data(updateData).Update()
			if updateErr != nil {
				g.Log().Errorf(ctx, "更新图片失败，ID: %d, 错误: %v", picture.Id, updateErr)
				return gerror.Newf("更新图片失败，ID: %d", picture.Id)
			}

			successCount++
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	g.Log().Infof(ctx, "批量编辑完成，成功处理了 %d 张图片", successCount)

	return &v1.PictureEditByBatchRes{
		Success: true,
	}, nil
}

// hasSpacePermission 检查用户是否有空间权限
func (s *sPicture) hasSpacePermission(ctx context.Context, spaceId int64, user *v1.GetLoginUserRes) bool {
	// 检查是否是空间创建者
	exist, err := dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, spaceId).
		Where(dao.Space.Columns().UserId, user.Id).
		Where(dao.Space.Columns().IsDelete, 0).Exist()
	if err != nil {
		g.Log().Errorf(ctx, "检查空间权限失败: %v", err)
		return false
	}

	return exist || user.UserRole == consts.Admin
}
