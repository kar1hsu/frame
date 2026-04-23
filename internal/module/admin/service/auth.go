package service

import (
	"errors"
	"fmt"
	"time"

	"frame/internal/app"
	"frame/internal/dao"
	"frame/internal/pkg/cache"
	jwtpkg "frame/internal/pkg/jwt"
	"frame/internal/pkg/utils"
)

type AuthService struct {
	userDAO *dao.UserDAO
}

func NewAuthService() *AuthService {
	return &AuthService{userDAO: dao.NewUserDAO()}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	if cache.IsLoginLocked(req.Username) {
		ttl := cache.GetLoginLockTTL(req.Username)
		minutes := int(ttl.Minutes()) + 1
		return nil, fmt.Errorf("登录失败次数过多，请 %d 分钟后重试", minutes)
	}

	user, err := s.userDAO.GetByUsername(req.Username)
	if err != nil {
		cache.IncrLoginFail(req.Username)
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		count, _ := cache.IncrLoginFail(req.Username)
		remaining := int64(5) - count
		if remaining > 0 {
			return nil, fmt.Errorf("用户名或密码错误，还可尝试 %d 次", remaining)
		}
		return nil, errors.New("登录失败次数过多，账户已被临时锁定")
	}

	cache.ClearLoginFail(req.Username)

	roleCode := "default"
	if len(user.Roles) > 0 {
		roleCode = user.Roles[0].Code
	}

	token, err := jwtpkg.GenerateToken(user.ID, user.Username, roleCode)
	if err != nil {
		return nil, errors.New("生成 Token 失败")
	}

	return &LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
	}, nil
}

func (s *AuthService) Logout(token string) error {
	expiration := time.Duration(app.Cfg.JWT.Expire) * time.Second
	return cache.BlacklistToken(token, expiration)
}
