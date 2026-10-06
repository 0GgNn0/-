package models

import (
	"crypto/rand"
	"math/big"
	"os"
	"time"

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

// RotateAdminPasswordFromEnv 仅在显式请求时（--rotate-admin-password）把环境变量里的密码
// 重新哈希写库，用于轮换已泄露的运维凭据。未设置环境变量时不做任何改动，避免误锁死后台。
func RotateAdminPasswordFromEnv(envPassword string) (string, error) {
	if envPassword == "" {
		return "skipped: ADMIN_PASSWORD 未设置", nil
	}
	if VerifyPassword(envPassword) {
		return "skipped: 与当前密码一致", nil
	}
	if err := ChangePassword(envPassword); err != nil {
		return "", err
	}
	return "rotated: 已按 ADMIN_PASSWORD 重置后台密码", nil
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
	if err := SetSetting("admin_password_hash", string(hash)); err != nil {
		return err
	}
	// 审计：记录最后一次改密时间与来源，避免"密码被谁改过"无从查证
	source := "admin_panel"
	if len(os.Args) > 1 {
		for _, a := range os.Args[1:] {
			if a == "--rotate-admin-password" {
				source = "env_rotation"
			}
		}
	}
	SetSetting("admin_password_changed_at", time.Now().Format("2006-01-02 15:04:05"))
	SetSetting("admin_password_changed_by", source)
	return nil
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
