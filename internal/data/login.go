package data

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"go-cs/internal/biz"
	"go-cs/internal/conf"
	"go-cs/internal/utils"
	"go-cs/internal/utils/auth"
	"go-cs/internal/utils/rand"
	"go-cs/pkg/stream"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
)

type loginRepo struct {
	baseRepo
	confJwt *conf.Jwt
}

func NewLoginRepo(c *conf.Jwt, data *Data, logger log.Logger) biz.LoginRepo {
	moduleName := "LoginRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &loginRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
		confJwt: c,
	}
}

func (l loginRepo) GenerateJwtToken(ctx context.Context, userId int64) (string, error) {

	token, err := auth.NewAuthJwtToken(userId, l.confJwt.Key)
	if err != nil {
		return "", err
	}

	tokenKey := "session:" + token.JwtTokenId + ":user_id"
	userKey := "session:" + cast.ToString(userId) + ":tokens"

	_, err = l.data.rdb.Pipelined(ctx, func(pipeline redis.Pipeliner) error {
		pipeline.Set(ctx, tokenKey, userId, token.ExpiresTime)

		// 添加到用户token列表
		pipeline.ZAdd(ctx, userKey, &redis.Z{
			Score:  float64(time.Now().Add(token.ExpiresTime).Unix()),
			Member: token.JwtTokenId,
		})

		// 移除过期的token
		pipeline.ZRemRangeByScore(ctx, userKey, "0", cast.ToString(time.Now().Unix()))
		return nil
	})
	if err != nil {
		return "", err
	}

	return token.JwdSigned, nil
}

func (l loginRepo) ClearAllJwtToken(ctx context.Context, userId int64) error {
	userKey := "session:" + cast.ToString(userId) + ":tokens"
	tokens := l.data.rdb.ZRange(ctx, userKey, 0, -1).Val()

	tokenKeys := stream.Map(tokens, func(v string) string {
		return "session:" + v + ":user_id"
	})

	allKeys := append(tokenKeys, userKey)

	err := l.data.rdb.Del(ctx, allKeys...).Err()

	return err
}

func (l loginRepo) ClearJwtTokenByToken(ctx context.Context, token string) error {
	tokenKey := "session:" + token + ":user_id"

	userId := l.data.rdb.Get(ctx, tokenKey).Val()

	userKey := "session:" + userId + ":tokens"

	_, err := l.data.rdb.Pipelined(ctx, func(pipeline redis.Pipeliner) error {
		pipeline.Del(ctx, tokenKey)
		pipeline.ZRem(ctx, userKey, token)
		return nil
	})

	return err
}

func (l loginRepo) GenerateVerificationCode(ctx context.Context, name string) (string, error) {

	kTransport := utils.GetTransPortFormCtx(ctx)
	if kTransport == nil {
		return "", errors.New("failed")
	}

	session, err := l.data.sessionStore.New(kTransport, name)
	if err != nil {
		return "", err
	}

	validCode := rand.S(5)
	session.Values = map[any]any{
		"code": validCode,
	}

	saveErr := session.Save(kTransport)
	if saveErr != nil {
		return "", saveErr
	}

	return validCode, nil

}

func (l loginRepo) IncrAndGetCountCode(ctx context.Context, name string) (int64, error) {
	kTransport := utils.GetTransPortFormCtx(ctx)
	if kTransport == nil {
		return 0, errors.New("failed")
	}
	countSession, err := l.data.sessionStore.New(kTransport, name)
	if err != nil {
		return 0, err
	}

	c := cast.ToInt64(countSession.Values["c"]) + 1
	countSession.Values["c"] = c
	countSession.Save(kTransport)

	return c, nil
}

func (l loginRepo) CleanVerificationCode(ctx context.Context, name string) error {
	kTransport := utils.GetTransPortFormCtx(ctx)
	if kTransport == nil {
		return errors.New("failed")
	}

	session, _ := l.data.sessionStore.New(kTransport, name)

	//用完清理掉
	session.Values = map[any]any{}
	session.Options.MaxAge = -1
	session.Save(kTransport)

	return nil
}

func (l loginRepo) GetVerificationCode(ctx context.Context, name string) (string, error) {

	kTransport := utils.GetTransPortFormCtx(ctx)
	if kTransport == nil {
		return "", errors.New("failed")
	}

	session, err := l.data.sessionStore.New(kTransport, name)
	if err != nil {
		return "", err
	}

	return cast.ToString(session.Values["code"]), nil
}

func (l loginRepo) SavePfTokenInfo(ctx context.Context, pfToken string, info string) error {
	key := "balala:pf_token:" + pfToken

	err := l.data.rdb.Set(ctx, key, info, time.Minute*5).Err()

	return err
}

func (l loginRepo) GetPfTokenInfo(ctx context.Context, pfToken string) string {
	key := "balala:pf_token:" + pfToken

	val := l.data.rdb.Get(ctx, key).Val()

	return val
}

func (l loginRepo) DelPfTokenInfo(ctx context.Context, pfToken string) {
	key := "balala:pf_token:" + pfToken

	l.data.rdb.Del(ctx, key)
}
