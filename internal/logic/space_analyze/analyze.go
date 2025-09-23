package space_analyze

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/consts"
	"cloud/internal/dao"
	"cloud/internal/model/entity"
	"cloud/internal/service"
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// CategoryAnalyze 空间分类分析
func (s *sSpaceAnalyze) CategoryAnalyze(ctx context.Context, req *v1.SpaceCategoryAnalyzeReq) (res *v1.SpaceCategoryAnalyzeRes, err error) {
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	model, err := buildPictureScopeModel(ctx, user.Id, user.UserRole == consts.Admin, req.SpaceId, req.QueryAll, req.QueryPublic)
	if err != nil {
		return nil, err
	}

	// 分类统计：按 category 分组统计数量
	type row struct {
		Category  string `json:"category"`
		Count     int64  `json:"count"`
		TotalSize int64  `json:"totalSize"`
	}
	var rows []row
	err = model.Where(dao.Picture.Columns().IsDelete, 0).
		Fields(dao.Picture.Columns().Category + " as category, COUNT(1) as count, IFNULL(SUM(" + dao.Picture.Columns().PicSize + "),0) as totalSize").
		Group(dao.Picture.Columns().Category).
		OrderDesc("count").
		Scan(&rows)
	if err != nil {
		return nil, err
	}

	records := make([]v1.SpaceCategoryAnalyzeResponse, 0, len(rows))
	for _, r := range rows {
		// 兼容空分类
		cat := r.Category
		if cat == "" {
			cat = "未分类"
		}
		records = append(records, v1.SpaceCategoryAnalyzeResponse{
			Category:  cat,
			Count:     r.Count,
			TotalSize: r.TotalSize,
		})
	}
	resArr := v1.SpaceCategoryAnalyzeRes(records)
	return &resArr, nil
}

// CheckSpaceAuth 检查空间权限
func (s *sSpaceAnalyze) CheckSpaceAuth(ctx context.Context, req *v1.SpaceCategoryAnalyzeReq) (err error) {
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return gerror.New("用户未登录")
	}

	if req.QueryAll || req.QueryPublic {
		if user.UserRole != consts.Admin {
			return gerror.New("无权限")
		}
	} else {
		var space *entity.Space
		err = dao.Space.Ctx(ctx).Where(dao.Space.Columns().Id, req.SpaceId).Where(dao.Space.Columns().IsDelete, 0).Scan(&space)
		if err != nil || space == nil {
			return gerror.New("空间不存在")
		}
		if space.UserId != user.Id && space.SpaceType == consts.SpaceTypePrivate && user.UserRole != consts.Admin {
			return gerror.New("无权限访问此空间")
		}
	}
	return nil
}

// TagAnalyze 空间标签分析
func (s *sSpaceAnalyze) TagAnalyze(ctx context.Context, req *v1.SpaceTagAnalyzeReq) (res *v1.SpaceTagAnalyzeRes, err error) {
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	g.Log().Infof(ctx, "用户 %d 请求空间标签分析", user.Id)

	model, err := buildPictureScopeModel(ctx, user.Id, user.UserRole == consts.Admin, req.SpaceId, req.QueryAll, req.QueryPublic)
	if err != nil {
		return nil, err
	}

	// 先取出满足条件的 tags 列，逐条解析 JSON 数组并计数
	var tagRows []struct{ Tags string }
	if err = model.Where(dao.Picture.Columns().IsDelete, 0).
		Fields(dao.Picture.Columns().Tags).
		Scan(&tagRows); err != nil {
		return nil, err
	}

	counter := make(map[string]int64)
	for _, r := range tagRows {
		if r.Tags == "" {
			continue
		}
		// 优先解析 JSON 数组
		j, jErr := gjson.LoadContent([]byte(r.Tags))
		if jErr == nil {
			if arr := j.Array(); len(arr) > 0 {
				for _, v := range arr {
					t := gconv.String(v)
					if t == "" {
						continue
					}
					counter[t]++
				}
				continue
			}
		}
		// 兜底：按逗号/空白分割
		for _, t := range gstr.SplitAndTrim(r.Tags, ",") {
			if t == "" {
				continue
			}
			counter[t]++
		}
	}

	records := make([]v1.SpaceTagAnalyzeResponse, 0, len(counter))
	for tag, cnt := range counter {
		records = append(records, v1.SpaceTagAnalyzeResponse{Tag: tag, Count: cnt})
	}
	resArr := v1.SpaceTagAnalyzeRes(records)
	return &resArr, nil
}

