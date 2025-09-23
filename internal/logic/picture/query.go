package picture

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/consts"
	"cloud/internal/dao"
	"cloud/internal/model/entity"
	"cloud/internal/service"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
)

// TagCategory 获取图片标签分类
func (s *sPicture) TagCategory(ctx context.Context, req *v1.PictureTagCategoryReq) (res *v1.PictureTagCategoryRes, err error) {
	// 返回常用的标签列表
	tagList := []string{
		// 主题类
		"风景", "人物", "动物", "建筑", "美食", "花卉", "植物",
		"城市", "自然", "海洋", "山川", "天空", "日落", "夜景",

		// 风格类
		"摄影", "插画", "设计", "艺术", "抽象", "复古", "现代",
		"简约", "文艺", "清新", "唯美", "梦幻", "科幻", "卡通",

		// 色彩类
		"黑白", "彩色", "暖色调", "冷色调", "高对比", "柔和",

		// 用途类
		"壁纸", "头像", "封面", "背景", "素材", "图标", "logo",
		"海报", "banner", "名片", "宣传", "广告",

		// 情感类
		"温馨", "浪漫", "激情", "宁静", "活力", "神秘", "优雅",

		// 技术类
		"高清", "4K", "矢量", "手绘", "数字艺术", "3D", "渲染",
	}

	// 返回常用的分类列表
	categoryList := []string{
		"摄影作品",
		"数字艺术",
		"插画设计",
		"平面设计",
		"UI设计",
		"网页设计",
		"品牌设计",
		"包装设计",
		"海报设计",
		"图标素材",
		"背景纹理",
		"矢量图形",
		"手绘作品",
		"3D渲染",
		"概念艺术",
		"游戏美术",
		"动漫插画",
		"儿童插画",
		"时尚摄影",
		"产品摄影",
		"建筑摄影",
		"风光摄影",
		"人像摄影",
		"街拍摄影",
		"其他",
	}

	return &v1.PictureTagCategoryRes{
		TagList:      tagList,
		CategoryList: categoryList,
	}, nil
}

// Get 获取图片详情
func (s *sPicture) Get(ctx context.Context, req *v1.PictureGetReq) (res *v1.PictureAdminGetRes, err error) {
	// 从数据库查询图片信息
	var picture *entity.Picture
	pic := dao.Picture.Columns()
	db := dao.Picture.Ctx(ctx).Where(pic.Id, req.Id).
		Where(pic.IsDelete, 0)
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户不存在")
	}
	err = db.Scan(&picture)
	if err != nil || picture == nil {
		return nil, gerror.New("图片不存在")
	}
	// 权限检查
	hasPermission := false

	// 1. 管理员有查看所有图片的权限
	if user.UserRole == consts.Admin {
		hasPermission = true
	}

	// 2. 图片所有者有查看权限
	if picture.UserId == user.Id {
		hasPermission = true
	}

	// 3. 已审核通过的公共图片（无空间）可以查看
	if picture.SpaceId == 0 && picture.ReviewStatus == 1 {
		hasPermission = true
	}

	// 4. 如果图片属于某个空间，检查用户是否是空间成员
	if picture.SpaceId > 0 {
		// 检查用户是否是空间所有者
		var space *entity.Space
		err = dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, picture.SpaceId).
			Where(dao.Space.Columns().IsDelete, 0).Scan(&space)
		if err == nil && space != nil && space.UserId == user.Id {
			hasPermission = true
		}

		// 检查用户是否是空间成员
		if !hasPermission {
			exists, checkErr := dao.SpaceUser.Ctx(ctx).Where(dao.SpaceUser.Columns().SpaceId, picture.SpaceId).
				Where(dao.SpaceUser.Columns().UserId, user.Id).Exist()
			if checkErr == nil && exists {
				hasPermission = true
			}
		}
	}

	if !hasPermission {
		return nil, gerror.New("无权查看该图片")
	}
	return &v1.PictureAdminGetRes{
		Picture: s.entityToPicture(ctx, picture),
	}, nil
}

