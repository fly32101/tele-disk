// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure for table users.
type Users struct {
	Id           int64       `json:"id"           orm:"id"            description:"主键"`         // 主键
	Username     string      `json:"username"     orm:"username"      description:"登录用户名"`      // 登录用户名
	PasswordHash string      `json:"passwordHash" orm:"password_hash" description:"密码哈希"`       // 密码哈希
	IsAdmin      int         `json:"isAdmin"      orm:"is_admin"      description:"是否管理员"`      // 是否管理员
	QuotaBytes   int64       `json:"quotaBytes"   orm:"quota_bytes"   description:"配额，0 表示不限"`  // 配额，0 表示不限
	UsedBytes    int64       `json:"usedBytes"    orm:"used_bytes"    description:"已用容量"`       // 已用容量
	FileCount    int         `json:"fileCount"    orm:"file_count"    description:"文件数量"`       // 文件数量
	Status       int         `json:"status"       orm:"status"        description:"状态：1启用，0停用"` // 状态：1启用，0停用
	LastLoginAt  *gtime.Time `json:"lastLoginAt"  orm:"last_login_at" description:"上次登录时间"`     // 上次登录时间
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:"创建时间"`       // 创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:"更新时间"`       // 更新时间
}
