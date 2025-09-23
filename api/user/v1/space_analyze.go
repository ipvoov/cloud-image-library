package v1

// SpaceCategoryAnalyzeReq 空间分类分析请求
type SpaceCategoryAnalyzeReq struct {
	SpaceId     int64 `json:"spaceId"`
	QueryAll    bool  `json:"queryAll"`
	QueryPublic bool  `json:"queryPublic"`
}

// SpaceCategoryAnalyzeRes 空间分类分析响应
type SpaceCategoryAnalyzeRes []SpaceCategoryAnalyzeResponse

// SpaceCategoryAnalyzeResponse 空间分类分析响应项
type SpaceCategoryAnalyzeResponse struct {
	Category  string `json:"category"`
	Count     int64  `json:"count"`
	TotalSize int64  `json:"totalSize"`
}

// SpaceTagAnalyzeReq 空间标签分析请求
type SpaceTagAnalyzeReq struct {
	SpaceId     int64 `json:"spaceId"`
	QueryAll    bool  `json:"queryAll"`
	QueryPublic bool  `json:"queryPublic"`
}

// SpaceTagAnalyzeRes 空间标签分析响应
type SpaceTagAnalyzeRes []SpaceTagAnalyzeResponse

// SpaceTagAnalyzeResponse 空间标签分析响应项
type SpaceTagAnalyzeResponse struct {
	Tag   string `json:"tag"`
	Count int64  `json:"count"`
}

// SpaceSizeAnalyzeReq 空间大小分析请求
type SpaceSizeAnalyzeReq struct {
	SpaceId     int64 `json:"spaceId"`
	QueryAll    bool  `json:"queryAll"`
	QueryPublic bool  `json:"queryPublic"`
}

// SpaceSizeAnalyzeRes 空间大小分析响应
type SpaceSizeAnalyzeRes []SpaceSizeAnalyzeResponse

// SpaceSizeAnalyzeResponse 空间大小分析响应项
type SpaceSizeAnalyzeResponse struct {
	SizeRange string `json:"sizeRange"`
	Count     int64  `json:"count"`
}

// SpaceUsageAnalyzeReq 空间使用情况分析请求
type SpaceUsageAnalyzeReq struct {
	SpaceId     int64 `json:"spaceId"`
	QueryAll    bool  `json:"queryAll"`
	QueryPublic bool  `json:"queryPublic"`
}

// SpaceUsageAnalyzeRes 空间使用情况分析响应
type SpaceUsageAnalyzeRes struct{ *SpaceUsageAnalyzeResponse }

// SpaceUsageAnalyzeResponse 空间使用情况分析响应项
type SpaceUsageAnalyzeResponse struct {
	MaxCount        int64   `json:"maxCount"`
	UsedCount       int64   `json:"usedCount"`
	CountUsageRatio float64 `json:"countUsageRatio"` // 0-100
	MaxSize         int64   `json:"maxSize"`
	UsedSize        int64   `json:"usedSize"`
	SizeUsageRatio  float64 `json:"sizeUsageRatio"` // 0-100
}

// SpaceUserAnalyzeReq 空间用户分析请求
type SpaceUserAnalyzeReq struct {
	SpaceId       int64  `json:"spaceId"`
	QueryAll      bool   `json:"queryAll"`
	QueryPublic   bool   `json:"queryPublic"`
	TimeDimension string `json:"timeDimension"` // day, week, month
}

// SpaceUserAnalyzeRes 空间用户分析响应
type SpaceUserAnalyzeRes []SpaceUserAnalyzeResponse

// SpaceUserAnalyzeResponse 空间用户分析响应项
type SpaceUserAnalyzeResponse struct {
	Period string `json:"period"`
	Count  int64  `json:"count"`
}

// SpaceRankAnalyzeReq 空间排行分析请求
type SpaceRankAnalyzeReq struct {
	TopN int `json:"topN" v:"between:1,100#排行数量必须在1-100之间"`
}

// SpaceRankAnalyzeRes 空间排行分析响应
type SpaceRankAnalyzeRes []Space
