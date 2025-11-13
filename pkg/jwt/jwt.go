package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

var (
	NoAuthHeaderError  = errors.New("no authorization header")
	BadSignMethodError = errors.New("bad sign method")
)

var secretKey []byte = []byte(os.Getenv("SECRET_KEY"))

type Claims struct {
	Login string `json:"login"`
	Role  string `json:"role"`
	ID    uint   `json:"id"`
	jwt.RegisteredClaims
}

type JwtResponse struct {
	Token string `json:"token"`
}

func ParseToken(r *http.Request) (*jwt.Token, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return nil, NoAuthHeaderError
	}
	inToken := strings.TrimPrefix(auth, "Bearer ")
	claims := new(Claims)
	return jwt.ParseWithClaims(inToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, BadSignMethodError
		}
		return secretKey, nil
	})
}

func GenerateToken(claims Claims) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(secretKey)
}
