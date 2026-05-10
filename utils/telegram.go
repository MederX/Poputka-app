package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
)

func VerifyTelegramAuth(initData string, botToken string) bool {
	vals, err := url.ParseQuery(initData)
	if err != nil {
		return false
	}

	hash := vals.Get("hash")
	if hash == "" {
		return false
	}
	vals.Del("hash")

	keys := make([]string, 0, len(vals))
	for k := range vals {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+vals.Get(k))
	}
	dataCheckString := strings.Join(parts, "\n")

	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(botToken))
	sk := secretKey.Sum(nil)

	mac := hmac.New(sha256.New, sk)
	mac.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expectedHash), []byte(hash))
}
