package v1

import "github.com/gogf/gf/v2/frame/g"

// RegisterReq 用户注册
// Password 需要在应用层加密存储
// swagger:route POST /api/auth/register Auth register
//
// Responses:
//   default: RegisterRes
//
type RegisterReq struct {
    g.Meta   `path:"/api/auth/register" tags:"Auth" method:"post" summary:"用户注册"`
    Username string `json:"username" v:"required|length:3,64"`
    Password string `json:"password" v:"required|length:6,64"`
}

type RegisterRes struct {
    Token    string `json:"token"`
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    IsAdmin  bool   `json:"is_admin"`
}

// LoginReq 用户登录
// swagger:route POST /api/auth/login Auth login
//
// Responses:
//   default: LoginRes
//
type LoginReq struct {
    g.Meta   `path:"/api/auth/login" tags:"Auth" method:"post" summary:"用户登录"`
    Username string `json:"username" v:"required|length:3,64"`
    Password string `json:"password" v:"required|length:6,64"`
}

type LoginRes struct {
    Token    string `json:"token"`
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    IsAdmin  bool   `json:"is_admin"`
}
