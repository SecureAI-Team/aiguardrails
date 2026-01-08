package auth

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SMSCode 验证码
type SMSCode struct {
	ID        string
	Phone     string
	Code      string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}

// SMSProvider 短信服务接口
type SMSProvider interface {
	Send(phone, code string) error
}

// SMSStore 短信验证码存储
type SMSStore struct {
	db       *sql.DB
	provider SMSProvider
}

// NewSMSStore 创建短信存储
func NewSMSStore(db *sql.DB, provider SMSProvider) *SMSStore {
	return &SMSStore{db: db, provider: provider}
}

// GenerateCode 生成6位验证码
func GenerateCode() string {
	bytes := make([]byte, 3)
	_, _ = rand.Read(bytes)
	code := fmt.Sprintf("%06d", int(bytes[0])<<16|int(bytes[1])<<8|int(bytes[2]))
	return code[:6]
}

// SendCode 发送验证码
func (s *SMSStore) SendCode(phone string) error {
	// 检查发送频率限制 (1分钟内只能发送1次)
	var count int
	s.db.QueryRow(`SELECT COUNT(*) FROM sms_codes WHERE phone=$1 AND created_at > NOW() - INTERVAL '1 minute'`, phone).Scan(&count)
	if count > 0 {
		return errors.New("发送太频繁，请稍后再试")
	}

	// 生成验证码
	code := GenerateCode()
	expiresAt := time.Now().Add(5 * time.Minute)
	id := uuid.NewString()

	// 存储验证码
	_, err := s.db.Exec(`INSERT INTO sms_codes (id, phone, code, expires_at) VALUES ($1, $2, $3, $4)`,
		id, phone, code, expiresAt)
	if err != nil {
		return err
	}

	// 发送短信
	if s.provider != nil {
		if err := s.provider.Send(phone, code); err != nil {
			return err
		}
	}

	return nil
}

// VerifyCode 验证验证码
func (s *SMSStore) VerifyCode(phone, code string) error {
	var id string
	err := s.db.QueryRow(`SELECT id FROM sms_codes WHERE phone=$1 AND code=$2 AND expires_at > NOW() AND used=false ORDER BY created_at DESC LIMIT 1`, phone, code).
		Scan(&id)
	if err != nil {
		return errors.New("验证码无效或已过期")
	}

	// 标记为已使用
	_, _ = s.db.Exec(`UPDATE sms_codes SET used=true WHERE id=$1`, id)
	return nil
}

// CleanExpired 清理过期验证码
func (s *SMSStore) CleanExpired() error {
	_, err := s.db.Exec(`DELETE FROM sms_codes WHERE expires_at < NOW() OR used=true`)
	return err
}

// AliyunSMSProvider 阿里云短信
type AliyunSMSProvider struct {
	AccessKeyID     string
	AccessKeySecret string
	SignName        string
	TemplateCode    string
}

// Send 发送短信 (阿里云)
func (p *AliyunSMSProvider) Send(phone, code string) error {
	// TODO: 实现阿里云短信发送
	// 这里是占位实现，实际需要接入阿里云短信SDK
	fmt.Printf("[SMS] Sending code %s to %s\n", code, phone)
	return nil
}

// TencentSMSProvider 腾讯云短信
type TencentSMSProvider struct {
	SecretID   string
	SecretKey  string
	AppID      string
	SignName   string
	TemplateID string
}

// Send 发送短信 (腾讯云)
func (p *TencentSMSProvider) Send(phone, code string) error {
	// TODO: 实现腾讯云短信发送
	// 这里是占位实现，实际需要接入腾讯云短信SDK
	fmt.Printf("[SMS] Sending code %s to %s\n", code, phone)
	return nil
}
