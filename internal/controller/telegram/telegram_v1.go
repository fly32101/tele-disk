package telegram

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "tele-disk/api/telegram/v1"
	"tele-disk/internal/consts"
	"tele-disk/internal/dao"
	"tele-disk/internal/logic/secret"
	"tele-disk/internal/model/entity"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func currentUserID(ctx context.Context) int64 {
	r := g.RequestFromCtx(ctx)
	return r.GetCtxVar("userId").Int64()
}

// CreateAccount 创建账号
func (c *ControllerV1) CreateAccount(ctx context.Context, req *v1.CreateAccountReq) (res *v1.AccountRes, err error) {
	userID := currentUserID(ctx)
	if userID == 0 {
		return nil, gerror.New("unauthorized")
	}

	isDefault := consts.FlagNo
	if req.IsDefault {
		isDefault = consts.FlagYes
	}
	status := req.Status
	if status != consts.StatusDisabled {
		status = consts.StatusEnabled
	}

	var newID int64
	encToken, err := secret.EncryptString(req.BotToken)
	if err != nil {
		return nil, err
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if isDefault == consts.FlagYes {
			if _, err := tx.Model(dao.TelegramAccounts.Table()).Data(g.Map{"is_default": consts.FlagNo}).Where(g.Map{
				"is_default": consts.FlagYes,
				"user_id":    userID,
			}).Update(); err != nil {
				return err
			}
		}
		id, err := tx.Model(dao.TelegramAccounts.Table()).Data(g.Map{
			"user_id":       userID,
			"name":          req.Name,
			"bot_username":  req.BotUsername,
			"bot_token_enc": encToken,
			"bot_id":        req.BotId,
			"chat_id":       req.ChatId,
			"status":        status,
			"is_default":    isDefault,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		newID = id
		return nil
	})
	if err != nil {
		return nil, err
	}

	var account entity.TelegramAccounts
	if err := dao.TelegramAccounts.Ctx(ctx).WherePri(newID).Scan(&account); err != nil {
		return nil, err
	}
	return toAccountRes(&account), nil
}

// UpdateAccount 更新账号
func (c *ControllerV1) UpdateAccount(ctx context.Context, req *v1.UpdateAccountReq) (res *v1.AccountRes, err error) {
	userID := currentUserID(ctx)
	if userID == 0 {
		return nil, gerror.New("unauthorized")
	}

	data := g.Map{}
	if req.Name != "" {
		data["name"] = req.Name
	}
	if req.BotUsername != "" {
		data["bot_username"] = req.BotUsername
	}
	if req.BotToken != "" {
		enc, err := secret.EncryptString(req.BotToken)
		if err != nil {
			return nil, err
		}
		data["bot_token_enc"] = enc
	}
	if req.BotId != 0 {
		data["bot_id"] = req.BotId
	}
	if req.ChatId != 0 {
		data["chat_id"] = req.ChatId
	}
	if req.Status != nil {
		data["status"] = *req.Status
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if req.IsDefault != nil && *req.IsDefault {
			if _, err := tx.Model(dao.TelegramAccounts.Table()).Data(g.Map{"is_default": consts.FlagNo}).Where(g.Map{
				"is_default": consts.FlagYes,
				"user_id":    userID,
			}).Update(); err != nil {
				return err
			}
			data["is_default"] = consts.FlagYes
		} else if req.IsDefault != nil {
			data["is_default"] = consts.FlagNo
		}
		if len(data) == 0 {
			return nil
		}
		_, err := tx.Model(dao.TelegramAccounts.Table()).Data(data).WherePri(req.Id).Where("user_id", userID).Update()
		return err
	})
	if err != nil {
		return nil, err
	}

	var account entity.TelegramAccounts
	if err := dao.TelegramAccounts.Ctx(ctx).WherePri(req.Id).Scan(&account); err != nil {
		return nil, err
	}
	return toAccountRes(&account), nil
}

// SetDefault 设置默认账号
func (c *ControllerV1) SetDefault(ctx context.Context, req *v1.SetDefaultReq) (res *v1.AccountRes, err error) {
	userID := currentUserID(ctx)
	if userID == 0 {
		return nil, gerror.New("unauthorized")
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.TelegramAccounts.Table()).Data(g.Map{"is_default": consts.FlagNo}).Where(g.Map{
			"is_default": consts.FlagYes,
			"user_id":    userID,
		}).Update(); err != nil {
			return err
		}
		_, err := tx.Model(dao.TelegramAccounts.Table()).Data(g.Map{"is_default": consts.FlagYes}).WherePri(req.Id).Where("user_id", userID).Update()
		return err
	})
	if err != nil {
		return nil, err
	}

	var account entity.TelegramAccounts
	if err := dao.TelegramAccounts.Ctx(ctx).WherePri(req.Id).Scan(&account); err != nil {
		return nil, err
	}
	return toAccountRes(&account), nil
}

// ListAccounts 列表
func (c *ControllerV1) ListAccounts(ctx context.Context, req *v1.ListAccountsReq) (res *v1.ListAccountsRes, err error) {
	userID := currentUserID(ctx)
	if userID == 0 {
		return nil, gerror.New("unauthorized")
	}
	var accounts []entity.TelegramAccounts
	if err := dao.TelegramAccounts.Ctx(ctx).Where("user_id", userID).OrderAsc(dao.TelegramAccounts.Columns().Id).Scan(&accounts); err != nil {
		return nil, err
	}
	out := make([]v1.AccountRes, 0, len(accounts))
	for _, a := range accounts {
		out = append(out, *toAccountRes(&a))
	}
	return &v1.ListAccountsRes{Accounts: out}, nil
}

func toAccountRes(a *entity.TelegramAccounts) *v1.AccountRes {
	if a == nil {
		return &v1.AccountRes{}
	}
	return &v1.AccountRes{
		// 不返回解密后的 token
		Id:          a.Id,
		Name:        a.Name,
		BotUsername: a.BotUsername,
		BotToken:    "",
		BotId:       a.BotId,
		ChatId:      a.ChatId,
		Status:      a.Status,
		IsDefault:   a.IsDefault == consts.FlagYes,
		LastError:   a.LastError,
		CreatedAt:   a.CreatedAt.String(),
		UpdatedAt:   a.UpdatedAt.String(),
	}
}
