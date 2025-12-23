// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Files is the golang structure of table files for DAO operations like Where/Data.
type Files struct {
	g.Meta               `orm:"table:files, do:true"`
	Id                   any         // 主键
	UserId               any         // 归属用户ID（无外键）
	TelegramAccountId    any         // 使用的 Telegram 账号ID（无外键）
	TelegramFileId       any         // Telegram file_id（下载用）
	TelegramFileUniqueId any         // Telegram file_unique_id
	TelegramMessageId    any         // 发送到 Telegram 的消息ID
	TelegramChatId       any         // 发送到的 chat_id
	FileName             any         // 原始文件名
	SizeBytes            any         // 文件大小（字节）
	MimeType             any         // MIME 类型
	Md5                  any         // MD5 去重
	Sha256               any         // sha256 校验/去重
	StorageProvider      any         // 存储提供方标识，默认 telegram
	ObjectKey            any         // 存储对象 key/路径，默认使用 Telegram file_id
	Status               any         // 状态：1=ok,2=uploading,3=failed
	IsDeleted            any         // 是否删除（软删）
	DeletedAt            *gtime.Time // 软删时间
	CreatedAt            *gtime.Time // 创建时间
	UpdatedAt            *gtime.Time // 更新时间
}
