package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type authKey struct{}

type AuthJwtClaims struct {
	UserId int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthJwtToken struct {
	UserId      int64
	JwtTokenId  string
	JwdSigned   string
	ExpiresTime time.Duration
}

func NewAuthJwtToken(userId int64, jwtKey string) (*AuthJwtToken, error) {

	//过期时间
	expiresTime := time.Hour * 24 * 30

	jwtID := uuid.New().String()
	//生成token，并入信息
	claims := AuthJwtClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        jwtID,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := jwtToken.SignedString([]byte(jwtKey))
	if err != nil {
		return nil, err
	}

	token := &AuthJwtToken{
		JwtTokenId:  jwtID,
		UserId:      userId,
		JwdSigned:   ss,
		ExpiresTime: expiresTime,
	}

	return token, nil
}

func ParseJwtToken(token string, jwtKey string) (*AuthJwtClaims, error) {
	claims := &AuthJwtClaims{}
	tokenClaims, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*AuthJwtClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func NewContext(ctx context.Context, info AuthJwtToken) context.Context {
	return context.WithValue(ctx, authKey{}, info)
}

// FromContext extract auth info from context
func FromContext(ctx context.Context) (tokenInfo AuthJwtToken, ok bool) {
	tokenInfo, ok = ctx.Value(authKey{}).(AuthJwtToken)
	return
}
