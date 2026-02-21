package player

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/teomiscia/hexbattle/internal/model"
)

// Session represents an active player connection.
type Session struct {
	ID        string    `json:"id"`    // UUIDv4
	Token     string    `json:"token"` // 32-byte hex string
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`

	// Current game association (set when player joins a game)
	GameID   string `json:"-"`
	RoomCode string `json:"-"`
}

var nicknameRe = regexp.MustCompile(`^[a-zA-Z0-9_]{3,16}$`)

// ValidateNickname checks if a nickname meets the requirements:
// 3-16 characters, alphanumeric + underscores only.
func ValidateNickname(nickname string) error {
	nickname = strings.TrimSpace(nickname)
	if !nicknameRe.MatchString(nickname) {
		return fmt.Errorf("nickname must be 3-16 characters, alphanumeric and underscores only")
	}
	return nil
}

// SanitizeNickname trims and validates a nickname.
func SanitizeNickname(nickname string) (string, error) {
	nickname = strings.TrimSpace(nickname)
	if err := ValidateNickname(nickname); err != nil {
		return "", err
	}
	return nickname, nil
}

// NewSession creates a new guest player session with a generated ID and token.
func NewSession(nickname string) (*Session, error) {
	nickname, err := SanitizeNickname(nickname)
	if err != nil {
		return nil, err
	}

	id, err := generateUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate player ID: %w", err)
	}

	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &Session{
		ID:        id,
		Token:     token,
		Nickname:  nickname,
		CreatedAt: time.Now(),
	}, nil
}

// generateUUID creates a UUIDv4.
func generateUUID() (string, error) {
	var uuid [16]byte
	if _, err := rand.Read(uuid[:]); err != nil {
		return "", err
	}
	// Set version (4) and variant (RFC 4122)
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16]), nil
}

// generateToken creates a 32-byte hex string using crypto/rand.
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GuestResponse is the JSON response for guest registration.
type GuestResponse struct {
	PlayerID string `json:"player_id"`
	Token    string `json:"token"`
	Nickname string `json:"nickname"`
}

// ToResponse converts a session to a guest registration response.
func (s *Session) ToResponse() GuestResponse {
	return GuestResponse{
		PlayerID: s.ID,
		Token:    s.Token,
		Nickname: s.Nickname,
	}
}

// GenerateGameID creates a unique game ID.
func GenerateGameID() (string, error) {
	return generateUUID()
}

// GenerateUnitID creates a unique unit ID for troops.
func GenerateUnitID() string {
	id, err := generateUUID()
	if err != nil {
		// Fallback: use timestamp-based ID (should never happen with crypto/rand)
		return fmt.Sprintf("unit-%d", time.Now().UnixNano())
	}
	return id
}

// GenerateStructureID creates a unique structure ID.
func GenerateStructureID() string {
	id, err := generateUUID()
	if err != nil {
		return fmt.Sprintf("struct-%d", time.Now().UnixNano())
	}
	return id
}

// Validate that model package is reachable (compile-time check).
var _ = model.TroopMarine
