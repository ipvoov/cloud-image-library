package v1

// PictureTagCategoryReq 获取图片标签分类请求
type PictureTagCategoryReq struct{}

// PictureTagCategoryRes 获取图片标签分类响应
type PictureTagCategoryRes struct {
	TagList      []string `json:"tagList"`
	CategoryList []string `json:"categoryList"`
}

// PictureEditReq 图片编辑请求
type PictureEditReq struct {
	Id           int64    `json:"id" v:"required#图片ID不能为空"`
	Name         string   `json:"name"`
	Introduction string   `json:"introduction"`
	Category     string   `json:"category"`
	Tags         []string `json:"tags"`
	SpaceId      int64    `json:"spaceId"`
}

// PictureEditRes 图片编辑响应
type PictureEditRes struct {
	Success bool `json:"success"`
}

// PictureUploadReq 图片上传请求
type PictureUploadReq struct {
	SpaceId int64 `json:"spaceId"`
}

// PictureUploadByUrlReq 通过URL上传图片请求
type PictureUploadByUrlReq struct {
	FileUrl  string `json:"fileUrl" v:"required#文件URL不能为空"`
	FileName string `json:"fileName"`
	SpaceId  int64  `json:"spaceId"`
}

// PictureVO 图片视图对象（用户视图）
type PictureVO struct {
	Id             int64    `json:"id"`
	Url            string   `json:"url"`
	Name           string   `json:"name"`
	Introduction   string   `json:"introduction"`
	Category       string   `json:"category"`
	Tags           []string `json:"tags"`
	PicSize        int64    `json:"picSize"`
	PicWidth       int      `json:"picWidth"`
	PicHeight      int      `json:"picHeight"`
	PicScale       float64  `json:"picScale"`
	PicFormat      string   `json:"picFormat"`
	UserId         int64    `json:"userId"`
	SpaceId        int64    `json:"spaceId"`
	CreateTime     string   `json:"createTime"`
	EditTime       string   `json:"editTime"`
	UpdateTime     string   `json:"updateTime"`
	ThumbnailUrl   string   `json:"thumbnailUrl"`
	PicColor       string   `json:"picColor"`
	User           *UserVO  `json:"user,omitempty"`
	PermissionList []string `json:"permissionList,omitempty"`
}

// Picture 图片实体对象（管理员视图，包含审核信息）
type Picture struct {
	Id            int64   `json:"id"`
	Url           string  `json:"url"`
	Name          string  `json:"name"`
	Introduction  string  `json:"introduction"`
	Category      string  `json:"category"`
	Tags          string  `json:"tags"` // 注意：前端期望JSON字符串格式 "[\"标签1\",\"标签2\"]"
	PicSize       int64   `json:"picSize"`
	PicWidth      int     `json:"picWidth"`
	PicHeight     int     `json:"picHeight"`
	PicScale      float64 `json:"picScale"`
	PicFormat     string  `json:"picFormat"`
	UserId        int64   `json:"userId"`
	SpaceId       int64   `json:"spaceId"`
	CreateTime    string  `json:"createTime"`
	EditTime      string  `json:"editTime"`
	UpdateTime    string  `json:"updateTime"`
	ThumbnailUrl  string  `json:"thumbnailUrl"`
	PicColor      string  `json:"picColor"`
	IsDelete      int     `json:"isDelete"`      // 删除状态
	ReviewStatus  int     `json:"reviewStatus"`  // 审核状态：0-待审核，1-通过，2-拒绝
	ReviewMessage string  `json:"reviewMessage"` // 审核信息
	ReviewerId    int64   `json:"reviewerId"`    // 审核人ID
	ReviewTime    string  `json:"reviewTime"`    // 审核时间
}

// PictureUploadRes 图片上传响应
type PictureUploadRes struct {
	*PictureVO
}

// PictureUploadByUrlRes 通过URL上传图片响应
type PictureUploadByUrlRes struct {
	*PictureVO
}

// DeleteReq 删除请求
type DeleteReq struct {
	Id int64 `json:"id" v:"required#ID不能为空"`
}

// DeleteRes 删除响应
type DeleteRes struct {
	Success bool `json:"success"`
}

// PictureGetReq 获取图片请求
type PictureGetReq struct {
	Id int64 `json:"id" v:"required#图片ID不能为空"`
}

// PictureGetRes 获取图片响应
type PictureGetRes struct {
	*PictureVO
}

// PictureGetRes 获取图片响应
type PictureAdminGetRes struct {
	*Picture
}

// PictureQueryReq 图片查询请求
type PictureQueryReq struct {
	Current      int      `json:"current" p:"current" v:"min:1#页码最小为1"`
	PageSize     int      `json:"pageSize" p:"pageSize" v:"between:1,100#页面大小为1-100"`
	Tags         []string `json:"tags" p:"tags"`
	SpaceId      string   `json:"spaceId" p:"spaceId"`
	SearchText   string   `json:"searchText" p:"searchText"`
	ReviewStatus *int     `json:"reviewStatus" p:"reviewStatus" v:"in:0,1,2" dc:"0:待审核;1:审核通过;2:审核未通过"`
	Category     string   `json:"category" p:"category"`
	SortField    string   `json:"sortField" p:"sortField"`
	SortOrder    string   `json:"sortOrder" p:"sortOrder"`
}

type PageInfo struct {
	Current int `json:"current"`
	Size    int `json:"size"`
	Total   int `json:"total"`
	Pages   int `json:"pages"`
}

