package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"aiguardrails/internal/auth"
)

// registerSocialAuthRoutes 注册社交登录路由
func (s *Server) registerSocialAuthRoutes(r chi.Router) {
	r.Get("/auth/oauth/wechat/url", s.getWeChatAuthURL)
	r.Get("/auth/oauth/alipay/url", s.getAlipayAuthURL)
	r.Post("/auth/login/wechat", s.loginWeChat)
	r.Post("/auth/login/alipay", s.loginAlipay)
	r.Post("/auth/sms/send", s.sendSMSCode)
	r.Post("/auth/sms/verify", s.verifySMSLogin)
}

type oauthURLRequest struct {
	RedirectURI string `json:"redirect_uri"`
}

type oauthLoginRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type smsRequest struct {
	Phone string `json:"phone"`
}

type smsVerifyRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

func (s *Server) getWeChatAuthURL(w http.ResponseWriter, r *http.Request) {
	redirectURI := r.URL.Query().Get("redirect_uri")
	if redirectURI == "" {
		redirectURI = s.cfg.SocialAuthCallbackURL + "/wechat"
	}

	state, err := s.socialAuth.GenerateState(auth.ProviderWeChat, redirectURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := s.socialAuth.GetWeChatAuthURL(redirectURI, state)
	s.writeJSON(w, http.StatusOK, map[string]string{"url": url, "state": state})
}

func (s *Server) getAlipayAuthURL(w http.ResponseWriter, r *http.Request) {
	redirectURI := r.URL.Query().Get("redirect_uri")
	if redirectURI == "" {
		redirectURI = s.cfg.SocialAuthCallbackURL + "/alipay"
	}

	state, err := s.socialAuth.GenerateState(auth.ProviderAlipay, redirectURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := s.socialAuth.GetAlipayAuthURL(redirectURI, state)
	s.writeJSON(w, http.StatusOK, map[string]string{"url": url, "state": state})
}

func (s *Server) loginWeChat(w http.ResponseWriter, r *http.Request) {
	var req oauthLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 验证state
	_, _, err := s.socialAuth.ValidateState(req.State)
	if err != nil {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	// 换取用户信息
	wxUser, err := s.socialAuth.ExchangeWeChatCode(req.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 查找或创建用户
	providerID := wxUser.OpenID
	if wxUser.UnionID != "" {
		providerID = wxUser.UnionID
	}

	socialAccount, err := s.socialAuth.FindByProviderID(auth.ProviderWeChat, providerID)
	var user *auth.User

	if err != nil {
		// 新用户，自动注册
		user, err = s.userStore.Create("wx_"+providerID[:8], providerID, "tenant_user")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 绑定社交账号
		_, _ = s.socialAuth.BindSocialAccount(user.ID, auth.ProviderWeChat, providerID, wxUser)
	} else {
		user, _ = s.userStore.GetByID(socialAccount.UserID)
	}

	// 生成token
	token, err := s.jwtSigner.Sign(user.Username, user.Role, 24*time.Hour)
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	s.audit.RecordStore(s.auditStore, "wechat_login", map[string]string{
		"user_id": user.ID,
		"openid":  wxUser.OpenID,
	})

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  map[string]string{"id": user.ID, "username": user.Username, "role": user.Role},
	})
}

func (s *Server) loginAlipay(w http.ResponseWriter, r *http.Request) {
	// TODO: 实现支付宝登录逻辑 (类似微信)
	http.Error(w, "alipay login not implemented yet", http.StatusNotImplemented)
}

func (s *Server) sendSMSCode(w http.ResponseWriter, r *http.Request) {
	var req smsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Phone) != 11 {
		http.Error(w, "invalid phone number", http.StatusBadRequest)
		return
	}

	if err := s.smsStore.SendCode(req.Phone); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]string{"status": "sent"})
}

func (s *Server) verifySMSLogin(w http.ResponseWriter, r *http.Request) {
	var req smsVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.smsStore.VerifyCode(req.Phone, req.Code); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// 查找或创建用户
	socialAccount, err := s.socialAuth.FindByProviderID(auth.ProviderPhone, req.Phone)
	var user *auth.User

	if err != nil {
		// 新用户
		user, err = s.userStore.Create("phone_"+req.Phone[7:], req.Phone, "tenant_user")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = s.socialAuth.BindSocialAccount(user.ID, auth.ProviderPhone, req.Phone, nil)
	} else {
		user, _ = s.userStore.GetByID(socialAccount.UserID)
	}

	// 生成token
	token, err := s.jwtSigner.Sign(user.Username, user.Role, 24*time.Hour)
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	s.audit.RecordStore(s.auditStore, "phone_login", map[string]string{
		"user_id": user.ID,
		"phone":   req.Phone[:3] + "****" + req.Phone[7:],
	})

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  map[string]string{"id": user.ID, "username": user.Username, "role": user.Role},
	})
}
