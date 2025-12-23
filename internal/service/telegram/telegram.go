package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gogf/gf/v2/errors/gerror"

	"tele-disk/internal/consts"
	"tele-disk/internal/dao"
	"tele-disk/internal/logic/secret"
	"tele-disk/internal/model/do"
	"tele-disk/internal/model/entity"
)

const (
	apiBaseURL  = "https://api.telegram.org/bot%s"
	fileBaseURL = "https://api.telegram.org/file/bot%s/%s"
)

type apiResponse struct {
	Ok          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	ErrorCode   int             `json:"error_code,omitempty"`
	Description string          `json:"description,omitempty"`
}

type telegramMessage struct {
	MessageID int          `json:"message_id"`
	Chat      telegramChat `json:"chat"`
	Document  *telegramDoc `json:"document,omitempty"`
}

type telegramChat struct {
	ID int64 `json:"id"`
}

type telegramDoc struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileName     string `json:"file_name,omitempty"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int64  `json:"file_size,omitempty"`
}

type telegramFile struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
}

type UploadResult struct {
	FileID       string
	FileUniqueID string
	FileName     string
	MimeType     string
	FileSize     int64
	MessageID    int64
	ChatID       int64
}

// GetDefaultAccount 获取用户默认账号，先找该用户 is_default=1，否则取该用户的第一条启用账号。
func GetDefaultAccount(ctx context.Context, userID int64) (*entity.TelegramAccounts, error) {
	var account entity.TelegramAccounts
	q := dao.TelegramAccounts.Ctx(ctx)
	if userID > 0 {
		q = q.Where(do.TelegramAccounts{UserId: userID})
	}
	err := q.Where(do.TelegramAccounts{Status: consts.StatusEnabled, IsDefault: consts.FlagYes}).
		OrderAsc(dao.TelegramAccounts.Columns().Id).
		Limit(1).
		Scan(&account)
	if err != nil {
		return nil, err
	}
	if account.Id == 0 {
		q2 := dao.TelegramAccounts.Ctx(ctx)
		if userID > 0 {
			q2 = q2.Where(do.TelegramAccounts{UserId: userID})
		}
		err = q2.Where(do.TelegramAccounts{Status: consts.StatusEnabled}).
			OrderAsc(dao.TelegramAccounts.Columns().Id).
			Limit(1).
			Scan(&account)
		if err != nil {
			return nil, err
		}
	}
	if account.Id == 0 {
		return nil, gerror.New("no telegram account available")
	}
	return &account, nil
}

// GetAccountByID 获取账号
func GetAccountByID(ctx context.Context, id int64) (*entity.TelegramAccounts, error) {
	var account entity.TelegramAccounts
	err := dao.TelegramAccounts.Ctx(ctx).WherePri(id).Scan(&account)
	if err != nil {
		return nil, err
	}
	if account.Id == 0 {
		return nil, gerror.New("telegram account not found")
	}
	return &account, nil
}

// SendDocument 上传文件到 Telegram
func SendDocument(ctx context.Context, account *entity.TelegramAccounts, filename string, reader io.Reader) (*UploadResult, error) {
	if account == nil || account.BotTokenEnc == "" || account.ChatId == 0 {
		return nil, gerror.New("invalid telegram account config")
	}

	botToken, err := secret.DecryptString(account.BotTokenEnc)
	if err != nil {
		return nil, gerror.Wrap(err, "decrypt bot token failed")
	}

	apiURL := fmt.Sprintf(apiBaseURL+"/sendDocument", botToken)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("chat_id", fmt.Sprintf("%d", account.ChatId)); err != nil {
		return nil, err
	}
	part, err := writer.CreateFormFile("document", filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, reader); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	if !apiResp.Ok {
		return nil, gerror.Newf("telegram api error: %s", apiResp.Description)
	}
	var msg telegramMessage
	if err := json.Unmarshal(apiResp.Result, &msg); err != nil {
		return nil, err
	}
	if msg.Document == nil || msg.Document.FileID == "" {
		return nil, gerror.New("telegram response missing document")
	}
	return &UploadResult{
		FileID:       msg.Document.FileID,
		FileUniqueID: msg.Document.FileUniqueID,
		FileName:     msg.Document.FileName,
		MimeType:     msg.Document.MimeType,
		FileSize:     msg.Document.FileSize,
		MessageID:    int64(msg.MessageID),
		ChatID:       msg.Chat.ID,
	}, nil
}

// GetFileURL 获取下载地址
func GetFileURL(ctx context.Context, botTokenEnc, fileID string) (string, error) {
	if botTokenEnc == "" || fileID == "" {
		return "", gerror.New("missing bot token or file id")
	}
	botToken, err := secret.DecryptString(botTokenEnc)
	if err != nil {
		return "", gerror.Wrap(err, "decrypt bot token failed")
	}
	apiURL := fmt.Sprintf(apiBaseURL+"/getFile", botToken)
	apiURL = fmt.Sprintf("%s?file_id=%s", apiURL, fileID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", err
	}
	if !apiResp.Ok {
		return "", gerror.Newf("telegram api error: %s", apiResp.Description)
	}
	var file telegramFile
	if err := json.Unmarshal(apiResp.Result, &file); err != nil {
		return "", err
	}
	if file.FilePath == "" {
		return "", gerror.New("telegram file path empty")
	}
	return fmt.Sprintf(fileBaseURL, botToken, file.FilePath), nil
}