// GetVO 获取图片详情VO
func (s *sPicture) GetVO(ctx context.Context, req *v1.PictureGetReq) (res *v1.PictureGetRes, err error) {
	// VO版本可以包含更多关联信息，如用户信息等
	resp, err := s.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	// 处理标签
	var tags []string
	if resp.Tags != "" {
		err := gjson.DecodeTo(resp.Tags, &tags)
		if err != nil {
			// 如果JSON解析失败，尝试按逗号分割
			tags = strings.Split(resp.Tags, ",")
			for i, tag := range tags {
				tags[i] = strings.TrimSpace(tag)
			}
		}
	}

	// 查询用户信息
	var userVO *v1.UserVO
	if resp.UserId > 0 {
		userResp, err := service.User().GetUserById(ctx, &v1.GetUserByIdReq{Id: resp.UserId})
		if err == nil && userResp != nil {
			userVO = &v1.UserVO{
				Id:         userResp.Id,
				UserName:   userResp.UserName,
				UserAvatar: userResp.UserAvatar,
			}
		}
	}

	// 获取当前登录用户
	user, userErr := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	var permissions []string
	if userErr == nil && user != nil {
		// 获取用户对该图片的权限
		permissions = s.getPicturePermissions(ctx, resp.Picture, user)
	}

	res = &v1.PictureGetRes{
		PictureVO: &v1.PictureVO{
			Id:             resp.Id,
			Url:            resp.Url,
			Introduction:   resp.Introduction,
			Name:           resp.Name,
			PicFormat:      resp.PicFormat,
			PicWidth:       resp.PicWidth,
			PicHeight:      resp.PicHeight,
			PicScale:       resp.PicScale,
			PicColor:       resp.PicColor,
			PicSize:        resp.PicSize,
			UserId:         resp.UserId,
			CreateTime:     resp.CreateTime,
			UpdateTime:     resp.UpdateTime,
			Category:       resp.Category,
			Tags:           tags,
			ThumbnailUrl:   resp.ThumbnailUrl,
			User:           userVO,
			PermissionList: permissions,
		},
	}
	return
}

