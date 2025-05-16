package server3authfunc

import (
	"context"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go-cs/internal/conf"
	"go-cs/internal/utils"
	"go-cs/internal/utils/auth"
	"strings"

	"github.com/spf13/cast"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-redis/redis/v8"
)

const WebsocketSession = "websocket-session"

const (

	// bearerWord the bearer key word for authorization
	bearerWord string = "Bearer"

	// bearerFormat authorization token format
	bearerFormat string = "Bearer %s"

	// authorizationKey holds the key used to store the JWT Token in the request tokenHeader.
	authorizationKey string = "Authorization"

	// reason holds the error reason.
	reason string = "UNAUTHORIZED"
)

var (
	ErrInternalServer = errors.InternalServer("InternalServer", "server error")
	ErrMissingToken   = errors.Unauthorized(reason, "Token is missing")
	ErrTokenInvalid   = errors.Unauthorized(reason, "Token is invalid")
	ErrWrongContext   = errors.Unauthorized(reason, "Wrong context for middleware")
)

// AuthToken 驗證token並寫入ctx後回傳, 沒被Server middleware保護的api透過此方式驗證
func AuthToken(httpCtx context.Context, jwtConf *conf.Jwt, rdb *redis.Client) (context.Context, error) {
	req := utils.GetRequestFromTransport(httpCtx)
	if req == nil {
		return nil, ErrMissingToken
	}

	token, err := AuthTokenFromHttpReq(req, jwtConf, rdb)
	if err != nil {
		return nil, err
	}

	return auth.NewContext(httpCtx, *token), nil
}

func AuthTokenFromHttpReq(req *http.Request, jwtConf *conf.Jwt, rdb *redis.Client) (*auth.AuthJwtToken, error) {
	var (
		server3Token string
		sessionId    string
	)

	// 先看看是不是websocket转发的请求
	value := req.Context().Value(WebsocketSession)
	if value != nil {
		if s, ok := value.([]string); ok {
			userId := cast.ToInt64(s[0])
			sessionId = cast.ToString(s[1])
			tokenInfo := &auth.AuthJwtToken{
				UserId:     userId,
				JwtTokenId: sessionId,
			}

			return tokenInfo, nil
		}
	}

	if token := req.URL.Query().Get("token"); token != "" {
		// 从 url 取 token
		server3Token = token
	} else if token = req.Header.Get(authorizationKey); token != "" {
		// 从 Authorization 取 token
		auths := strings.SplitN(token, " ", 2)
		if len(auths) != 2 || !strings.EqualFold(auths[0], bearerWord) {
			return nil, ErrMissingToken
		}
		server3Token = auths[1]
	} else {
		// 从 cookie 取 token
		c, err := req.Cookie("token")
		if err != nil {
			return nil, ErrMissingToken
		}

		server3Token = c.Value
	}

	token, err := auth.ParseJwtToken(server3Token, jwtConf.Key)
	if err != nil {
		return nil, ErrMissingToken
	}

	sessionId = token.RegisteredClaims.ID

	// 從redis取userId, 取不到則代表未登入
	userId, err := rdb.Get(req.Context(), "session:"+sessionId+":user_id").Int64()
	if err == redis.Nil {
		return nil, ErrMissingToken
	}
	if err != nil {
		return nil, ErrInternalServer
	}

	if userId <= 0 {
		return nil, ErrTokenInvalid
	}
	if userId != token.UserId {
		return nil, ErrTokenInvalid
	}

	tokenInfo := &auth.AuthJwtToken{
		UserId:     userId,
		JwtTokenId: sessionId,
	}

	return tokenInfo, nil
}
