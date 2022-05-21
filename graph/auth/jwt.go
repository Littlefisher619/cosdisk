package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JwtCustomClaim struct {
	jwt.StandardClaims
}

type JwtManager struct {
	secret     []byte
	expireTime time.Duration
}

type contextKey string

const (
	contextKeyJwt = contextKey("JWT-AUTH")
	contextKeyUID = contextKey("UID")
)

func ExtractJwtFromContext(ctx context.Context) string {
	r, _ := ctx.Value(contextKeyJwt).(string)
	if r == "" {
		return ""
	}
	return r
}

func ExtractUserIdFromContext(ctx context.Context) string {
	r, _ := ctx.Value(contextKeyUID).(string)
	if r == "" {
		return ""
	}
	return r
}

func New() *JwtManager {
	return &JwtManager{
		secret:     []byte("secret"),
		expireTime: time.Hour * 24,
	}
}

func (j *JwtManager) GenerateToken(ctx context.Context, userID string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtCustomClaim{
		jwt.StandardClaims{
			Id:        userID,
			ExpiresAt: time.Now().Add(j.expireTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	token, err := t.SignedString(j.secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j *JwtManager) Validate(ctx context.Context, token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &JwtCustomClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there's a problem with the signing method")
		}
		return j.secret, nil
	})
}

func (j *JwtManager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuthToken := r.Header.Get("Authorization")
		bearer := "Bearer "
		if gotAuthToken == "" || len(gotAuthToken) <= len(bearer) {
			next.ServeHTTP(w, r)
			return
		}

		gotAuthToken = gotAuthToken[len(bearer):]
		ctx := context.WithValue(r.Context(), contextKeyJwt, gotAuthToken)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