// ListByPage 分页查询图片
func (s *sPicture) ListByPage(ctx context.Context, req *v1.PictureQueryReq) (res *v1.PictureAdminQueryRes, err error) {
	g.Log().Infof(ctx, "图片列表查询请求参数: %+v", req)

	// 构建查询条件
	var queryW string
	for index, tag := range req.Tags {
		if index != 0 {
			queryW += " AND "
		}
		queryW += fmt.Sprintf("JSON_CONTAINS(tags, '\"%s\"')", tag)
	}

	pic := dao.Picture.Columns()
	db := dao.Picture.Ctx(ctx).Where(pic.IsDelete, 0)
	g.Log().Infof(ctx, "已添加删除过滤条件: isDelete = 0")

	user, _ := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if user == nil || req.ReviewStatus == nil && user.UserRole == consts.DefaultRole && req.SpaceId == "" {
		db = db.Where(pic.ReviewStatus, 1)
	} else if req.ReviewStatus != nil && user.UserRole == consts.Admin {
		db = db.Where(pic.ReviewStatus, *req.ReviewStatus)
	}

	if req.Category != "" {
		db = db.Where(pic.Category, req.Category)
	}

	if len(req.Tags) > 0 {
		db = db.Where(queryW)
	}

	if req.SearchText != "" {
		db = db.WhereLike(pic.Introduction, "%"+req.SearchText+"%")
	}
	if req.SpaceId == "[object Object]" {
		g.Log().Infof(ctx, "处理特殊spaceId: [object Object]")
		user, _ := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
		var spaceId int
		space := dao.Space.Columns()
		id, err := dao.Space.Ctx(ctx).Fields(space.Id).Where(space.UserId, user.Id).Value()
		if err != nil || id.IsEmpty() {
			return nil, gerror.New("空间不存在")
		}
		spaceId = id.Int()
		g.Log().Infof(ctx, "用户默认空间ID: %d", spaceId)
		db = db.Where(pic.SpaceId, spaceId)
	} else if req.SpaceId != "" {
		spaceIdInt := gconv.Int64(req.SpaceId)
		g.Log().Infof(ctx, "过滤指定空间ID: %s -> %d", req.SpaceId, spaceIdInt)
		db = db.Where(pic.SpaceId, spaceIdInt)
	} else {
		g.Log().Infof(ctx, "查询所有空间的图片（包括公共图片）")
	}

	// 查询总数
	total, err := db.Count()
	if err != nil {
		return nil, gerror.New("查询失败")
	}

	// 排序
	var orderBy string
	if req.SortField == "" {
		orderBy = "createTime DESC"
	} else {
		var sort string
		switch req.SortOrder {
		case "ascend":
			sort = "ASC"
		case "descend":
			sort = "DESC"
		default:
			sort = "DESC"
		}
		orderBy = req.SortField + " " + sort
	}

	// 计算总页数
	pages := (total + req.PageSize - 1) / req.PageSize
	var rKey string
	var records []v1.Picture

	if req.SpaceId == "" {
		rKey = fmt.Sprintf("picture:page:%d:%d:%s:%s:%s:%s:%d:%s:%v", req.Current, req.PageSize, req.Category, req.SearchText, req.SortField, req.SortOrder, req.ReviewStatus, req.Tags, req.SpaceId)
		rKey, err = gmd5.Encrypt(rKey)
		if err != nil {
			return nil, gerror.New("md5 encrypt failed")
		}
		// 从缓存中获取数据
		if value, err := g.Redis().Get(ctx, rKey); err == nil && !value.IsEmpty() {
			err = gjson.DecodeTo(value, &records)
			if err != nil {
				return nil, gerror.New("gjson decode to records failed")
			}
			return &v1.PictureAdminQueryRes{
				Records: records,
				PageInfo: &v1.PageInfo{
					Current: req.Current,
					Size:    req.PageSize,
					Total:   total,
					Pages:   pages,
				},
			}, nil
		}
	}

	// 分页查询
	var pictures []entity.Picture
	err = db.Page(req.Current, req.PageSize).Order(orderBy).Scan(&pictures)
	if err != nil {
		return nil, gerror.New("查询失败")
	}
	g.Log().Infof(ctx, "数据库查询结果数量: %d", len(pictures))
	for i, picture := range pictures {
		g.Log().Infof(ctx, "图片%d: ID=%d, Name=%s, SpaceId=%d, IsDelete=%d", i+1, picture.Id, picture.Name, picture.SpaceId, picture.IsDelete)
		records = append(records, *s.entityToPicture(ctx, &picture))
	}
	if req.SpaceId == "" {
		data, err := gjson.Encode(records)
		if err == nil {
			err = g.Redis().SetEX(ctx, rKey, data, int64(300+grand.Intn(121)))
			if err != nil {
				return nil, gerror.New("set redis failed")
			}
		}
	}

	return &v1.PictureAdminQueryRes{
		Records: records,
		PageInfo: &v1.PageInfo{
			Current: req.Current,
			Size:    req.PageSize,
			Total:   total,
			Pages:   pages,
		},
	}, nil
}

