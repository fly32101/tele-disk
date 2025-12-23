package v1

import "github.com/gogf/gf/v2/frame/g"

// CreateAccountReq 创建 Telegram 账号
// swagger:route POST /api/telegram/accounts Telegram createAccount
//
// Responses:
//   default: AccountRes
//
type CreateAccountReq struct {
    g.Meta      `path:"/api/telegram/accounts" tags:"Telegram" method:"post" summary:"创建 Telegram 账号"`
    Name        string `json:"name" v:"required|length:1,64"`
    BotUsername string `json:"bot_username" v:"required|length:1,64"`
    BotToken    string `json:"bot_token" v:"required|length:10,128"`
    BotId       int64  `json:"bot_id"`
    ChatId      int64  `json:"chat_id"`
    IsDefault   bool   `json:"is_default"`
    Status      int    `json:"status" d:"1"`
}

// UpdateAccountReq 更新 Telegram 账号
// swagger:route PUT /api/telegram/accounts/{id} Telegram updateAccount
//
// Responses:
//   default: AccountRes
//
type UpdateAccountReq struct {
    g.Meta      `path:"/api/telegram/accounts/{id}" tags:"Telegram" method:"put" summary:"更新 Telegram 账号"`
    Id          int64  `json:"id" in:"path"`
    Name        string `json:"name"`
    BotUsername string `json:"bot_username"`
    BotToken    string `json:"bot_token"`
    BotId       int64  `json:"bot_id"`
    ChatId      int64  `json:"chat_id"`
    IsDefault   *bool  `json:"is_default"`
    Status      *int   `json:"status"`
}

// SetDefaultReq 设为默认账号
// swagger:route POST /api/telegram/accounts/{id}/default Telegram setDefault
//
// Responses:
//   default: AccountRes
//
type SetDefaultReq struct {
    g.Meta `path:"/api/telegram/accounts/{id}/default" tags:"Telegram" method:"post" summary:"设置默认 Telegram 账号"`
    Id     int64 `json:"id" in:"path"`
}

// ListAccountsReq 列表
// swagger:route GET /api/telegram/accounts Telegram listAccounts
//
// Responses:
//   default: ListAccountsRes
//
type ListAccountsReq struct {
    g.Meta `path:"/api/telegram/accounts" tags:"Telegram" method:"get" summary:"Telegram 账号列表"`
}

type AccountRes struct {
    Id          int64  `json:"id"`
    Name        string `json:"name"`
    BotUsername string `json:"bot_username"`
    BotToken    string `json:"bot_token"`
    BotId       int64  `json:"bot_id"`
    ChatId      int64  `json:"chat_id"`
    Status      int    `json:"status"`
    IsDefault   bool   `json:"is_default"`
    LastError   string `json:"last_error"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}

type ListAccountsRes struct {
    Accounts []AccountRes `json:"accounts"`
}
