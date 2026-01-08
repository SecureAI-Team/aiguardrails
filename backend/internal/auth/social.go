package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// SocialProvider 社交登录提供商
type SocialProvider string

const (
	ProviderWeChat SocialProvider = "wechat"
	ProviderAlipay SocialProvider = "alipay"
	ProviderPhone  SocialProvider = "phone"
	ProviderWecom  SocialProvider = "wecom"
)

// SocialAccount 社交账号
type SocialAccount struct {
	ID         string          `json:"id"`
	UserID     string          `json:"user_id"`
	Provider   SocialProvider  `json:"provider"`
	ProviderID string          `json:"provider_id"`
	Profile    json.RawMessage `json:"profile"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// OAuthConfig OAuth配置
type OAuthConfig struct {
	WeChatAppID      string
	WeChatAppSecret  string
	AlipayAppID      string
	AlipayPrivateKey string
	AlipayPublicKey  string
	RedirectBaseURL  string
}

// SocialAuthStore 社交登录存储
type SocialAuthStore struct {
	db     *sql.DB
	config OAuthConfig
}

// NewSocialAuthStore 创建社交登录存储
func NewSocialAuthStore(db *sql.DB, config OAuthConfig) *SocialAuthStore {
	return &SocialAuthStore{db: db, config: config}
}

// GenerateState 生成OAuth state
func (s *SocialAuthStore) GenerateState(provider SocialProvider, redirectURL string) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	state := hex.EncodeToString(bytes)

	expiresAt := time.Now().Add(10 * time.Minute)
	_, err := s.db.Exec(`INSERT INTO oauth_states (state, provider, redirect_url, expires_at) VALUES ($1, $2, $3, $4)`,
		state, provider, redirectURL, expiresAt)
	if err != nil {
		return "", err
	}
	return state, nil
}

// ValidateState 验证OAuth state
func (s *SocialAuthStore) ValidateState(state string) (*SocialProvider, string, error) {
	var provider SocialProvider
	var redirectURL sql.NullString
	err := s.db.QueryRow(`SELECT provider, redirect_url FROM oauth_states WHERE state=$1 AND expires_at > NOW()`, state).
		Scan(&provider, &redirectURL)
	if err != nil {
		return nil, "", errors.New("invalid or expired state")
	}
	// Delete used state
	_, _ = s.db.Exec(`DELETE FROM oauth_states WHERE state=$1`, state)
	return &provider, redirectURL.String, nil
}

// GetWeChatAuthURL 获取微信授权URL
func (s *SocialAuthStore) GetWeChatAuthURL(redirectURI, state string) string {
	params := url.Values{}
	params.Set("appid", s.config.WeChatAppID)
	params.Set("redirect_uri", redirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "snsapi_userinfo")
	params.Set("state", state)
	return "https://open.weixin.qq.com/connect/oauth2/authorize?" + params.Encode() + "#wechat_redirect"
}

// GetAlipayAuthURL 获取支付宝授权URL
func (s *SocialAuthStore) GetAlipayAuthURL(redirectURI, state string) string {
	params := url.Values{}
	params.Set("app_id", s.config.AlipayAppID)
	params.Set("redirect_uri", redirectURI)
	params.Set("scope", "auth_user")
	params.Set("state", state)
	return "https://openauth.alipay.com/oauth2/publicAppAuthorize.htm?" + params.Encode()
}

// ExchangeWeChatCode 微信code换取用户信息
func (s *SocialAuthStore) ExchangeWeChatCode(code string) (*WeChatUser, error) {
	// Step 1: Get access token
	tokenURL := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		s.config.WeChatAppID, s.config.WeChatAppSecret, code)

	resp, err := http.Get(tokenURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		OpenID      string `json:"openid"`
		UnionID     string `json:"unionid"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}
	if tokenResp.ErrCode != 0 {
		return nil, fmt.Errorf("wechat error: %s", tokenResp.ErrMsg)
	}

	// Step 2: Get user info
	userURL := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		tokenResp.AccessToken, tokenResp.OpenID)

	resp2, err := http.Get(userURL)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()

	body2, _ := io.ReadAll(resp2.Body)
	var user WeChatUser
	if err := json.Unmarshal(body2, &user); err != nil {
		return nil, err
	}
	user.UnionID = tokenResp.UnionID
	return &user, nil
}

// WeChatUser 微信用户信息
type WeChatUser struct {
	OpenID     string `json:"openid"`
	UnionID    string `json:"unionid"`
	Nickname   string `json:"nickname"`
	HeadImgURL string `json:"headimgurl"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
}

// BindSocialAccount 绑定社交账号
func (s *SocialAuthStore) BindSocialAccount(userID string, provider SocialProvider, providerID string, profile interface{}) (*SocialAccount, error) {
	profileJSON, _ := json.Marshal(profile)
	id := uuid.NewString()
	now := time.Now().UTC()

	_, err := s.db.Exec(`INSERT INTO user_social_accounts (id, user_id, provider, provider_id, profile, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (provider, provider_id) DO UPDATE SET user_id=$2, profile=$5, updated_at=$7`,
		id, userID, provider, providerID, profileJSON, now, now)
	if err != nil {
		return nil, err
	}
	return &SocialAccount{ID: id, UserID: userID, Provider: provider, ProviderID: providerID, Profile: profileJSON, CreatedAt: now}, nil
}

// FindByProviderID 根据社交账号查找用户
func (s *SocialAuthStore) FindByProviderID(provider SocialProvider, providerID string) (*SocialAccount, error) {
	var sa SocialAccount
	err := s.db.QueryRow(`SELECT id, user_id, provider, provider_id, profile, created_at FROM user_social_accounts WHERE provider=$1 AND provider_id=$2`,
		provider, providerID).Scan(&sa.ID, &sa.UserID, &sa.Provider, &sa.ProviderID, &sa.Profile, &sa.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &sa, nil
}

// ListByUserID 列出用户绑定的社交账号
func (s *SocialAuthStore) ListByUserID(userID string) ([]SocialAccount, error) {
	rows, err := s.db.Query(`SELECT id, user_id, provider, provider_id, profile, created_at FROM user_social_accounts WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []SocialAccount
	for rows.Next() {
		var sa SocialAccount
		if err := rows.Scan(&sa.ID, &sa.UserID, &sa.Provider, &sa.ProviderID, &sa.Profile, &sa.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, sa)
	}
	return accounts, nil
}