// ListVOByPage 分页查询图片VO
func (s *sPicture) ListVOByPage(ctx context.Context, req *v1.PictureQueryReq) (res *v1.PictureQueryRes, err error) {
	// VO版本可以包含更多关联信息
	resp, err := s.ListByPage(ctx, req)
	if err != nil {
		return nil, err
	}
	res = &v1.PictureQueryRes{
		Records:  make([]v1.PictureVO, len(resp.Records)),
		PageInfo: resp.PageInfo,
	}

	// 收集所有用户ID，批量查询用户信息
	userIds := make([]int64, 0)
	userIdSet := make(map[int64]bool)
	for _, record := range resp.Records {
		if record.UserId > 0 && !userIdSet[record.UserId] {
			userIds = append(userIds, record.UserId)
			userIdSet[record.UserId] = true
		}
	}

	// 批量查询用户信息
	userMap := make(map[int64]*v1.UserVO)
	for _, userId := range userIds {
		userResp, err := service.User().GetUserById(ctx, &v1.GetUserByIdReq{Id: userId})
		if err == nil && userResp != nil {
			userMap[userId] = &v1.UserVO{
				Id:         userResp.Id,
				UserName:   userResp.UserName,
				UserAvatar: userResp.UserAvatar,
			}
		}
	}

	for i, record := range resp.Records {
		res.Records[i] = v1.PictureVO{
			Id:           record.Id,
			Url:          record.Url,
			Introduction: record.Introduction,
			Name:         record.Name,
			PicFormat:    record.PicFormat,
			PicWidth:     record.PicWidth,
			PicHeight:    record.PicHeight,
			PicScale:     record.PicScale,
			PicColor:     record.PicColor,
			PicSize:      record.PicSize,
			UserId:       record.UserId,
			CreateTime:   record.CreateTime,
			UpdateTime:   record.UpdateTime,
			Category:     record.Category,
			ThumbnailUrl: record.ThumbnailUrl,
			SpaceId:      record.SpaceId,
			EditTime:     record.EditTime,
			User:         userMap[record.UserId],
		}
	}
	for index, picture := range resp.Records {
		var tags []string
		if picture.Tags != "" {
			err := gjson.DecodeTo(picture.Tags, &tags)
			if err != nil {
				// 如果JSON解析失败，尝试按逗号分割
				tags = strings.Split(picture.Tags, ",")
				for i, tag := range tags {
					tags[i] = strings.TrimSpace(tag)
				}
			}
			res.Records[index].Tags = tags
		}
	}
	return
}