// SizeAnalyze 空间大小分析
func (s *sSpaceAnalyze) SizeAnalyze(ctx context.Context, req *v1.SpaceSizeAnalyzeReq) (res *v1.SpaceSizeAnalyzeRes, err error) {
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}
	g.Log().Infof(ctx, "用户 %d 请求空间大小分析", user.Id)

	model, err := buildPictureScopeModel(ctx, user.Id, user.UserRole == consts.Admin, req.SpaceId, req.QueryAll, req.QueryPublic)
	if err != nil {
		return nil, err
	}

	// 定义区间（字节）
	const (
		mb1  int64 = 1 << 20
		mb5  int64 = 5 << 20
		mb10 int64 = 10 << 20
	)

	// 分别统计不同区间数量
	countRange := func(cond string, args ...interface{}) (int64, error) {
		m := model.Clone().Where(dao.Picture.Columns().IsDelete, 0)
		if cond != "" {
			m = m.Where(cond, args...)
		}
		nInt, err := m.CountColumn(dao.Picture.Columns().Id)
		if err != nil {
			return 0, err
		}
		return gconv.Int64(nInt), nil
	}

	lt1, err := countRange(dao.Picture.Columns().PicSize+" < ?", mb1)
	if err != nil {
		return nil, err
	}
	btw1to5, err := countRange(dao.Picture.Columns().PicSize+" >= ? AND "+dao.Picture.Columns().PicSize+" < ?", mb1, mb5)
	if err != nil {
		return nil, err
	}
	btw5to10, err := countRange(dao.Picture.Columns().PicSize+" >= ? AND "+dao.Picture.Columns().PicSize+" < ?", mb5, mb10)
	if err != nil {
		return nil, err
	}
	gt10, err := countRange(dao.Picture.Columns().PicSize+" >= ?", mb10)
	if err != nil {
		return nil, err
	}

	records := []v1.SpaceSizeAnalyzeResponse{
		{SizeRange: "< 1MB", Count: lt1},
		{SizeRange: "1-5MB", Count: btw1to5},
		{SizeRange: "5-10MB", Count: btw5to10},
		{SizeRange: "> 10MB", Count: gt10},
	}
	resArr := v1.SpaceSizeAnalyzeRes(records)
	return &resArr, nil
}

// UsageAnalyze 空间使用情况分析
func (s *sSpaceAnalyze) UsageAnalyze(ctx context.Context, req *v1.SpaceUsageAnalyzeReq) (res *v1.SpaceUsageAnalyzeRes, err error) {
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	g.Log().Infof(ctx, "用户 %d 请求空间使用情况分析", user.Id)

	// 三种模式：指定空间 / 全部空间 / 公共空间
	if req.QueryAll {
		if user.UserRole != consts.Admin {
			return nil, gerror.New("无权限")
		}
		// 聚合所有空间用量
		var rows []entity.Space
		if err := dao.Space.Ctx(ctx).
			Where(dao.Space.Columns().IsDelete, 0).
			Scan(&rows); err != nil {
			return nil, err
		}
		var maxCount, totalCount, maxSize, totalSize int64
		for _, s := range rows {
			maxCount += s.MaxCount
			totalCount += s.TotalCount
			maxSize += s.MaxSize
			totalSize += s.TotalSize
		}
		response := &v1.SpaceUsageAnalyzeResponse{
			MaxCount:        maxCount,
			UsedCount:       totalCount,
			CountUsageRatio: ratio(totalCount, maxCount) * 100,
			MaxSize:         maxSize,
			UsedSize:        totalSize,
			SizeUsageRatio:  ratio(totalSize, maxSize) * 100,
		}
		return &v1.SpaceUsageAnalyzeRes{SpaceUsageAnalyzeResponse: response}, nil
	}

	if req.QueryPublic {
		if user.UserRole != consts.Admin {
			return nil, gerror.New("无权限")
		}
		// 公共空间：直接聚合图片
		model, err := buildPictureScopeModel(ctx, user.Id, true, 0, false, true)
		if err != nil {
			return nil, err
		}
		nInt, err := model.Where(dao.Picture.Columns().IsDelete, 0).CountColumn(dao.Picture.Columns().Id)
		if err != nil {
			return nil, err
		}
		totalCount := gconv.Int64(nInt)
		v, err := model.Where(dao.Picture.Columns().IsDelete, 0).Fields("IFNULL(SUM(" + dao.Picture.Columns().PicSize + "),0) as s").Value("s")
		if err != nil {
			return nil, err
		}
		totalSize := gconv.Int64(v)
		response := &v1.SpaceUsageAnalyzeResponse{
			MaxCount:        0,
			UsedCount:       totalCount,
			CountUsageRatio: 0,
			MaxSize:         0,
			UsedSize:        totalSize,
			SizeUsageRatio:  0,
		}
		return &v1.SpaceUsageAnalyzeRes{SpaceUsageAnalyzeResponse: response}, nil
	}

	// 指定空间
	var space *entity.Space
	if err := dao.Space.Ctx(ctx).
		Where(dao.Space.Columns().Id, req.SpaceId).
		Where(dao.Space.Columns().IsDelete, 0).
		Scan(&space); err != nil || space == nil {
		return nil, gerror.New("空间不存在")
	}
	if user.UserRole != consts.Admin && space.UserId != user.Id {
		return nil, gerror.New("无权限访问此空间")
	}

	response := &v1.SpaceUsageAnalyzeResponse{
		MaxCount:        space.MaxCount,
		UsedCount:       space.TotalCount,
		CountUsageRatio: ratio(space.TotalCount, space.MaxCount) * 100,
		MaxSize:         space.MaxSize,
		UsedSize:        space.TotalSize,
		SizeUsageRatio:  ratio(space.TotalSize, space.MaxSize) * 100,
	}
	return &v1.SpaceUsageAnalyzeRes{SpaceUsageAnalyzeResponse: response}, nil
}

