package auth

import (
	"context"
	"time"
)

type tokenKey int

const key tokenKey = 1

func SetSpotifyToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, key, token)
}

func GetSpotifyToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(key).(string)
	return token, ok
}

type session struct {
	Token    string
	ExpireAt time.Time
}

func (s session) IsExpired() bool {
	return s.ExpireAt.Before(time.Now())
}