// SearchByPicture 以图搜图
func (s *sPicture) SearchByPicture(ctx context.Context, req *v1.SearchPictureByPictureReq) (res []v1.SearchPictureByPictureRes, err error) {
	// 获取当前登录用户
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	// 获取要搜索的图片信息
	var picture *entity.Picture
	err = dao.Picture.Ctx(ctx).Fields(dao.Picture.Columns().Url, dao.Picture.Columns().UserId).Where(dao.Picture.Columns().Id, req.PictureId).
		Where(dao.Picture.Columns().IsDelete, 0).Scan(&picture)
	if err != nil || picture == nil {
		return nil, gerror.New("图片不存在")
	}

	// 暂时简化处理，只检查是否是图片上传者
	if picture.UserId != user.Id {
		// TODO: 检查用户是否是图片所属空间的成员
		// 这里可以添加更复杂的权限检查逻辑
		return nil, gerror.New("无权访问该图片")
	}

	// 调用具体的以图搜图逻辑
	// 这里返回一个接口，具体的实现由用户完成
	results, err := s.performImageSearch(ctx, picture)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// performImageSearch 执行图片搜索的具体逻辑
// 调用百度识图API进行以图搜图
func (s *sPicture) performImageSearch(ctx context.Context, picture *entity.Picture) ([]v1.SearchPictureByPictureRes, error) {
	// 构建请求URL，添加时间戳参数
	requestUrl := "https://graph.baidu.com/upload?uptime=" + gconv.String(gtime.Now().UnixMilli())

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 60 * time.Second, // 60秒超时
	}

	// 创建multipart表单数据
	var buf strings.Builder
	writer := multipart.NewWriter(&buf)

	// 添加表单字段 - 按照成功的curl命令
	formFields := map[string]string{
		"image":        picture.Url, // 图片URL
		"tn":           "pc",
		"from":         "pc",
		"image_source": "PC_UPLOAD_URL",
		"sdkParams":    `{"data":"c59dff933e7b5af7ec996e1db24ea811a1124f43b982686351512eb362403e24439dcdb0a7e89c6230b2b17e8b2605b277b21291f05e0b9ad4fb53772d99656e3e38bf3f0f73cfb07f77c9379fb0be11","key_id":"23","sign":"5174a3fe"}`, // 从curl命令中复制的sdkParams
	}

	for key, value := range formFields {
		err := writer.WriteField(key, value)
		if err != nil {
			g.Log().Error(ctx, "写入表单字段失败:", err)
			return nil, gerror.New("构建请求数据失败")
		}
	}

	// 关闭writer
	err := writer.Close()
	if err != nil {
		g.Log().Error(ctx, "关闭multipart writer失败:", err)
		return nil, gerror.New("构建请求数据失败")
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", requestUrl, strings.NewReader(buf.String()))
	if err != nil {
		g.Log().Error(ctx, "创建HTTP请求失败:", err)
		return nil, gerror.New("创建请求失败")
	}

	// 设置请求头 - 完全按照成功的curl命令
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-CN;q=0.8,en;q=0.7")
	req.Header.Set("Acs-Token", "1758428294445_1758482047547_fgUhsYwJFQyisJuBsqSBD9DDQkm+eispxH93g21ochZ0hnvCP4+pIEevGWJxiJz0tbneA5QEdfCDpaUYcqDeNg1UsPqRAKe5gJtkR4SO2oP/PIKIwaxz1TNib/V7JnYcVCZ5TODwCZLFmzyNIHOI/AEr0HD/NqxczuSabKtT28G/4DFGGAIaLb+AOCagb+zEoE9wI1ti2Nk3QYqDGqI3wJMVpp3SJkhJepAvYVBU6FELpiNyjNL99T44lOXK0SY82oYWO7bqVl7u85LPswfyG21nqLz46h/VYl+TgDDj4Ah8zVOlQA6FS16yOSjgTYBZAM0uogF407+fKIcJfrDl0QkmblmknQ9gqmsPKccSZXP2FoZQAoS2rLd63Ymr5CcXDePTF26N8EVBjxLR1uwEYrSPfFWCBGfxEMpvIipeosngVKd/aFGi9NYgo/0VFh4eO8gLZ1dOuCpH7PFKiXDC/D7HkrniQMWslQOvBTic7nQ=")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Origin", "https://graph.baidu.com")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://graph.baidu.com/pcpage/index?tpl_from=pc")
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="140", "Not=A?Brand";v="24", "Google Chrome";v="140"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	// 发送请求
	response, err := client.Do(req)
	if err != nil {
		return nil, gerror.New("图片搜索服务暂时不可用")
	}
	defer response.Body.Close()

	// 检查响应状态
	g.Log().Infof(ctx, "百度识图API响应状态码: %d", response.StatusCode)
	if response.StatusCode != 200 {
		// 读取错误响应内容
		errorText, _ := io.ReadAll(response.Body)
		g.Log().Errorf(ctx, "百度识图API返回错误状态码: %d, 响应内容: %s", response.StatusCode, string(errorText))
		return nil, gerror.New("图片搜索服务暂时不可用")
	}

	// 读取响应体
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, gerror.New("读取响应失败")
	}

	// 解析响应JSON
	var result map[string]interface{}
	err = gjson.DecodeTo(string(body), &result)
	if err != nil {
		return nil, gerror.New("解析搜索结果失败")
	}

	// 检查响应状态
	status, ok := result["status"].(float64)
	if !ok || status != 0 {
		return nil, gerror.New("解析搜索结果失败")
	}

	// 获取data字段
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return nil, gerror.New("百度识图API响应格式错误")
	}

	// 获取搜索页面URL
	searchUrl, ok := data["url"].(string)
	if !ok || searchUrl == "" {
		return nil, gerror.New("未返回有效的搜索结果地址")
	}

	_, ok = data["sign"].(string)
	if !ok {
		return nil, gerror.New("未返回有效的搜索结果签名")
	}
	// URL解码
	// decodedUrl, err := url.QueryUnescape(searchUrl)
	// if err != nil {
	// 	g.Log().Errorf(ctx, "URL解码失败: %v", err)
	// 	decodedUrl = searchUrl // 如果解码失败，使用原始URL
	// }

	results, err := s.parseBaiduImageSearchResult(searchUrl)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// parseBaiduImageSearchResult 解析百度识图API的响应结果