// PictureQueryRes 图片查询响应
type PictureQueryRes struct {
	Records []PictureVO `json:"records"`
	*PageInfo
}

// PictureAdminQueryRes 图片管理查询响应（管理员视图）
type PictureAdminQueryRes struct {
	Records []Picture `json:"records"`
	*PageInfo
}

// PictureUpdateReq 图片更新请求
type PictureUpdateReq struct {
	Id           int64    `json:"id" v:"required#图片ID不能为空"`
	Name         string   `json:"name"`
	Introduction string   `json:"introduction"`
	Category     string   `json:"category"`
	Tags         []string `json:"tags"`
}

// PictureUpdateRes 图片更新响应
type PictureUpdateRes struct {
	Success bool `json:"success"`
}

// PictureReviewReq 图片审核请求
type PictureReviewReq struct {
	Id            int64  `json:"id" v:"required#图片ID不能为空"`
	ReviewStatus  int    `json:"reviewStatus" v:"required|in:0,1,2#审核状态不能为空且必须为0,1,2"`
	ReviewMessage string `json:"reviewMessage" v:"required#审核信息不能为空"`
}

// PictureReviewRes 图片审核响应
type PictureReviewRes struct {
	Success bool `json:"success"`
}

// PictureUploadByBatchReq 批量上传图片请求
type PictureUploadByBatchReq struct {
	SearchText string `json:"searchText" v:"required#搜索关键词不能为空"`
	Count      int    `json:"count" v:"required|between:1,30#抓取数量不能为空且必须在1-30之间"`
	NamePrefix string `json:"namePrefix"`
	SpaceId    int64  `json:"spaceId"`
}

// PictureUploadByBatchRes 批量上传图片响应
type PictureUploadByBatchRes struct {
	SuccessCount int `json:"successCount"`
}

// PictureEditByBatchReq 批量编辑图片请求
type PictureEditByBatchReq struct {
	PictureIdList []int64  `json:"pictureIdList" v:"required#图片ID列表不能为空"`
	Category      string   `json:"category"`
	Tags          []string `json:"tags"`
	NameRule      string   `json:"nameRule"`
	SpaceId       int64    `json:"spaceId"`
}

// PictureEditByBatchRes 批量编辑图片响应
type PictureEditByBatchRes struct {
	Success bool `json:"success"`
}

// SearchPictureByPictureReq 以图搜图请求
type SearchPictureByPictureReq struct {
	PictureId int64 `json:"pictureId" v:"required#图片ID不能为空"`
}

// SearchPictureByPictureRes 以图搜图响应
type SearchPictureByPictureRes struct {
	FromUrl  string `json:"fromUrl"`  // 图片来源URL
	ThumbUrl string `json:"thumbUrl"` // 缩略图URL
}

// SearchPictureByColorReq 按颜色搜索图片请求
type SearchPictureByColorReq struct {
	PicColor string `json:"picColor" v:"required#颜色值不能为空"`
	SpaceId  int64  `json:"spaceId"`
}

// SearchPictureByColorRes 按颜色搜索图片响应
type SearchPictureByColorRes struct {
	*PictureVO
}

// CreatePictureOutPaintingReq 创建图片扩图请求
type CreatePictureOutPaintingReq struct {
	PictureId  int64                  `json:"pictureId" v:"required#图片ID不能为空"`
	Prompt     string                 `json:"prompt" v:"required#扩图描述不能为空"`
	Parameters *OutPaintingParameters `json:"parameters"`
}

// OutPaintingParameters 扩图参数
type OutPaintingParameters struct {
	XScale int `json:"xScale"`
	YScale int `json:"yScale"`
}

// CreatePictureOutPaintingRes 创建图片扩图响应
type CreatePictureOutPaintingRes struct {
	OutputImageUrl string `json:"outputImageUrl"` // 输出图片URL
	Success        bool   `json:"success"`        // 是否成功
	Message        string `json:"message"`        // 消息
}

// CreatePictureAIEditingTaskReq AI编辑图片任务创建请求
type CreatePictureAIEditingTaskReq struct {
	PictureId int64  `json:"pictureId" v:"required#图片ID不能为空"`
	Prompt    string `json:"prompt" v:"required#编辑描述不能为空"`
}

// CreatePictureAIEditingTaskRes AI编辑图片任务创建响应
type CreatePictureAIEditingTaskRes struct {
	Output    *CreateAIEditingTaskResponse `json:"output"`
	RequestId string                       `json:"requestId"`
}

// CreateAIEditingTaskResponse AI编辑任务响应
type CreateAIEditingTaskResponse struct {
	TaskId string `json:"taskId"`
}

// GetPictureAIEditingTaskReq 获取AI编辑任务请求
type GetPictureAIEditingTaskReq struct {
	TaskId string `json:"taskId" v:"required#任务ID不能为空"`
}

// GetPictureAIEditingTaskRes 获取AI编辑任务响应
type GetPictureAIEditingTaskRes struct {
	Output    *GetAIEditingTaskResponse `json:"output"`
	RequestId string                    `json:"requestId"`
}

// GetAIEditingTaskResponse 获取AI编辑任务响应
type GetAIEditingTaskResponse struct {
	TaskId         string `json:"taskId"`
	TaskStatus     string `json:"taskStatus"`     // SUCCEEDED, FAILED, RUNNING
	OutputImageUrl string `json:"outputImageUrl"` // 输出图片URL
	ErrorMessage   string `json:"errorMessage"`   // 错误信息
}
