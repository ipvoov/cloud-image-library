package picture

import (
	"cloud/internal/dao"
	"cloud/internal/model/entity"
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
)

// GetPictureById 根据ID获取图片 - WebSocket专用
func GetPictureById(ctx context.Context, pictureID int64) (*entity.Picture, error) {
	var picture *entity.Picture
	err := dao.Picture.Ctx(ctx).Where(dao.Picture.Columns().Id, pictureID).
		Where(dao.Picture.Columns().IsDelete, 0).Scan(&picture)

	if err != nil {
		return nil, gerror.New("查询图片失败")
	}

	if picture == nil {
		return nil, gerror.New("图片不存在")
	}

	return picture, nil
}
