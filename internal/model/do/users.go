// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure of table users for DAO operations like Where/Data.
type Users struct {
	g.Meta       `orm:"table:users, do:true"`
	Id           any         // 主键
	Username     any         // 登录用户名
	PasswordHash any         // 密码哈希
	IsAdmin      any         // 是否管理员
	QuotaBytes   any         // 配额，0 表示不限
	UsedBytes    any         // 已用容量
	FileCount    any         // 文件数量
	Status       any         // 状态：1启用，0停用
	LastLoginAt  *gtime.Time // 上次登录时间
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time // 更新时间
}
