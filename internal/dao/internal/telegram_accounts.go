// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TelegramAccountsDao is the data access object for the table telegram_accounts.
type TelegramAccountsDao struct {
	table    string                  // table is the underlying table name of the DAO.
	group    string                  // group is the database configuration group name of the current DAO.
	columns  TelegramAccountsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler      // handlers for customized model modification.
}

// TelegramAccountsColumns defines and stores column names for the table telegram_accounts.
type TelegramAccountsColumns struct {
	Id          string // 主键
	UserId      string // 所属用户ID（无外键）
	Name        string // 账号名称
	BotUsername string // Bot 用户名（@xxx）
	BotTokenEnc string // 加密后的 Bot Token（对称加密，可逆）
	BotId       string // Bot ID
	ChatId      string // 存储文件用的 chat_id
	Status      string // 状态：1启用，0停用
	IsDefault   string // 是否默认上传账号
	LastError   string // 最近一次错误信息
	CreatedAt   string // 创建时间
	UpdatedAt   string // 更新时间
}

// telegramAccountsColumns holds the columns for the table telegram_accounts.
var telegramAccountsColumns = TelegramAccountsColumns{
	Id:          "id",
	UserId:      "user_id",
	Name:        "name",
	BotUsername: "bot_username",
	BotTokenEnc: "bot_token_enc",
	BotId:       "bot_id",
	ChatId:      "chat_id",
	Status:      "status",
	IsDefault:   "is_default",
	LastError:   "last_error",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewTelegramAccountsDao creates and returns a new DAO object for table data access.
func NewTelegramAccountsDao(handlers ...gdb.ModelHandler) *TelegramAccountsDao {
	return &TelegramAccountsDao{
		group:    "default",
		table:    "telegram_accounts",
		columns:  telegramAccountsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TelegramAccountsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TelegramAccountsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TelegramAccountsDao) Columns() TelegramAccountsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TelegramAccountsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TelegramAccountsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *TelegramAccountsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