func (s *sPicture) parseBaiduImageSearchResult(url string) ([]v1.SearchPictureByPictureRes, error) {
	var results []v1.SearchPictureByPictureRes

	// 第一步：访问百度识图搜索页面，获取HTML内容
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	g.Log().Infof(context.Background(), "访问百度识图搜索页面: %s", url)

	// 创建请求访问搜索页面
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		g.Log().Error(context.Background(), "创建搜索页面请求失败:", err)
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-CN;q=0.8,en;q=0.7")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="140", "Not=A?Brand";v="24", "Google Chrome";v="140"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)

	// 发送请求获取HTML页面
	response, err := client.Do(req)
	if err != nil {
		g.Log().Error(context.Background(), "访问搜索页面失败:", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		g.Log().Errorf(context.Background(), "搜索页面返回错误状态码: %d", response.StatusCode)
		return nil, gerror.New("访问搜索页面失败")
	}

	// 读取HTML内容
	htmlBody, err := io.ReadAll(response.Body)
	if err != nil {
		g.Log().Error(context.Background(), "读取搜索页面内容失败:", err)
		return nil, err
	}

	htmlContent := string(htmlBody)
	g.Log().Infof(context.Background(), "获取到HTML页面内容，长度: %d", len(htmlContent))

	// 第二步：从HTML中提取firstUrl
	firstUrl := s.extractFirstUrlFromHTML(htmlContent)
	if firstUrl == "" {
		g.Log().Error(context.Background(), "未找到firstUrl")
		return nil, gerror.New("未找到图片搜索API地址")
	}

	g.Log().Infof(context.Background(), "提取到firstUrl: %s", firstUrl)

	// 第三步：访问firstUrl获取图片搜索结果
	results, err = s.fetchImageSearchResults(client, firstUrl)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// extractFirstUrlFromHTML 从HTML内容中提取firstUrl
func (s *sPicture) extractFirstUrlFromHTML(htmlContent string) string {
	// 使用正则表达式查找firstUrl
	// 匹配模式: "firstUrl":"http://..."
	pattern := `"firstUrl"\s*:\s*"(.*?)"`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(htmlContent)

	if len(matches) > 1 {
		firstUrl := matches[1]
		// 处理转义字符
		firstUrl = strings.ReplaceAll(firstUrl, "\\/", "/")
		firstUrl = strings.ReplaceAll(firstUrl, "\\\"", "\"")
		firstUrl = strings.ReplaceAll(firstUrl, "\\\\", "\\")
		return firstUrl
	}

	return ""
}

// fetchImageSearchResults 访问firstUrl获取图片搜索结果
func (s *sPicture) fetchImageSearchResults(client *http.Client, firstUrl string) ([]v1.SearchPictureByPictureRes, error) {

	g.Log().Infof(context.Background(), "访问图片搜索API: %s", firstUrl)

	// 创建请求访问图片搜索API
	req, err := http.NewRequest("GET", firstUrl, nil)
	if err != nil {
		g.Log().Error(context.Background(), "创建图片搜索请求失败:", err)
		return nil, err
	}

	// 设置请求头 - 按照你提供的curl命令
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-CN;q=0.8,en;q=0.7")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Cookie", "BIDUPSID=28FC22613AC3D6C07E4BC6C306EB2826; PSTM=1758480992; BAIDUID=28FC22613AC3D6C07AC0D8F17614F2EC:FG=1; H_PS_PSSID=63140_64559_64661_64746_64700_64818_64839_64910_64927_64989_65007_65080_65122_65142_65139_65137_65188_65204_65223_65247_65254_65143_65323_65367_65374_65384; BAIDUID_BFESS=28FC22613AC3D6C07AC0D8F17614F2EC:FG=1; BA_HECTOR=al8ga080802l0481242h20258kal851kd0ij424; ZFY=FsywclYYUFPwu7waofrA:BvAVpXwGoJnpzK5tnYPBk:BQ:C; PSINO=7; delPer=0; H_WISE_SIDS=63140_64559_64661_64746_64700_64818_64839_64910_64927_64989_65007_65080_65122_65142_65139_65137_65188_65204_65223_65247_65254_65143_65323_65367_65374_65384; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; antispam_key_id=23; antispam_data=c59dff933e7b5af7ec996e1db24ea811a1124f43b982686351512eb362403e24439dcdb0a7e89c6230b2b17e8b2605b277b21291f05e0b9ad4fb53772d99656e3e38bf3f0f73cfb07f77c9379fb0be11; ab_sr=1.0.1_YjJjNGNhNjQ5OWEyODZiMzNhYmZiNjZjMjBhMGIxYmM4NmM1YzQyYjJlZjEzNGJiNzE4ZTk4MTk2OGYxNjQ4M2RlOGY2YjMzODM5MDcwYTJjYjk2MWE5NDhiM2I2MzMxYTUzZGNlYzFjODYxZjBiNjk1MDIyYjIxNjgyN2JhMzA2OTg4OWQyYTdlODllMzVkYmYzNDM1MjMzMDg0NjQ3Zg==")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://graph.baidu.com/s?card_key=&entrance=GENERAL&extUiData%5BisLogoShow%5D=1&f=all&isLogoShow=1&session_id=17494005287103596248&sign=1261e0f3ee1497bfe556101758482048&tpl_from=pc")
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="140", "Not=A?Brand";v="24", "Google Chrome";v="140"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")

	// 发送请求
	response, err := client.Do(req)
	if err != nil {
		g.Log().Error(context.Background(), "发送图片搜索请求失败:", err)
		return nil, err
	}
	defer response.Body.Close()

	// 检查响应状态
	g.Log().Infof(context.Background(), "图片搜索API响应状态码: %d", response.StatusCode)
	if response.StatusCode != 200 {
		// 读取错误响应内容
		errorText, _ := io.ReadAll(response.Body)
		g.Log().Errorf(context.Background(), "图片搜索API返回错误状态码: %d, 响应内容: %s", response.StatusCode, string(errorText))
		return nil, gerror.New("图片搜索API返回错误状态")
	}

	// 读取响应体
	body, err := io.ReadAll(response.Body)
	if err != nil {
		g.Log().Error(context.Background(), "读取图片搜索响应失败:", err)
		return nil, err
	}

	var result map[string]interface{}
	err = gjson.DecodeTo(body, &result)
	if err != nil {
		g.Log().Error(context.Background(), "解析百度搜索响应失败:", err)
		return nil, err
	}

	// 获取搜索结果
	data := gconv.Map(result["data"])

	length := gconv.Int(data["length"])

	results := make([]v1.SearchPictureByPictureRes, 0, length)

	images := gconv.SliceMap(data["list"])

	for _, image := range images {
		results = append(results, v1.SearchPictureByPictureRes{
			FromUrl:  gconv.String(image["fromUrl"]),
			ThumbUrl: gconv.String(image["thumbUrl"]),
		})
	}

	return results, nil
}

// SearchByColor 按颜色搜索图片（基于欧氏距离的相似度搜索）
func (s *sPicture) SearchByColor(ctx context.Context, req *v1.SearchPictureByColorReq) (res []v1.SearchPictureByColorRes, err error) {
	// 获取当前登录用户
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	// 验证颜色格式
	if req.PicColor == "" {
		return nil, gerror.New("颜色值不能为空")
	}

	// 解析请求的颜色值
	targetRGB, err := s.parseHexColor(req.PicColor)
	if err != nil {
		return nil, gerror.New("颜色格式无效，请使用十六进制格式如 #FF5733")
	}

	// 构建查询条件
	query := dao.Picture.Ctx(ctx).Where(dao.Picture.Columns().IsDelete, 0)

	// 如果指定了空间ID，只搜索该空间的图片
	if req.SpaceId > 0 {
		query = query.Where(dao.Picture.Columns().SpaceId, req.SpaceId)
		// 检查用户是否有权限访问该空间
		// TODO: 这里可以添加更复杂的权限检查逻辑
	} else if user.UserRole != consts.Admin {
		// 如果没有指定空间ID，只搜索用户自己的图片
		query = query.Where(dao.Picture.Columns().UserId, user.Id)
	}

	// 只查询有主色调的图片
	query = query.Where(dao.Picture.Columns().PicColor+" != ?", "")

	// 执行查询
	var pictures []entity.Picture
	err = query.Scan(&pictures)
	if err != nil {
		g.Log().Errorf(ctx, "查询图片失败: %v", err)
		return nil, gerror.New("查询图片失败")
	}

	// 计算颜色相似度并排序
	similarPictures := s.calculateColorSimilarity(ctx, pictures, targetRGB)

	// 返回前12个最相似的图片
	limit := 12
	if len(similarPictures) > limit {
		similarPictures = similarPictures[:limit]
	}

	// 转换为响应对象
	records := make([]v1.SearchPictureByColorRes, 0, len(similarPictures))
	for _, item := range similarPictures {
		pictureVO := s.entityToVO(ctx, item.Picture)
		records = append(records, v1.SearchPictureByColorRes{
			PictureVO: pictureVO,
		})
	}

	g.Log().Infof(ctx, "按颜色相似度搜索图片成功，找到 %d 张相似图片", len(records))

	return records, nil
}

// getPicturePermissions 获取用户对图片的权限
func (s *sPicture) getPicturePermissions(ctx context.Context, picture *v1.Picture, user *v1.GetLoginUserRes) []string {
	permissions := make([]string, 0)

	// 1. 管理员拥有所有权限
	if user.UserRole == consts.Admin {
		return []string{
			"picture:view",
			"picture:upload",
			"picture:edit",
			"picture:delete",
			"spaceUser:manage",
		}
	}

	// 2. 图片所有者拥有所有权限（除了空间管理）
	if picture.UserId == user.Id {
		return []string{
			"picture:view",
			"picture:upload",
			"picture:edit",
			"picture:delete",
		}
	}

	// 3. 如果图片属于空间，检查用户在空间中的权限
	if picture.SpaceId > 0 {
		// 检查用户是否是空间所有者
		var space *entity.Space
		err := dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, picture.SpaceId).
			Where(dao.Space.Columns().IsDelete, 0).Scan(&space)
		if err == nil && space != nil && space.UserId == user.Id {
			return []string{
				"picture:view",
				"picture:upload",
				"picture:edit",
				"picture:delete",
				"spaceUser:manage",
			}
		}

		// 检查用户在空间中的角色
		var spaceUser *entity.SpaceUser
		err = dao.SpaceUser.Ctx(ctx).Where(dao.SpaceUser.Columns().SpaceId, picture.SpaceId).
			Where(dao.SpaceUser.Columns().UserId, user.Id).Scan(&spaceUser)
		if err == nil && spaceUser != nil {
			switch spaceUser.SpaceRole {
			case "admin":
				return []string{
					"picture:view",
					"picture:upload",
					"picture:edit",
					"picture:delete",
					"spaceUser:manage",
				}
			case "editor":
				return []string{
					"picture:view",
					"picture:upload",
					"picture:edit",
				}
			case "viewer":
				return []string{
					"picture:view",
				}
			}
		}
	}

	// 4. 默认权限：只能查看已审核的公共图片
	if picture.SpaceId == 0 && picture.ReviewStatus == 1 {
		permissions = append(permissions, "picture:view")
	}

	return permissions
}
