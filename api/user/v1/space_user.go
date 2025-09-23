package v1

// SpaceUserListMyReq 获取我的团队空间请求
type SpaceUserListMyReq struct{}

// SpaceUserVO 空间用户视图对象
type SpaceUserVO struct {
	Id         int64    `json:"id"`
	UserId     int64    `json:"userId"`
	SpaceId    int64    `json:"spaceId"`
	SpaceRole  string   `json:"spaceRole"`
	CreateTime string   `json:"createTime"`
	UpdateTime string   `json:"updateTime"`
	User       *UserVO  `json:"user,omitempty"`
	Space      *SpaceVO `json:"space,omitempty"`
}

// UserVO 用户视图对象
type UserVO struct {
	Id         int64  `json:"id"`
	UserName   string `json:"userName"`
	UserAvatar string `json:"userAvatar"`
}

// SpaceUserListMyRes 获取我的团队空间响应
type SpaceUserListMyRes []SpaceUserVO

// SpaceUserAddReq 添加空间用户请求
type SpaceUserAddReq struct {
	SpaceId int64 `json:"spaceId" v:"required#空间ID不能为空"`
	UserId  int64 `json:"userId" v:"required#用户ID不能为空"`
	// SpaceRole string `json:"spaceRole" v:"required|in:viewer,editor,admin#空间角色必须为viewer,editor或admin"`
}

// SpaceUserAddRes 添加空间用户响应
type SpaceUserAddRes struct {
	Id int64 `json:"id"`
}

// SpaceUserEditReq 编辑空间用户请求
type SpaceUserEditReq struct {
	Id        int64  `json:"id" v:"required#空间用户ID不能为空"`
	SpaceRole string `json:"spaceRole" v:"required|in:viewer,editor,admin#空间角色必须为viewer,editor或admin"`
}

// SpaceUserEditRes 编辑空间用户响应
type SpaceUserEditRes struct {
	Success bool `json:"success"`
}

// SpaceUserDeleteReq 删除空间用户请求
type SpaceUserDeleteReq struct {
	Id int64 `json:"id" v:"required#空间用户ID不能为空"`
}

// SpaceUserDeleteRes 删除空间用户响应
type SpaceUserDeleteRes struct {
	Success bool `json:"success"`
}

// SpaceUserGetReq 获取空间用户请求
type SpaceUserGetReq struct {
	Id      int64 `json:"id"`
	SpaceId int64 `json:"spaceId"`
	UserId  int64 `json:"userId"`
}

// SpaceUserGetRes 获取空间用户响应
type SpaceUserGetRes struct {
	*SpaceUser
}

// SpaceUserListReq 获取空间用户列表请求
type SpaceUserListReq struct {
	SpaceId int64 `json:"spaceId" v:"required#空间ID不能为空"`
}

// SpaceUserListRes 获取空间用户列表响应
type SpaceUserListRes []SpaceUserVO

// SpaceUser 空间用户实体对象
type SpaceUser struct {
	Id         int64  `json:"id"`
	UserId     int64  `json:"userId"`
	SpaceId    int64  `json:"spaceId"`
	SpaceRole  string `json:"spaceRole"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}
