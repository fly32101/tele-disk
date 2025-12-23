package files

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "tele-disk/api/files/v1"
	"tele-disk/internal/consts"
	"tele-disk/internal/dao"
	"tele-disk/internal/model/do"
	"tele-disk/internal/model/entity"
	telegramSvc "tele-disk/internal/service/telegram"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func getScheme(r *ghttp.Request) string {
	scheme := r.URL.Scheme
	if scheme == "" {
		if r.TLS != nil {
			scheme = "https"
		} else if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
			scheme = proto
		} else {
			scheme = "http"
		}
	}
	return scheme
}

// Upload 上传文件
func (c *ControllerV1) Upload(ctx context.Context, req *v1.UploadReq) (res *v1.UploadRes, err error) {
	r := g.RequestFromCtx(ctx)
	userID := r.GetCtxVar("userId").Int64()
	if userID == 0 {
		return nil, gerror.New("unauthorized")
	}

	uploadFile := r.GetUploadFile("file")
	if uploadFile == nil {
		return nil, gerror.New("missing file")
	}
	if uploadFile.Size > consts.MaxUploadSizeBytes {
		return nil, gerror.Newf("file too large, max %d bytes", consts.MaxUploadSizeBytes)
	}

	fileReader, err := uploadFile.Open()
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	data, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, err
	}

	hashSha := sha256.Sum256(data)
	sha256Hex := hex.EncodeToString(hashSha[:])
	hashMd5 := md5.Sum(data)
	md5Hex := hex.EncodeToString(hashMd5[:])

	mimeType := uploadFile.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType(data)
	}

	// 去重：同用户同 MD5 直接返回已有记录
	var existing entity.Files
	err = dao.Files.Ctx(ctx).Where(g.Map{
		"user_id":    userID,
		"md5":        md5Hex,
		"is_deleted": consts.FlagNo,
	}).Scan(&existing)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if existing.Id > 0 {
		scheme := getScheme(r)
		proxyURL := fmt.Sprintf("%s://%s/api/files/proxy/%s", scheme, r.Host, existing.Md5)
		return &v1.UploadRes{
			Id:        existing.Id,
			Md5:       existing.Md5,
			FileName:  existing.FileName,
			SizeBytes: existing.SizeBytes,
			MimeType:  existing.MimeType,
			ProxyUrl:  proxyURL,
		}, nil
	}

	account, err := telegramSvc.GetDefaultAccount(ctx, userID)
	if err != nil {
		return nil, err
	}

	uploadResult, err := telegramSvc.SendDocument(ctx, account, uploadFile.Filename, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	fileName := uploadFile.Filename
	if uploadResult.FileName != "" {
		fileName = uploadResult.FileName
	}
	sizeBytes := uploadResult.FileSize
	if sizeBytes == 0 {
		sizeBytes = uploadFile.Size
	}

	var newID int64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		id, err := tx.Model(dao.Files.Table()).Data(g.Map{
			"user_id":                 userID,
			"telegram_account_id":     account.Id,
			"telegram_file_id":        uploadResult.FileID,
			"telegram_file_unique_id": uploadResult.FileUniqueID,
			"telegram_message_id":     uploadResult.MessageID,
			"telegram_chat_id":        uploadResult.ChatID,
			"file_name":               fileName,
			"size_bytes":              sizeBytes,
			"mime_type":               mimeType,
			"md5":                     md5Hex,
			"sha256":                  sha256Hex,
			"storage_provider":        "telegram",
			"object_key":              uploadResult.FileID,
			"status":                  consts.FileStatusOk,
			"is_deleted":              consts.FlagNo,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		newID = id

		_, err = tx.Model(dao.Users.Table()).WherePri(userID).Data(g.Map{
			"used_bytes": &gdb.Counter{Field: "used_bytes", Value: float64(sizeBytes)},
			"file_count": &gdb.Counter{Field: "file_count", Value: float64(1)},
		}).Update()
		return err
	})
	if err != nil {
		return nil, err
	}

	scheme := getScheme(r)
	proxyURL := fmt.Sprintf("%s://%s/api/files/proxy/%s", scheme, r.Host, md5Hex)

	res = &v1.UploadRes{
		Id:        newID,
		Md5:       md5Hex,
		FileName:  fileName,
		SizeBytes: sizeBytes,
		MimeType:  mimeType,
		ProxyUrl:  proxyURL,
	}
	return
}

