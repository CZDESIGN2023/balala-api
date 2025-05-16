package result

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-cs/internal/server/file/result/errs"
	"net/http"
)

func Ok(ctx *gin.Context, data any) {
	Custom(ctx, 200, data, "ok")
}

func Fail(ctx *gin.Context, err error) {

	var e *errs.Err
	switch {
	case errors.As(err, &e):
		Custom(ctx, e.Code, e.Data, e.Msg)
	default:
		Custom(ctx, 500, nil, err.Error())
	}
}

func Custom(ctx *gin.Context, code int, data any, msg string) {
	ctx.JSON(200, gin.H{
		"code":    code,
		"data":    data,
		"message": msg,
	})
}

func Forbidden(ctx *gin.Context) {
	ctx.Writer.WriteHeader(http.StatusForbidden)
	return
}
