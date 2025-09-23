// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	v1 "cloud/api/user/v1"
	"context"
)

type (
	ISpaceAnalyze interface {
		// CategoryAnalyze 空间分类分析
		CategoryAnalyze(ctx context.Context, req *v1.SpaceCategoryAnalyzeReq) (res *v1.SpaceCategoryAnalyzeRes, err error)
		// TagAnalyze 空间标签分析
		TagAnalyze(ctx context.Context, req *v1.SpaceTagAnalyzeReq) (res *v1.SpaceTagAnalyzeRes, err error)
		// SizeAnalyze 空间大小分析
		SizeAnalyze(ctx context.Context, req *v1.SpaceSizeAnalyzeReq) (res *v1.SpaceSizeAnalyzeRes, err error)
		// UsageAnalyze 空间使用情况分析
		UsageAnalyze(ctx context.Context, req *v1.SpaceUsageAnalyzeReq) (res *v1.SpaceUsageAnalyzeRes, err error)
		// UserAnalyze 空间用户分析
		UserAnalyze(ctx context.Context, req *v1.SpaceUserAnalyzeReq) (res *v1.SpaceUserAnalyzeRes, err error)
		// RankAnalyze 空间排行分析
		RankAnalyze(ctx context.Context, req *v1.SpaceRankAnalyzeReq) (res *v1.SpaceRankAnalyzeRes, err error)
	}
)

var (
	localSpaceAnalyze ISpaceAnalyze
)

func SpaceAnalyze() ISpaceAnalyze {
	if localSpaceAnalyze == nil {
		panic("implement not found for interface ISpaceAnalyze, forgot register?")
	}
	return localSpaceAnalyze
}

func RegisterSpaceAnalyze(i ISpaceAnalyze) {
	localSpaceAnalyze = i
}
