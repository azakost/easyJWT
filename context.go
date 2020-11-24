package easyWeb

import (
	"context"
	"net/http"
)

type ctxKey string
type ctxData interface{}

func SetContext(r *http.Request, data ctxData) *http.Request {
	var key ctxKey = "data"
	ctx := context.WithValue(r.Context(), key, data)
	return r.WithContext(ctx)
}

func Ctx(r *http.Request) ctxData {
	if v := r.Context().Value("data"); v != nil {
		return v
	}
	return nil
}
