package jwt

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/golang-jwt/jwt/v4"
)

// JWT secret key
var jwtKey = []byte("tiktok")

// Claims结构体用于定义JWT载荷
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateJWT生成JWT
func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute) // 设置过期时间
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // 设置过期时间
			Issuer:    "tiktok",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Println("new token")
	return token.SignedString(jwtKey) // 使用密钥签名JWT
}

// ValidateJWT验证JWT
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil // 返回用于验证的密钥
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

// AuthMiddleware是JWT验证的中间件
func AuthMiddleware(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	if token == "" {
		token = c.PostForm("token")
	}
	if token == "" {
		c.JSON(http.StatusUnauthorized, "Missing token")
		log.Println("Missing token")
		c.Abort()
		return
	}

	claims, err := ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Invalid token")
		log.Println("Invalid token", err)
		c.Abort()
		return
	}

	// 将解析后的claims存储在上下文中，供后续使用
	c.Set("username", claims.Username)
	c.Next(ctx)
}
