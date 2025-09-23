package controller

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/service"
	"context"
)

var SpaceAnalyze = cSpaceAnalyze{}

type cSpaceAnalyze struct{}

// CategoryAnalyze 空间分类分析
func (c *cSpaceAnalyze) CategoryAnalyze(ctx context.Context, req *v1.SpaceCategoryAnalyzeReq) (res *v1.SpaceCategoryAnalyzeRes, err error) {
	return service.SpaceAnalyze().CategoryAnalyze(ctx, req)
}

// TagAnalyze 空间标签分析
func (c *cSpaceAnalyze) TagAnalyze(ctx context.Context, req *v1.SpaceTagAnalyzeReq) (res *v1.SpaceTagAnalyzeRes, err error) {
	return service.SpaceAnalyze().TagAnalyze(ctx, req)
}

// SizeAnalyze 空间大小分析
func (c *cSpaceAnalyze) SizeAnalyze(ctx context.Context, req *v1.SpaceSizeAnalyzeReq) (res *v1.SpaceSizeAnalyzeRes, err error) {
	return service.SpaceAnalyze().SizeAnalyze(ctx, req)
}

// UsageAnalyze 空间使用情况分析
func (c *cSpaceAnalyze) UsageAnalyze(ctx context.Context, req *v1.SpaceUsageAnalyzeReq) (res *v1.SpaceUsageAnalyzeRes, err error) {
	return service.SpaceAnalyze().UsageAnalyze(ctx, req)
}

// UserAnalyze 空间用户分析
func (c *cSpaceAnalyze) UserAnalyze(ctx context.Context, req *v1.SpaceUserAnalyzeReq) (res *v1.SpaceUserAnalyzeRes, err error) {
	return service.SpaceAnalyze().UserAnalyze(ctx, req)
}

// RankAnalyze 空间排行分析
func (c *cSpaceAnalyze) RankAnalyze(ctx context.Context, req *v1.SpaceRankAnalyzeReq) (res *v1.SpaceRankAnalyzeRes, err error) {
	return service.SpaceAnalyze().RankAnalyze(ctx, req)
}
