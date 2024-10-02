package middles

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// FlowControlTag 统一用metadata封装，从http头部获取流控标志，并设置到context
func FlowControlTag() Middle {
	return func(c *gin.Context) {
		ctx := FromHTTPRequest(c.Request.Context(), c.Request)
		c.Request = c.Request.WithContext(ctx)
	}
}

func FromHTTPRequest(ctx context.Context, req *http.Request) context.Context {
	header := req.Header.Get(metadataStrKey)
	if header == "" {
		if c, err := req.Cookie(metadataStrKey); err == nil {
			header = c.Value
		}
	}

	if header == "" {
		return WithMetadata(ctx, make(Metadata))
	}

	//ret := make(Metadata)
	//_ = json.Unmarshal([]byte(header), &ret)
	ret := string2Md(header)
	return WithMetadata(ctx, ret)
}

type Metadata map[string]string
type metadataCtxKey struct{}

const (
	metadataStrKey      = "metadata"
	mdParisSeparator    = "||" // k-v pairs 之间的间隔符
	mdKeyValueSeparator = "="  // k-v pair, k和v的间隔符
)

func WithMetadata(ctx context.Context, metadata Metadata) context.Context {
	return context.WithValue(ctx, metadataCtxKey{}, metadata)
}

// string s format: k1=v1||k2=v2||k3=v3...
func string2Md(s string) Metadata {
	items := strings.Split(s, mdParisSeparator)
	n := len(items)
	md := make(Metadata, n)
	for _, item := range items {
		ss := strings.Split(item, mdKeyValueSeparator)
		if len(ss) != 2 {
			continue
		}
		md[ss[0]] = ss[1]
	}

	return md
}
