package user

import (
	e "pigpq/internal/pkg/errors"
	"pigpq/internal/service"
	"time"

	c "pigpq/config"
	"pigpq/internal/model"
	j "pigpq/internal/pkg/untils/jwt"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   int64  `json:"expires_at"`
}
type UserService struct {
	service.Base
}

func NewUserService() *UserService {
	return &UserService{}
}

func (us UserService) Login(phone string) (*TokenResponse, error) {
	// 查询用户信息
	user := model.NewUser().GetUserByPhone(phone)

	if user == nil {
		err := e.NewBusinessError(e.UserDoesNotExist)
		return nil, err
	}

	// 判断用户状态是否禁用
	if user.Status != 1 {
		err := e.NewBusinessError(e.UserDoesNotExist)
		return nil, err
	}

	// 生成用户token
	claims := us.newCustomClaims(user)

	accessToken, err := j.Generate(claims)
	if err != nil {
		return nil, e.NewBusinessError(e.FAILURE, "生成Token失败")
	}

	return &TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresAt:   claims.ExpiresAt.Unix(),
	}, nil
}

// Refresh 刷新Token
func (us *UserService) Refresh(id uint) (*TokenResponse, error) {
	// 查询用户是否存在
	userModel := model.NewUser()
	user := userModel.GetUserById(id)
	if user == nil {
		return nil, e.NewBusinessError(e.FAILURE, "更新用户异常")
	}

	claims := us.newCustomClaims(user)
	accessToken, err := j.Refresh(claims)
	if err != nil {
		return nil, e.NewBusinessError(e.FAILURE, "生成Token失败")
	}

	return &TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresAt:   claims.ExpiresAt.Unix(),
	}, nil
}

// newAdminCustomClaims 初始化AdminCustomClaims
func (us UserService) newCustomClaims(user *model.User) j.CustomClaims[*model.User] {
	now := time.Now()
	expiresAt := now.Add(time.Second * c.Config.Jwt.TTL)
	return j.NewCustomClaims(user, expiresAt)
}
