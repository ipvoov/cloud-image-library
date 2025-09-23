package v1

// SpaceAddReq 添加空间请求
type SpaceAddReq struct {
	SpaceName  string `json:"spaceName" v:"required#空间名称不能为空"`
	SpaceLevel int    `json:"spaceLevel" v:"required|in:0,1,2#空间级别必须为0,1,2"`
	SpaceType  int    `json:"spaceType" v:"required|in:0,1#空间类型必须为0,1"`
}

// SpaceAddRes 添加空间响应
type SpaceAddRes struct {
	Id int64 `json:"id"`
}

// SpaceEditReq 编辑空间请求
type SpaceEditReq struct {
	Id        int64  `json:"id" v:"required#空间ID不能为空"`
	SpaceName string `json:"spaceName" v:"required#空间名称不能为空"`
}

// SpaceEditRes 编辑空间响应
type SpaceEditRes struct {
	Success bool `json:"success"`
}

// SpaceUpdateReq 更新空间请求
type SpaceUpdateReq struct {
	Id       int64 `json:"id" v:"required#空间ID不能为空"`
	MaxCount int64 `json:"maxCount"`
	MaxSize  int64 `json:"maxSize"`
}

// SpaceUpdateRes 更新空间响应
type SpaceUpdateRes struct {
	Success bool `json:"success"`
}

// SpaceGetReq 获取空间请求
type SpaceGetReq struct {
	Id int64 `json:"id" v:"required#空间ID不能为空"`
}

// SpaceGetRes 获取空间响应
type SpaceGetRes struct {
	*Space
}

// SpaceGetVOReq 获取空间VO请求
type SpaceGetVOReq struct {
	Id string `json:"id" v:"required#空间ID不能为空"`
}

// SpaceGetVORes 获取空间VO响应
type SpaceGetVORes struct {
	*SpaceVO
}

// SpaceQueryReq 空间查询请求
type SpaceQueryReq struct {
	Current   int    `json:"current" p:"current" v:"min:1#页码最小为1"`
	PageSize  int    `json:"pageSize" p:"pageSize" v:"between:1,100#页面大小为1-100"`
	SpaceId   string `json:"spaceId" p:"spaceId"`
	SpaceName string `json:"spaceName" p:"spaceName"`
	UserId    int64  `json:"userId" p:"userId"`
}

// SpaceQueryRes 空间查询响应
type SpaceQueryRes struct {
	Records []Space `json:"records"`
	*PageInfo
}

// SpaceQueryVORes 空间查询VO响应
type SpaceQueryVORes struct {
	Records []SpaceVO `json:"records"`
	*PageInfo
}

// SpaceDeleteReq 删除空间请求
type SpaceDeleteReq struct {
	Id int64 `json:"id" v:"required#空间ID不能为空"`
}

// SpaceDeleteRes 删除空间响应
type SpaceDeleteRes struct {
	Success bool `json:"success"`
}

// SpaceLevelListReq 获取空间级别列表请求
type SpaceLevelListReq struct{}

// SpaceLevel 空间级别
type SpaceLevel struct {
	Level       int    `json:"level"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MaxSize     int64  `json:"maxSize"`
	MaxCount    int64  `json:"maxCount"`
}

// SpaceLevelListRes 获取空间级别列表响应
type SpaceLevelListRes struct {
	Records []SpaceLevel `json:"records"`
}

// Space 空间实体对象（管理员视图）
type Space struct {
	Id         int64  `json:"id"`
	SpaceName  string `json:"spaceName"`
	SpaceLevel int    `json:"spaceLevel"`
	MaxSize    int64  `json:"maxSize"`
	MaxCount   int64  `json:"maxCount"`
	TotalSize  int64  `json:"totalSize"`
	TotalCount int64  `json:"totalCount"`
	UserId     int64  `json:"userId"`
	CreateTime string `json:"createTime"`
	EditTime   string `json:"editTime"`
	UpdateTime string `json:"updateTime"`
	IsDelete   int    `json:"isDelete"`
	SpaceType  int    `json:"spaceType"`
}

// SpaceVO 空间视图对象（用户视图）
type SpaceVO struct {
	Id             int64    `json:"id"`
	SpaceName      string   `json:"spaceName"`
	SpaceLevel     int      `json:"spaceLevel"`
	MaxSize        int64    `json:"maxSize"`
	MaxCount       int64    `json:"maxCount"`
	TotalSize      int64    `json:"totalSize"`
	TotalCount     int64    `json:"totalCount"`
	UserId         int64    `json:"userId"`
	CreateTime     string   `json:"createTime"`
	EditTime       string   `json:"editTime"`
	UpdateTime     string   `json:"updateTime"`
	SpaceType      int      `json:"spaceType"`
	PermissionList []string `json:"permissionList,omitempty"`
	User           *UserVO  `json:"user,omitempty"`
}
