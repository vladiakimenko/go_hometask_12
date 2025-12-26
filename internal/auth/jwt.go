package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"secure-service/internal/service"
)

var jwtSecret []byte

func InitAuth() {
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) < 32 {
		panic("JWT_SECRET must be at least 32 characters long")
	}
}

func GenerateToken(user service.User) (string, error) {
	claims := service.Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign a token: %w", err)
	}
	return signed, nil
}

// Реализация KeyFunc для валидации токена
func CheckAlgorithm(token *jwt.Token) (any, error) {
	// проверяем алгоритм из хедера jwt
	// jwt токены можно создавать/редактировать вручную - это просто json в b64;
	// если в хедере указать "alg": "none", то jwt не будет проверять подпись вообще!
	if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
		return nil, fmt.Errorf("wrong encryption algorythm: %v", token.Header["alg"])
	}
	// вернуть надо секрет для чтения подписи
	return jwtSecret, nil
}

func ValidateToken(tokenString string) (*service.Claims, error) {
	claims := &service.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, CheckAlgorithm)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
