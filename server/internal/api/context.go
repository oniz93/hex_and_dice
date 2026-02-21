package api

import (
	"context"

	"github.com/teomiscia/hexbattle/internal/player"
)

type contextKey string

const sessionKey contextKey = "session"

// withSession adds the player session to the request context.
func withSession(ctx context.Context, session *player.Session) context.Context {
	return context.WithValue(ctx, sessionKey, session)
}

// SessionFromContext retrieves the player session from the request context.
func SessionFromContext(ctx context.Context) *player.Session {
	session, _ := ctx.Value(sessionKey).(*player.Session)
	return session
}
