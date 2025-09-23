// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserDao is the data access object for the table user.
type UserDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UserColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UserColumns defines and stores column names for the table user.
type UserColumns struct {
	Id            string // id
	UserAccount   string // 账号
	UserPassword  string // 密码
	UserName      string // 用户昵称
	UserAvatar    string // 用户头像
	UserProfile   string // 用户简介
	UserRole      string // 用户角色：user/admin
	EditTime      string // 编辑时间
	CreateTime    string // 创建时间
	UpdateTime    string // 更新时间
	IsDelete      string // 是否删除
	VipExpireTime string // 会员过期时间
	VipCode       string // 会员兑换码
	VipNumber     string // 会员编号
}

// userColumns holds the columns for the table user.
var userColumns = UserColumns{
	Id:            "id",
	UserAccount:   "userAccount",
	UserPassword:  "userPassword",
	UserName:      "userName",
	UserAvatar:    "userAvatar",
	UserProfile:   "userProfile",
	UserRole:      "userRole",
	EditTime:      "editTime",
	CreateTime:    "createTime",
	UpdateTime:    "updateTime",
	IsDelete:      "isDelete",
	VipExpireTime: "vipExpireTime",
	VipCode:       "vipCode",
	VipNumber:     "vipNumber",
}

// NewUserDao creates and returns a new DAO object for table data access.
func NewUserDao(handlers ...gdb.ModelHandler) *UserDao {
	return &UserDao{
		group:    "default",
		table:    "user",
		columns:  userColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserDao) Columns() UserColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *UserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
