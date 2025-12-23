// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// TelegramAccounts is the golang structure of table telegram_accounts for DAO operations like Where/Data.
type TelegramAccounts struct {
	g.Meta      `orm:"table:telegram_accounts, do:true"`
	Id          any         // 主键
	UserId      any         // 所属用户ID（无外键）
	Name        any         // 账号名称
	BotUsername any         // Bot 用户名（@xxx）
	BotTokenEnc any         // 加密后的 Bot Token（对称加密，可逆）
	BotId       any         // Bot ID
	ChatId      any         // 存储文件用的 chat_id
	Status      any         // 状态：1启用，0停用
	IsDefault   any         // 是否默认上传账号
	LastError   any         // 最近一次错误信息
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
}
