package picture

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/consts"
	"cloud/internal/dao"
	"cloud/internal/model/do"
	"cloud/internal/service"
	"context"

	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

// Upload 上传图片
func (s *sPicture) Upload(ctx context.Context, req *v1.PictureUploadReq, file *ghttp.UploadFile) (res *v1.PictureUploadRes, err error) {
	// 解析图片信息
	width, height, format, scale, parseErr := s.parseImageInfo(ctx, file)
	if parseErr != nil {
		g.Log().Warningf(ctx, "解析图片信息失败，使用默认值: %v", parseErr)
		// 如果解析失败，使用默认值
		width = 0
		height = 0
		scale = 0
		format = s.getFormatFromFilename(file.Filename)
	}

	// 提取图片主色调
	picColor := "#000000" // 默认黑色
	if parseErr == nil {
		// 只有图片解析成功时才提取主色调
		extractedColor, colorErr := s.extractDominantColor(ctx, file)
		if colorErr != nil {
			g.Log().Warningf(ctx, "提取图片主色调失败，使用默认值: %v", colorErr)
		} else {
			picColor = extractedColor
		}
	}

	// 调用bucket服务上传文件
	bucketReq := &v1.BucketUploadReq{
		File: file,
	}
	bucketRes, err := service.Bucket().Upload(ctx, bucketReq)
	if err != nil {
		g.Log().Errorf(ctx, "文件上传失败: %v", err)
		return nil, err
	}

	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, err
	}

	// 使用事务：插入图片 + 更新空间统计
	var id int64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1) 插入图片
		resp, inErr := dao.Picture.Ctx(ctx).TX(tx).Data(do.Picture{
			Url:          consts.BucketURL + bucketRes.FileAddress,
			Name:         file.Filename,
			Introduction: "",
			Category:     "默认",
			Tags:         "[]",
			PicSize:      file.Size,
			PicWidth:     width,
			PicHeight:    height,
			PicScale:     scale,
			PicFormat:    format,
			UserId:       user.Id,
			SpaceId:      req.SpaceId,
			ReviewStatus: consts.DefRwStatus,
			ThumbnailUrl: consts.BucketURL + bucketRes.FileAddress,
			PicColor:     picColor,
		}).Insert()
		if inErr != nil {
			g.Log().Errorf(ctx, "保存图片信息失败: %v", inErr)
			return inErr
		}
		lastID, idErr := resp.LastInsertId()
		if idErr != nil {
			return idErr
		}
		id = lastID

		// 2) 累加空间统计
		if req.SpaceId > 0 {
			space := dao.Space.Columns()
			_, upErr := dao.Space.Ctx(ctx).TX(tx).
				Where(space.Id, req.SpaceId).
				Data(g.Map{
					space.TotalSize:  gdb.Raw(fmt.Sprintf("totalSize + %d", file.Size)),
					space.TotalCount: gdb.Raw("totalCount + 1"),
				}).
				Update()
			if upErr != nil {
				g.Log().Errorf(ctx, "更新空间统计失败 spaceId=%d: %v", req.SpaceId, upErr)
				return upErr
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 创建图片VO对象
	pictureVO := &v1.PictureVO{
		Id:           id,                                       // 数据库生成的ID
		Url:          consts.BucketURL + bucketRes.FileAddress, // 使用实际上传后的URL
		Name:         file.Filename,
		Introduction: "默认",
		Category:     "",
		Tags:         []string{},
		PicSize:      file.Size, // 使用实际文件大小
		PicWidth:     width,     // 解析得到的图片宽度
		PicHeight:    height,    // 解析得到的图片高度
		PicScale:     scale,     // 计算得到的宽高比例
		PicFormat:    format,    // 解析得到的图片格式
		UserId:       user.Id,   // 从上下文获取当前用户ID
		SpaceId:      req.SpaceId,
		CreateTime:   gtime.Now().Format(consts.Y_m_d_His),
		EditTime:     gtime.Now().Format(consts.Y_m_d_His),
		UpdateTime:   gtime.Now().Format(consts.Y_m_d_His),
		ThumbnailUrl: consts.BucketURL + bucketRes.FileAddress, // 暂时使用原图URL
		PicColor:     picColor,                                 // 提取的主色调
	}

	return &v1.PictureUploadRes{
		PictureVO: pictureVO,
	}, nil
}

// UploadByUrl 通过URL上传图片
func (s *sPicture) UploadByUrl(ctx context.Context, req *v1.PictureUploadByUrlReq) (res *v1.PictureUploadByUrlRes, err error) {
	// 解析图片信息
	width, height, format, scale, fileSize, parseErr := s.parseImageInfoFromURL(ctx, req.FileUrl)
	if parseErr != nil {
		g.Log().Warningf(ctx, "解析URL图片信息失败，使用默认值: %v", parseErr)
		// 如果解析失败，使用默认值
		width = 0
		height = 0
		scale = 0
		fileSize = 0
		format = s.getFormatFromFilename(req.FileUrl)
	}

	// 提取图片主色调
	picColor := "#000000" // 默认黑色
	if parseErr == nil {
		// 只有图片解析成功时才提取主色调
		extractedColor, colorErr := s.extractDominantColorFromURL(ctx, req.FileUrl)
		if colorErr != nil {
			g.Log().Warningf(ctx, "提取URL图片主色调失败，使用默认值: %v", colorErr)
		} else {
			picColor = extractedColor
		}
	}

	if req.FileName != "" {
		req.FileName += service.Bucket().GetFileExtFromUrl(req.FileUrl)
	}

	// 调用bucket服务上传URL文件
	bucketRes, err := service.Bucket().UploadByUrl(ctx, &v1.BucketUploadByUrlReq{
		FileUrl:  req.FileUrl,
		FileName: req.FileName,
		SpaceId:  req.SpaceId,
	})
	if err != nil {
		g.Log().Errorf(ctx, "URL文件上传失败: %v", err)
		return nil, err
	}

	user, err := service.User().GetLoginUser(ctx, &v1.GetLoginUserReq{})
	if err != nil {
		return nil, err
	}

	// 使用bucket服务返回的文件名，如果没有则从URL提取
	filename := bucketRes.FileName
	if filename == "" {
		filename = req.FileName
	}

	// 优先使用bucket服务返回的文件大小，如果为0则使用解析得到的大小
	finalFileSize := bucketRes.FileSize
	if finalFileSize == 0 {
		finalFileSize = fileSize
	}

	// 使用事务：插入图片 + 更新空间统计
	var id int64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 1) 插入图片
		resp, inErr := dao.Picture.Ctx(ctx).TX(tx).Data(do.Picture{
			Url:          consts.BucketURL + bucketRes.FileAddress,
			Name:         filename,
			Introduction: "",
			Category:     "默认",
			Tags:         "[]",
			PicSize:      finalFileSize,
			PicWidth:     width,
			PicHeight:    height,
			PicScale:     scale,
			PicFormat:    format,
			UserId:       user.Id,
			SpaceId:      req.SpaceId,
			ReviewStatus: consts.DefRwStatus,
			ThumbnailUrl: consts.BucketURL + bucketRes.FileAddress,
			PicColor:     picColor,
		}).Insert()
		if inErr != nil {
			g.Log().Errorf(ctx, "保存图片信息失败: %v", inErr)
			return inErr
		}
		lastID, idErr := resp.LastInsertId()
		if idErr != nil {
			return idErr
		}
		id = lastID

		// 2) 累加空间统计
		if req.SpaceId > 0 {
			space := dao.Space.Columns()
			_, upErr := dao.Space.Ctx(ctx).TX(tx).
				Where(space.Id, req.SpaceId).
				Data(g.Map{
					space.TotalSize:  gdb.Raw(fmt.Sprintf("totalSize + %d", finalFileSize)),
					space.TotalCount: gdb.Raw("totalCount + 1"),
				}).
				Update()
			if upErr != nil {
				g.Log().Errorf(ctx, "更新空间统计失败 spaceId=%d: %v", req.SpaceId, upErr)
				return upErr
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	pictureVO := &v1.PictureVO{
		Id:           id,                                       // 数据库生成的ID
		Url:          consts.BucketURL + bucketRes.FileAddress, // 使用实际上传后的URL
		Name:         filename,
		Introduction: "",
		Category:     "默认",
		Tags:         []string{},
		PicSize:      finalFileSize, // 使用实际文件大小
		PicWidth:     width,         // 解析得到的图片宽度
		PicHeight:    height,        // 解析得到的图片高度
		PicScale:     scale,         // 计算得到的宽高比例
		PicFormat:    format,        // 解析得到的图片格式
		UserId:       user.Id,       // 从上下文获取当前用户ID
		SpaceId:      req.SpaceId,
		CreateTime:   gtime.Now().Format(consts.Y_m_d_His),
		EditTime:     gtime.Now().Format(consts.Y_m_d_His),
		UpdateTime:   gtime.Now().Format(consts.Y_m_d_His),
		ThumbnailUrl: consts.BucketURL + bucketRes.FileAddress, // 暂时使用原图URL
		PicColor:     picColor,                                 // 提取的主色调
	}

	return &v1.PictureUploadByUrlRes{
		PictureVO: pictureVO,
	}, nil
}
