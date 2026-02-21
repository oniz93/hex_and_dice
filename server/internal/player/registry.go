package player

import (
	"fmt"
	"sync"
)

// Registry manages active player sessions in memory.
// It is safe for concurrent use.
type Registry struct {
	mu      sync.RWMutex
	byToken map[string]*Session // token -> session
	byID    map[string]*Session // player_id -> session
}

// NewRegistry creates an empty player registry.
func NewRegistry() *Registry {
	return &Registry{
		byToken: make(map[string]*Session),
		byID:    make(map[string]*Session),
	}
}

// Register adds a new session to the registry.
func (r *Registry) Register(session *Session) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byToken[session.Token] = session
	r.byID[session.ID] = session
}

// GetByToken returns the session for the given token, or nil if not found.
func (r *Registry) GetByToken(token string) *Session {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.byToken[token]
}

// GetByID returns the session for the given player ID, or nil if not found.
func (r *Registry) GetByID(id string) *Session {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.byID[id]
}

// Remove removes a session from the registry.
func (r *Registry) Remove(token string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if session, ok := r.byToken[token]; ok {
		delete(r.byID, session.ID)
		delete(r.byToken, token)
	}
}

// Count returns the number of active sessions.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.byToken)
}

// Authenticate validates a token and returns the session.
// Returns an error if the token is not found.
func (r *Registry) Authenticate(token string) (*Session, error) {
	session := r.GetByToken(token)
	if session == nil {
		return nil, fmt.Errorf("invalid or expired token")
	}
	return session, nil
}