// List 文件列表
func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	r := g.RequestFromCtx(ctx)
	userID := r.GetCtxVar("userId").Int64()
	if userID == 0 {
		return nil, gerror.New("unauthorized")
	}
	page := req.Page
	if page <= 0 {
		page = consts.DefaultPage
	}
	pageSize := req.PageSize
	if pageSize <= 0 || pageSize > consts.MaxPageSize {
		pageSize = consts.DefaultPageSize
	}

	var totalInt int
	totalInt, err = dao.Files.Ctx(ctx).
		Where(do.Files{UserId: userID, IsDeleted: consts.FlagNo}).
		Count()
	if err != nil {
		return nil, err
	}

	var items []entity.Files
	err = dao.Files.Ctx(ctx).
		Where(do.Files{UserId: userID, IsDeleted: consts.FlagNo}).
		OrderDesc(dao.Files.Columns().Id).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&items)
	if err != nil {
		return nil, err
	}

	scheme := getScheme(r)
	list := make([]v1.FileItem, 0, len(items))
	for _, item := range items {
		list = append(list, v1.FileItem{
			Id:        item.Id,
			Md5:       item.Md5,
			FileName:  item.FileName,
			SizeBytes: item.SizeBytes,
			MimeType:  item.MimeType,
			CreatedAt: item.CreatedAt.String(),
			ProxyUrl:  fmt.Sprintf("%s://%s/api/files/proxy/%s", scheme, r.Host, item.Md5),
		})
	}

	res = &v1.ListRes{Total: int64(totalInt), Page: page, Files: list}
	return
}

// Stats 统计
func (c *ControllerV1) Stats(ctx context.Context, req *v1.StatsReq) (res *v1.StatsRes, err error) {
	r := g.RequestFromCtx(ctx)
	userID := r.GetCtxVar("userId").Int64()
	if userID == 0 {
		return nil, gerror.New("unauthorized")
	}
	var user entity.Users
	if err := dao.Users.Ctx(ctx).WherePri(userID).Scan(&user); err != nil {
		return nil, err
	}
	if user.Id == 0 {
		return nil, gerror.New("user not found")
	}
	res = &v1.StatsRes{UsedBytes: user.UsedBytes, FileCount: user.FileCount, Quota: user.QuotaBytes}
	return
}

// Proxy 代理下载
func (c *ControllerV1) Proxy(ctx context.Context, req *v1.ProxyReq) (res *v1.ProxyRes, err error) {
	r := g.RequestFromCtx(ctx)
	userID := r.GetCtxVar("userId").Int64()
	isAdmin := r.GetCtxVar("isAdmin").Bool()

	if req.Md5 == "" {
		return nil, gerror.New("missing md5")
	}

	var file entity.Files
	if err := dao.Files.Ctx(ctx).
		Where(g.Map{"md5": req.Md5, "is_deleted": consts.FlagNo}).
		OrderDesc(dao.Files.Columns().Id).
		Limit(1).
		Scan(&file); err != nil {
		return nil, err
	}
	if file.Id == 0 || file.IsDeleted == consts.FlagYes {
		return nil, gerror.New("file not found")
	}
	if userID > 0 && !isAdmin && file.UserId != userID {
		return nil, gerror.New("forbidden")
	}

	accountID := file.TelegramAccountId
	var account *entity.TelegramAccounts
	if accountID > 0 {
		account, err = telegramSvc.GetAccountByID(ctx, accountID)
	} else {
		account, err = telegramSvc.GetDefaultAccount(ctx, userID)
	}
	if err != nil {
		return nil, err
	}

	fileURL, err := telegramSvc.GetFileURL(ctx, account.BotTokenEnc, file.TelegramFileId)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contentType := file.MimeType
	if contentType == "" {
		contentType = resp.Header.Get("Content-Type")
	}
	if contentType != "" {
		r.Response.Header().Set("Content-Type", contentType)
	}
	if cl := resp.Header.Get("Content-Length"); cl != "" {
		r.Response.Header().Set("Content-Length", cl)
	}
	r.Response.Header().Set("Cache-Control", "public, max-age=31536000")

	// Decide inline vs attachment
	inlinePrefixes := []string{
		"image/",
		"text/",
		"application/pdf",
		"video/",
		"audio/",
	}
	disposition := fmt.Sprintf("attachment; filename=\"%s\"", file.FileName)
	for _, p := range inlinePrefixes {
		if strings.HasPrefix(contentType, p) {
			disposition = fmt.Sprintf("inline; filename=\"%s\"", file.FileName)
			break
		}
	}
	r.Response.Header().Set("Content-Disposition", disposition)

	r.Response.WriteStatus(http.StatusOK)
	_, _ = io.Copy(r.Response.Writer, resp.Body)
	r.Exit()
	return nil, nil
}
