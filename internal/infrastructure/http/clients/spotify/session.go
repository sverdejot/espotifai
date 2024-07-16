package spotify

import (
	"context"
	"time"
)

type tokenKey int

const key tokenKey = 1

func SetToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, key, token)
}

func GetToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(key).(string)
	return token, ok
}

type Session struct {
	Token    string
	ExpireAt time.Time
}

func (s Session) IsExpired() bool {
	return s.ExpireAt.Before(time.Now())
}
