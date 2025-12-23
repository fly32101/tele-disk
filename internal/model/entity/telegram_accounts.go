// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// TelegramAccounts is the golang structure for table telegram_accounts.
type TelegramAccounts struct {
	Id          int64       `json:"id"          orm:"id"            description:"主键"`                      // 主键
	UserId      int64       `json:"userId"      orm:"user_id"       description:"所属用户ID（无外键）"`             // 所属用户ID（无外键）
	Name        string      `json:"name"        orm:"name"          description:"账号名称"`                    // 账号名称
	BotUsername string      `json:"botUsername" orm:"bot_username"  description:"Bot 用户名（@xxx）"`           // Bot 用户名（@xxx）
	BotTokenEnc string      `json:"botTokenEnc" orm:"bot_token_enc" description:"加密后的 Bot Token（对称加密，可逆）"` // 加密后的 Bot Token（对称加密，可逆）
	BotId       int64       `json:"botId"       orm:"bot_id"        description:"Bot ID"`                  // Bot ID
	ChatId      int64       `json:"chatId"      orm:"chat_id"       description:"存储文件用的 chat_id"`          // 存储文件用的 chat_id
	Status      int         `json:"status"      orm:"status"        description:"状态：1启用，0停用"`              // 状态：1启用，0停用
	IsDefault   int         `json:"isDefault"   orm:"is_default"    description:"是否默认上传账号"`                // 是否默认上传账号
	LastError   string      `json:"lastError"   orm:"last_error"    description:"最近一次错误信息"`                // 最近一次错误信息
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"    description:"创建时间"`                    // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"    description:"更新时间"`                    // 更新时间
}
