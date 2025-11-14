package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

var (
	NoAuthHeaderError       = errors.New("no authorization header")
	BadSignMethodError      = errors.New("bad sign method")
	NoSecretKeyError        = errors.New("no secret key")
	InvalidTokenClaimsError = errors.New("invalid token claims")
)

var secretKey []byte

type Claims struct {
	Login string `json:"login"`
	Role  string `json:"role"`
	ID    uint   `json:"id"`
	jwt.RegisteredClaims
}

type JwtResponse struct {
	Token string `json:"token"`
}

func Init() error {
	secretKey = []byte(os.Getenv("SECRET_KEY"))
	if string(secretKey) == "" {
		return NoSecretKeyError
	}
	return nil
}

func ParseToken(r *http.Request) (*jwt.Token, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return nil, NoAuthHeaderError
	}
	inToken := strings.TrimPrefix(auth, "Bearer ")
	return jwt.Parse(inToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, BadSignMethodError
		}
		return secretKey, nil
	})
}

func ParseClaims(token *jwt.Token) (Claims, error) {
	claims := new(Claims)
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return *claims, InvalidTokenClaimsError
	}
	claims.Login = mapClaims["login"].(string)
	claims.Role = mapClaims["role"].(string)
	claims.ID = uint(mapClaims["id"].(float64))
	return *claims, nil
}

func GenerateToken(claims Claims) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(secretKey)
}
