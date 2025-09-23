// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SpaceUserDao is the data access object for the table space_user.
type SpaceUserDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SpaceUserColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SpaceUserColumns defines and stores column names for the table space_user.
type SpaceUserColumns struct {
	Id         string // id
	SpaceId    string // 空间 id
	UserId     string // 用户 id
	SpaceRole  string // 空间角色：viewer/editor/admin
	CreateTime string // 创建时间
	UpdateTime string // 更新时间
}

// spaceUserColumns holds the columns for the table space_user.
var spaceUserColumns = SpaceUserColumns{
	Id:         "id",
	SpaceId:    "spaceId",
	UserId:     "userId",
	SpaceRole:  "spaceRole",
	CreateTime: "createTime",
	UpdateTime: "updateTime",
}

// NewSpaceUserDao creates and returns a new DAO object for table data access.
func NewSpaceUserDao(handlers ...gdb.ModelHandler) *SpaceUserDao {
	return &SpaceUserDao{
		group:    "default",
		table:    "space_user",
		columns:  spaceUserColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SpaceUserDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SpaceUserDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SpaceUserDao) Columns() SpaceUserColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SpaceUserDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SpaceUserDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SpaceUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
