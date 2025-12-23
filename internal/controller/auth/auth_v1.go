package auth

import (
    "context"
    "errors"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	v1 "tele-disk/api/auth/v1"
	"tele-disk/internal/consts"
	"tele-disk/internal/dao"
	"tele-disk/internal/logic/auth"
	"tele-disk/internal/model/do"
	"tele-disk/internal/model/entity"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
    return &ControllerV1{}
}

// Register 用户注册
func (c *ControllerV1) Register(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err error) {
    // 查重
    count, err := dao.Users.Ctx(ctx).Where(do.Users{Username: req.Username}).Count()
    if err != nil {
        return nil, err
    }
    if count > 0 {
        return nil, gerror.New("用户名已存在")
    }

	hash := auth.HashPassword(req.Username, req.Password)
	data := do.Users{
		Username:     req.Username,
		PasswordHash: hash,
		Status:       consts.StatusEnabled,
	}
    uid, err := dao.Users.Ctx(ctx).Data(data).InsertAndGetId()
    if err != nil {
        return nil, err
    }
    token, err := auth.GenerateToken(uid, req.Username, false)
    if err != nil {
        return nil, err
    }
    res = &v1.RegisterRes{Token: token, UserID: uid, Username: req.Username, IsAdmin: false}
    return
}

// Login 用户登录
func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
    var user entity.Users
    err = dao.Users.Ctx(ctx).Where(do.Users{Username: req.Username}).Scan(&user)
    if err != nil {
        return nil, err
    }
    if user.Id == 0 {
        return nil, gerror.New("用户名或密码错误")
    }
	if user.Status == consts.StatusDisabled {
		return nil, errors.New("用户已停用")
	}
    if !auth.VerifyPassword(req.Username, req.Password, user.PasswordHash) {
        return nil, gerror.New("用户名或密码错误")
    }
    // 更新最后登录时间（忽略错误）
    _, _ = dao.Users.Ctx(ctx).Data(g.Map{"last_login_at": gtime.Now()}).WherePri(user.Id).Update()

	token, err := auth.GenerateToken(user.Id, user.Username, user.IsAdmin == consts.FlagYes)
	if err != nil {
		return nil, err
	}
	res = &v1.LoginRes{Token: token, UserID: user.Id, Username: user.Username, IsAdmin: user.IsAdmin == consts.FlagYes}
	return
}
