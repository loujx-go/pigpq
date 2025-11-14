package token

import (
	"errors"
	"pigpq/global"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"pigpq/config"
	e "pigpq/internal/pkg/errors"
)

// Generate 生成JWT Token
func Generate(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 生成签名字符串
	tokenStr, err := token.SignedString([]byte(config.Config.Jwt.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// Refresh 刷新JWT Token
func Refresh(claims jwt.Claims) (string, error) {
	return Generate(claims)
}

// Parse 解析token
func Parse(accessToken string, claims jwt.Claims, options ...jwt.ParserOption) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (i interface{}, err error) {
		// 确保算法正确
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		if config.Config.Jwt.SecretKey == "" {
			return nil, errors.New("JWT secret key is empty")
		}
		return []byte(config.Config.Jwt.SecretKey), nil
	}, options...)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, e.NewBusinessError(1, "invalid token")
	}

	return token, nil
}

// GetAccessToken 获取jwt的Token
func GetAccessToken(authorization string) (accessToken string, err error) {
	if authorization == "" {
		return "", errors.New("authorization header is missing")
	}

	// 检查 Authorization 头的格式
	if !strings.HasPrefix(authorization, "Bearer ") {
		return "", errors.New("invalid Authorization header format")
	}

	// 提取 Token 的值
	accessToken = strings.TrimPrefix(authorization, "Bearer ")
	return
}

// CustomClaims 自定义格式内容
type CustomClaims[T any] struct {
	Payload              T `json:"payload"`
	jwt.RegisteredClaims   // 内嵌标准的声明
}

// NewCustomClaims 初始化CustomClaims
func NewCustomClaims[T any](payload T, expiresAt time.Time) CustomClaims[T] {
	now := time.Now()
	return CustomClaims[T]{
		Payload: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    global.Issuer,                 // 签发人
			IssuedAt:  jwt.NewNumericDate(now),       // 签发时间
			ExpiresAt: jwt.NewNumericDate(expiresAt), // 定义过期时间
			NotBefore: jwt.NewNumericDate(now),       // 生效时间
			//ID:        rand.Int63n(1000000000),
			//Subject: global.Subject, // 签发主体
		},
	}
}
