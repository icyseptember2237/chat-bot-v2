package middles

import (
	"context"
	"gotemplate/logger"
	"net/http"
)

// parseDefaultServiceCode 从返回 error 中解析服务 code（默认解析规则）
// 需要调用放保证 err 底层值非空
func parseDefaultServiceCode(ctx context.Context, err error) (code int64, msg string, ok bool) {
	if err == nil {
		return
	}
	defer func() {
		// 以防万一接口底层值为空时panic
		if e := recover(); e != nil {
			code = int64(http.StatusInternalServerError)
			ok = false
			logger.Warnf(ctx, "parseDefaultServiceCode recover:%v", e)
		}
	}()
	switch e := err.(type) {
	case CodeError:
		if e != nil {
			return int64(e.GetCode()), e.Error(), true
		}
	case CodeError64:
		if e != nil {
			return e.GetCode(), e.Error(), true
		}
	default:
		if e != nil {
			return int64(http.StatusInternalServerError), e.Error(), true
		}
	}
	return
}

// 支持int int64类型code
type CodeError interface {
	GetCode() int
	error
}
type CodeError64 interface {
	GetCode() int64
	error
}
