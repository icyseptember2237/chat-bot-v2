package middles

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"gotemplate/logger"
	"gotemplate/server/httpserver/middles/status"
	"net/http"
	"reflect"
)

type requestKey struct{}
type responseKey struct{}
type ginContextKey struct{}

type Initer interface {
	Init(ctx context.Context)
}

var (
	EmptyStr = ""

	ErrMustPtr           = errors.New("param must be ptr")
	ErrMustPointToStruct = errors.New("param must point to struct")
	ErrMustHasThreeParam = errors.New("method must has three input")
	ErrMustFunc          = errors.New("method must be func")
	ErrMustValid         = errors.New("method must be valid")
	ErrMustError         = errors.New("method ret must be error or xderror")
	ErrMustOneOut        = errors.New("method must has one out")
	ErrWrongMethodType   = errors.New("method 格式不对")

	initerType     = reflect.TypeOf((*Initer)(nil)).Elem()
	replyErrorType = reflect.TypeOf((*error)(nil)).Elem()

	RequestKey    = requestKey{}
	ResponseKey   = responseKey{}
	GinContextKey = ginContextKey{}
)

// NewHandlerFuncFrom 从 GRPC handler 中创建 gin handler
// method 格式为 Method(ctx context.Context, req *ReqObj) (rsp *RspObj, err error)
func NewHandlerFuncFrom(method interface{}, opts ...Option) gin.HandlerFunc {
	return NewHandlerFuncWithLoggerFrom(method, opts...)
}

type createHandlerOptions struct {
	hasDataKey   bool
	notWriteResp bool
}

type Option func(opt *createHandlerOptions)

func getOption(opts ...Option) createHandlerOptions {
	option := createHandlerOptions{
		hasDataKey: true,
	}

	for _, opt := range opts {
		opt(&option)
	}

	return option
}

type RespStruct struct {
	Code    int64       `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func isImplementIniter(v reflect.Value) bool {
	return v.Type().Implements(initerType)
}

func callFieldInit(ctx context.Context, v reflect.Value) {
	elem := v.Elem()
	vT := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		ev := elem.Field(i)
		if isImplementIniter(ev) {
			if ev.CanSet() {
				ev.Set(reflect.New(vT.Field(i).Type.Elem()))
				initer := ev.Interface().(Initer)
				initer.Init(ctx)
			}
		}
	}
}

// MutateRequest mutate request
var MutateRequest func(r *http.Request) = func(r *http.Request) {}

// NewHandlerFuncWithLoggerFrom 从 method 创建 gin.HandlerFunc
// method 格式为 (h *Handler) Method(ctx context.Context, req *ReqObj) (rsp *RspObj, err error)
func NewHandlerFuncWithLoggerFrom(method interface{}, opts ...Option) gin.HandlerFunc {
	option := getOption(opts...)

	mV, reqT, err := check22Method(method)
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		MutateRequest(c.Request)
		ctx := c.Request.Context()
		req := reflect.New(reqT)

		if err := c.ShouldBind(req.Interface()); err != nil {
			logger.Warnf(ctx, "req: %v, err: %v\nbind param failed\n", c.Request.URL.Path, err)
			c.JSON(http.StatusBadRequest, RespStruct{Code: 1, Message: err.Error()})
			return
		}
		ctx = context.WithValue(ctx, RequestKey, c.Request)
		ctx = context.WithValue(ctx, ResponseKey, c.Writer)
		ctx = context.WithValue(ctx, GinContextKey, c)
		callFieldInit(ctx, req)

		logger.Debugf(ctx, "req: %v, func: %s\ninvoke handler\n", req, mV.Type().String())

		results := mV.Call([]reflect.Value{reflect.ValueOf(ctx), req})
		errValue := results[1]
		if errValue.IsValid() && !errValue.IsZero() && errValue.CanInterface() && errValue.Elem().IsValid() && !errValue.Elem().IsZero() {
			logger.Debugf(ctx, "url: %s\nhandler err: %v\n", c.Request.URL.Path, errValue)
			err := errValue.Interface().(error)
			if code, msg, ok := parseDefaultServiceCode(ctx, err); ok {
				sta := status.GetCode(err)
				c.JSON(sta, RespStruct{Code: code, Message: msg})
				return
			}
		}
		ret := results[0].Interface()
		if option.notWriteResp {
			return
		}
		if option.hasDataKey {
			ret = RespStruct{Code: 0, Data: ret}
		}
		if statusCode := c.Writer.Status(); statusCode != 0 {
			c.PureJSON(statusCode, ret)
			return
		}
		c.PureJSON(http.StatusOK, ret)
	}
}

// method 格式为 (h *Handler) Method(ctx context.Context, req *ReqObj) (rsp *RspObj, err error)
func check22Method(method interface{}) (mV reflect.Value, reqT reflect.Type, err error) {
	mV = reflect.ValueOf(method)
	if !mV.IsValid() {
		err = ErrMustValid
		return
	}

	mT := mV.Type()
	if mT.Kind() != reflect.Func {
		err = ErrMustFunc
		return
	}

	if mT.NumIn() != 2 {
		err = ErrWrongMethodType
		return
	}

	reqT = mT.In(1)
	if reqT.Kind() != reflect.Ptr {
		err = ErrMustPtr
		return
	}

	if reqT.Elem().Kind() != reflect.Struct {
		err = ErrMustPointToStruct
		return
	}
	reqT = reqT.Elem()

	if mT.NumOut() != 2 {
		err = ErrMustOneOut
		return
	}
	rspT := mT.Out(0)
	if rspT.Kind() != reflect.Ptr || rspT.Elem().Kind() != reflect.Struct {
		err = ErrMustPointToStruct
		return
	}

	errT := mT.Out(1)
	if errT != replyErrorType {
		err = ErrMustError
		return
	}
	return mV, reqT, err
}
