package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

// getKey returns 32-byte key from config security.encryptKey (pad/truncate).
func getKey() []byte {
	def := "change-me-32bytes-change-me-32byte"
	cfg, _ := g.Cfg().Get(nil, "security.encryptKey")
	key := cfg.String()
	if key == "" {
		key = def
	}
	// normalize length to 32
	keyBytes := []byte(key)
	if len(keyBytes) >= 32 {
		return keyBytes[:32]
	}
	buf := make([]byte, 32)
	copy(buf, keyBytes)
	for i := len(keyBytes); i < 32; i++ {
		buf[i] = '0'
	}
	return buf
}

// EncryptString AES-GCM encrypts plaintext and returns base64(nonce|ciphertext).
func EncryptString(plain string) (string, error) {
	key := getKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nil, nonce, []byte(plain), nil)
	full := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(full), nil
}

// DecryptString decrypts base64(nonce|ciphertext) to plaintext.
func DecryptString(enc string) (string, error) {
	if strings.TrimSpace(enc) == "" {
		return "", nil
	}
	key := getKey()
	raw, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(raw) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := raw[:nonceSize], raw[nonceSize:]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
