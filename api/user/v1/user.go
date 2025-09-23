package v1

// UserRegisterReq 用户注册请求
// UserRegisterReq 用户注册请求
type UserRegisterReq struct {
	UserAccount   string `json:"userAccount" v:"required#请输入账号|length:4,20|账号长度为4-20位"`
	UserPassword  string `json:"userPassword" v:"required#请输入密码|length:8,20|密码长度为8-20位"`
	CheckPassword string `json:"checkPassword" v:"required#请输入确认密码|length:8,20|确认密码长度为8-20位"`
}

// UserRegisterRes 用户注册响应
type UserRegisterRes struct {
	UserId int64 `json:"userId"`
}

// UserLoginReq 用户登录请求
type UserLoginReq struct {
	UserAccount  string `json:"userAccount" v:"required#请输入账号|length:4,20|账号长度为4-20位"`
	UserPassword string `json:"userPassword" v:"required#请输入密码|length:8,20|密码长度为8-20位"`
}

// UserLoginRes 用户登录响应
type UserLoginRes struct {
	*LoginUserVO
}

// LoginUserVO 登录用户信息
type LoginUserVO struct {
	Id          int64  `json:"id"`
	UserName    string `json:"userName"`
	UserAccount string `json:"userAccount"`
	UserAvatar  string `json:"userAvatar"`
	UserRole    string `json:"userRole"`
	CreateTime  string `json:"createTime"`
	UpdateTime  string `json:"updateTime"`
}

// GetLoginUserReq 获取登录用户请求
type GetLoginUserReq struct{}

// GetLoginUserRes 获取登录用户响应
type GetLoginUserRes struct {
	*LoginUserVO
}

// UserLogoutReq 用户登出请求
type UserLogoutReq struct{}

// UserLogoutRes 用户登出响应
type UserLogoutRes struct {
	Success bool `json:"success"`
}

// UserAddReq 添加用户请求
type UserAddReq struct {
	UserAccount  string `json:"userAccount" v:"required#请输入账号|length:4,20|账号长度为4-20位"`
	UserPassword string `json:"userPassword" v:"required#请输入密码|length:8,20|密码长度为8-20位"`
	UserName     string `json:"userName" v:"length:1,20|用户名长度为1-20位"`
	UserRole     string `json:"userRole" v:"in:user,admin#用户角色只能是user或admin"`
}

// UserAddRes 添加用户响应
type UserAddRes struct {
	UserId int64 `json:"userId"`
}

// UserUpdateReq 更新用户请求
type UserUpdateReq struct {
	Id         int64  `json:"id" v:"required#请输入用户ID"`
	UserName   string `json:"userName" v:"length:1,20|用户名长度为1-20位"`
	UserAvatar string `json:"userAvatar"`
	UserRole   string `json:"userRole" v:"in:user,admin#用户角色只能是user或admin"`
}

// UserUpdateRes 更新用户响应
type UserUpdateRes struct {
	Success bool `json:"success"`
}

// // DeleteReq 删除请求
// type DeleteReq struct {
// 	Id int64 `json:"id" v:"required#请输入ID"`
// }

// // DeleteRes 删除响应
// type DeleteRes struct {
// 	Success bool `json:"success"`
// }

// GetUserByIdReq 根据ID获取用户请求
type GetUserByIdReq struct {
	Id int64 `json:"id" v:"required#请输入用户ID"`
}

// GetUserByIdRes 根据ID获取用户响应
type GetUserByIdRes struct {
	*LoginUserVO
}

// UserQueryReq 用户查询请求
type UserQueryReq struct {
	Current     int    `json:"current" v:"min:1#页码最小为1"`
	PageSize    int    `json:"pageSize" v:"between:1,100#页面大小为1-100"`
	SortField   string `json:"sortField"`
	SortOrder   string `json:"sortOrder"`
	UserAccount string `json:"userAccount"`
	UserName    string `json:"userName"`
	UserRole    string `json:"userRole"`
}

// UserQueryRes 用户查询响应
type UserQueryRes struct {
	Records []LoginUserVO `json:"records"`
	Total   int           `json:"total"`
	Current int           `json:"current"`
	Size    int           `json:"size"`
}

// GetUserProfileReq 获取用户详细信息请求
type GetUserProfileReq struct{}

// GetUserProfileRes 获取用户详细信息响应
type GetUserProfileRes struct {
	*UserProfileVO
}

// UserProfileVO 用户详细信息VO
type UserProfileVO struct {
	Id            int64  `json:"id"`
	UserAccount   string `json:"userAccount"`
	UserName      string `json:"userName"`
	UserAvatar    string `json:"userAvatar"`
	UserProfile   string `json:"userProfile"`
	UserRole      string `json:"userRole"`
	CreateTime    string `json:"createTime"`
	UpdateTime    string `json:"updateTime"`
	VipExpireTime string `json:"vipExpireTime"`
	VipCode       string `json:"vipCode"`
	VipNumber     int64  `json:"vipNumber"`
	// 统计信息
	PictureCount int64 `json:"pictureCount"` // 上传图片数量
	SpaceCount   int64 `json:"spaceCount"`   // 创建空间数量
	TeamCount    int64 `json:"teamCount"`    // 加入团队数量
}

// UpdateUserProfileReq 更新用户资料请求
type UpdateUserProfileReq struct {
	UserName    string `json:"userName" v:"length:1,20|用户名长度为1-20位"`
	UserAvatar  string `json:"userAvatar"`
	UserProfile string `json:"userProfile" v:"max-length:512|个人简介最多512个字符"`
}

// UpdateUserProfileRes 更新用户资料响应
type UpdateUserProfileRes struct {
	Success bool `json:"success"`
}
