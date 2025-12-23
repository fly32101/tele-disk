// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FilesDao is the data access object for the table files.
type FilesDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  FilesColumns       // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// FilesColumns defines and stores column names for the table files.
type FilesColumns struct {
	Id                   string // 主键
	UserId               string // 归属用户ID（无外键）
	TelegramAccountId    string // 使用的 Telegram 账号ID（无外键）
	TelegramFileId       string // Telegram file_id（下载用）
	TelegramFileUniqueId string // Telegram file_unique_id
	TelegramMessageId    string // 发送到 Telegram 的消息ID
	TelegramChatId       string // 发送到的 chat_id
	FileName             string // 原始文件名
	SizeBytes            string // 文件大小（字节）
	MimeType             string // MIME 类型
	Md5                  string // MD5 去重
	Sha256               string // sha256 校验/去重
	StorageProvider      string // 存储提供方标识，默认 telegram
	ObjectKey            string // 存储对象 key/路径，默认使用 Telegram file_id
	Status               string // 状态：1=ok,2=uploading,3=failed
	IsDeleted            string // 是否删除（软删）
	DeletedAt            string // 软删时间
	CreatedAt            string // 创建时间
	UpdatedAt            string // 更新时间
}

// filesColumns holds the columns for the table files.
var filesColumns = FilesColumns{
	Id:                   "id",
	UserId:               "user_id",
	TelegramAccountId:    "telegram_account_id",
	TelegramFileId:       "telegram_file_id",
	TelegramFileUniqueId: "telegram_file_unique_id",
	TelegramMessageId:    "telegram_message_id",
	TelegramChatId:       "telegram_chat_id",
	FileName:             "file_name",
	SizeBytes:            "size_bytes",
	MimeType:             "mime_type",
	Md5:                  "md5",
	Sha256:               "sha256",
	StorageProvider:      "storage_provider",
	ObjectKey:            "object_key",
	Status:               "status",
	IsDeleted:            "is_deleted",
	DeletedAt:            "deleted_at",
	CreatedAt:            "created_at",
	UpdatedAt:            "updated_at",
}

// NewFilesDao creates and returns a new DAO object for table data access.
func NewFilesDao(handlers ...gdb.ModelHandler) *FilesDao {
	return &FilesDao{
		group:    "default",
		table:    "files",
		columns:  filesColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *FilesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *FilesDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *FilesDao) Columns() FilesColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *FilesDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *FilesDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *FilesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
