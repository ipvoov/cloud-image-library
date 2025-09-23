package user

import (
	v1 "cloud/api/user/v1"
	"cloud/internal/consts"
	"cloud/internal/model/entity"
	"crypto/md5"
	"fmt"
)

// encryptPassword 密码加密
func (s *sUser) encryptPassword(password string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(consts.Salt+password)))
}

// getLoginUserVO 获取登录用户VO
func (s *sUser) getLoginUserVO(user *entity.User) *v1.LoginUserVO {
	if user == nil {
		return nil
	}

	return &v1.LoginUserVO{
		Id:          user.Id,
		UserName:    user.UserName,
		UserAccount: user.UserAccount,
		UserAvatar:  user.UserAvatar,
		UserRole:    user.UserRole,
		CreateTime:  user.CreateTime.String(),
		UpdateTime:  user.UpdateTime.String(),
	}
}
