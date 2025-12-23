package auth

import (
    "context"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "errors"
    "strings"
    "time"

    "github.com/gogf/gf/v2/frame/g"
)

// Claims 自定义简易 Token 负载
// 仅用于内部接口，生产建议替换为成熟 JWT 方案。
type Claims struct {
    UserID   int64  `json:"uid"`
    Username string `json:"u"`
    IsAdmin  bool   `json:"adm"`
    ExpireAt int64  `json:"exp"`
}

func tokenSecret() string {
    secret, _ := g.Cfg().Get(context.Background(), "auth.tokenSecret")
    if s := secret.String(); s != "" {
        return s
    }
    return "change-me"
}

func tokenTTL() time.Duration {
    ttl, _ := g.Cfg().Get(context.Background(), "auth.tokenTTLHours")
    if ttl != nil {
        if v := ttl.Int64(); v > 0 {
            return time.Duration(v) * time.Hour
        }
    }
    return 7 * 24 * time.Hour
}

// HashPassword 使用 sha256(secret + username + password)
func HashPassword(username, password string) string {
    h := sha256.Sum256([]byte(tokenSecret() + ":" + username + ":" + password))
    return hex.EncodeToString(h[:])
}

func VerifyPassword(username, password, storedHash string) bool {
    return HashPassword(username, password) == storedHash
}

// GenerateToken 创建签名 token
func GenerateToken(userID int64, username string, isAdmin bool) (string, error) {
    payload := Claims{
        UserID:   userID,
        Username: username,
        IsAdmin:  isAdmin,
        ExpireAt: time.Now().Add(tokenTTL()).Unix(),
    }
    data, err := json.Marshal(payload)
    if err != nil {
        return "", err
    }
    encoded := base64.RawURLEncoding.EncodeToString(data)
    mac := hmac.New(sha256.New, []byte(tokenSecret()))
    mac.Write([]byte(encoded))
    sig := mac.Sum(nil)
    token := encoded + "." + hex.EncodeToString(sig)
    return token, nil
}

// ParseToken 校验并返回 Claims
func ParseToken(token string) (*Claims, error) {
    parts := strings.Split(token, ".")
    if len(parts) != 2 {
        return nil, errors.New("invalid token format")
    }
    payloadB64, sigHex := parts[0], parts[1]
    payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadB64)
    if err != nil {
        return nil, errors.New("invalid token payload")
    }
    expectedMac := hmac.New(sha256.New, []byte(tokenSecret()))
    expectedMac.Write([]byte(payloadB64))
    expectedSig := expectedMac.Sum(nil)
    gotSig, err := hex.DecodeString(sigHex)
    if err != nil {
        return nil, errors.New("invalid token signature")
    }
    if !hmac.Equal(gotSig, expectedSig) {
        return nil, errors.New("signature mismatch")
    }
    var claims Claims
    if err := json.Unmarshal(payloadBytes, &claims); err != nil {
        return nil, errors.New("invalid token claims")
    }
    if claims.ExpireAt > 0 && time.Now().Unix() > claims.ExpireAt {
        return nil, errors.New("token expired")
    }
    return &claims, nil
}
