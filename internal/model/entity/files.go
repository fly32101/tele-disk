// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Files is the golang structure for table files.
type Files struct {
	Id                   int64       `json:"id"                   orm:"id"                      description:"主键"`                                // 主键
	UserId               int64       `json:"userId"               orm:"user_id"                 description:"归属用户ID（无外键）"`                       // 归属用户ID（无外键）
	TelegramAccountId    int64       `json:"telegramAccountId"    orm:"telegram_account_id"     description:"使用的 Telegram 账号ID（无外键）"`            // 使用的 Telegram 账号ID（无外键）
	TelegramFileId       string      `json:"telegramFileId"       orm:"telegram_file_id"        description:"Telegram file_id（下载用）"`             // Telegram file_id（下载用）
	TelegramFileUniqueId string      `json:"telegramFileUniqueId" orm:"telegram_file_unique_id" description:"Telegram file_unique_id"`           // Telegram file_unique_id
	TelegramMessageId    int64       `json:"telegramMessageId"    orm:"telegram_message_id"     description:"发送到 Telegram 的消息ID"`                // 发送到 Telegram 的消息ID
	TelegramChatId       int64       `json:"telegramChatId"       orm:"telegram_chat_id"        description:"发送到的 chat_id"`                      // 发送到的 chat_id
	FileName             string      `json:"fileName"             orm:"file_name"               description:"原始文件名"`                             // 原始文件名
	SizeBytes            int64       `json:"sizeBytes"            orm:"size_bytes"              description:"文件大小（字节）"`                          // 文件大小（字节）
	MimeType             string      `json:"mimeType"             orm:"mime_type"               description:"MIME 类型"`                           // MIME 类型
	Md5                  string      `json:"md5"                  orm:"md5"                     description:"MD5 去重"`                            // MD5 去重
	Sha256               string      `json:"sha256"               orm:"sha256"                  description:"sha256 校验/去重"`                      // sha256 校验/去重
	StorageProvider      string      `json:"storageProvider"      orm:"storage_provider"        description:"存储提供方标识，默认 telegram"`               // 存储提供方标识，默认 telegram
	ObjectKey            string      `json:"objectKey"            orm:"object_key"              description:"存储对象 key/路径，默认使用 Telegram file_id"` // 存储对象 key/路径，默认使用 Telegram file_id
	Status               int         `json:"status"               orm:"status"                  description:"状态：1=ok,2=uploading,3=failed"`      // 状态：1=ok,2=uploading,3=failed
	IsDeleted            int         `json:"isDeleted"            orm:"is_deleted"              description:"是否删除（软删）"`                          // 是否删除（软删）
	DeletedAt            *gtime.Time `json:"deletedAt"            orm:"deleted_at"              description:"软删时间"`                              // 软删时间
	CreatedAt            *gtime.Time `json:"createdAt"            orm:"created_at"              description:"创建时间"`                              // 创建时间
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"              description:"更新时间"`                              // 更新时间
}
