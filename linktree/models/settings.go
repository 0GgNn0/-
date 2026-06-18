package models

import (
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

func GetSetting(key string) (string, error) {
	var value string
	err := DB.QueryRow("SELECT value FROM settings WHERE key=?", key).Scan(&value)
	return value, err
}

func GetSettingOrDefault(key, defaultVal string) string {
	val, err := GetSetting(key)
	if err != nil {
		return defaultVal
	}
	return val
}

func SetSetting(key, value string) error {
	_, err := DB.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)", key, value)
	return err
}

func InitAdminPassword(envPassword string) (string, error) {
	_, err := GetSetting("admin_password_hash")
	if err == nil {
		return "", nil // already initialized
	}

	password := envPassword
	if password == "" {
		password = generateRandomPassword(12)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	if err := SetSetting("admin_password_hash", string(hash)); err != nil {
		return "", err
	}

	return password, nil
}

func VerifyPassword(password string) bool {
	hash, err := GetSetting("admin_password_hash")
	if err != nil {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func ChangePassword(newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return SetSetting("admin_password_hash", string(hash))
}

func generateRandomPassword(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		result[i] = chars[n.Int64()]
	}
	return string(result)
}

func MaskAPIKey(key string) string {
	if key == "" {
		return ""
	}
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

func GetAIBaseURL() string {
	return GetSettingOrDefault("ai_base_url", "https://api.deepseek.com")
}

func GetAIAPIKey() string {
	val, _ := GetSetting("ai_api_key")
	return val
}

func GetAIModel() string {
	return GetSettingOrDefault("ai_model", "deepseek-v4-flash")
}