// UserAnalyze 空间用户分析
func (s *sSpaceAnalyze) UserAnalyze(ctx context.Context, req *v1.SpaceUserAnalyzeReq) (res *v1.SpaceUserAnalyzeRes, err error) {
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	g.Log().Infof(ctx, "用户 %d 请求空间用户分析", user.Id)

	// 依据时间维度聚合上传数量（按 picture.createTime）
	model, err := buildPictureScopeModel(ctx, user.Id, user.UserRole == consts.Admin, req.SpaceId, req.QueryAll, req.QueryPublic)
	if err != nil {
		return nil, err
	}

	timeFmt := "%Y-%m-%d"
	switch req.TimeDimension {
	case "week":
		// 以年-周作为分组
		timeFmt = "%x-第%v周"
	case "month":
		timeFmt = "%Y-%m"
	}

	type row struct {
		Period string `json:"period"`
		Count  int64  `json:"count"`
	}
	var rows []row
	if err := model.Where(dao.Picture.Columns().IsDelete, 0).
		Fields("DATE_FORMAT(" + dao.Picture.Columns().CreateTime + ", '" + timeFmt + "') as period, COUNT(1) as count").
		Group("period").
		Order("period asc").
		Scan(&rows); err != nil {
		return nil, err
	}

	records := make([]v1.SpaceUserAnalyzeResponse, 0, len(rows))
	for _, r := range rows {
		records = append(records, v1.SpaceUserAnalyzeResponse{Period: r.Period, Count: r.Count})
	}
	resArr := v1.SpaceUserAnalyzeRes(records)
	return &resArr, nil
}

// RankAnalyze 空间排行分析
func (s *sSpaceAnalyze) RankAnalyze(ctx context.Context, req *v1.SpaceRankAnalyzeReq) (res *v1.SpaceRankAnalyzeRes, err error) {
	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, gerror.New("用户未登录")
	}

	g.Log().Infof(ctx, "用户 %d 请求空间排行分析，TopN: %d", user.Id, req.TopN)

	// 仅管理员可查看全站排行
	if user.UserRole != consts.Admin {
		return nil, gerror.New("无权限")
	}

	// 取 topN，按图片数量降序，其次按总大小降序
	var rows []entity.Space
	if err := dao.Space.Ctx(ctx).
		Where(dao.Space.Columns().IsDelete, 0).
		OrderDesc(dao.Space.Columns().TotalCount).
		OrderDesc(dao.Space.Columns().TotalSize).
		Limit(req.TopN).
		Scan(&rows); err != nil {
		return nil, err
	}

	records := make([]v1.Space, 0, len(rows))
	for _, s := range rows {
		records = append(records, v1.Space{
			Id:         s.Id,
			SpaceName:  s.SpaceName,
			SpaceLevel: s.SpaceLevel,
			MaxSize:    s.MaxSize,
			MaxCount:   s.MaxCount,
			TotalSize:  s.TotalSize,
			TotalCount: s.TotalCount,
			UserId:     s.UserId,
			CreateTime: "",
			EditTime:   "",
			UpdateTime: "",
			IsDelete:   s.IsDelete,
			SpaceType:  s.SpaceType,
		})
	}

	resArr := v1.SpaceRankAnalyzeRes(records)
	return &resArr, nil
}

// buildPictureScopeModel 构造图片查询模型，按 QueryAll / QueryPublic / SpaceId 权限过滤
func buildPictureScopeModel(ctx context.Context, loginUserId int64, isAdmin bool, spaceId int64, queryAll, queryPublic bool) (*gdb.Model, error) {
	m := dao.Picture.Ctx(ctx)
	// 公共空间
	if queryPublic {
		if !isAdmin {
			return nil, gerror.New("无权限")
		}
		return m.Where("" + dao.Picture.Columns().SpaceId + " IS NULL"), nil
	}
	// 全部空间
	if queryAll {
		if !isAdmin {
			return nil, gerror.New("无权限")
		}
		return m.Where("" + dao.Picture.Columns().SpaceId + " IS NOT NULL"), nil
	}
	// 指定空间
	var space *entity.Space
	if err := dao.Space.Ctx(ctx).
		Where(dao.Space.Columns().Id, spaceId).
		Where(dao.Space.Columns().IsDelete, 0).
		Scan(&space); err != nil || space == nil {
		return nil, gerror.New("空间不存在")
	}
	if !isAdmin && space.UserId != loginUserId {
		return nil, gerror.New("无权限访问此空间")
	}
	return m.Where(dao.Picture.Columns().SpaceId, spaceId), nil
}

// ratio 计算使用率（分母为0时返回0）
func ratio(used, max int64) float64 {
	if max <= 0 {
		return 0
	}
	return float64(used) / float64(max)
}
