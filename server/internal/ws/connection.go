package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"nhooyr.io/websocket"

	"github.com/teomiscia/hexbattle/internal/model"
)

const (
	// SendChanSize is the buffer size for the outbound message channel.
	SendChanSize = 64

	// MaxMessageSize is the maximum allowed incoming message size in bytes.
	MaxMessageSize = 4096
)

// Connection wraps a WebSocket connection with read/write goroutines.
type Connection struct {
	PlayerID string
	Conn     *websocket.Conn
	SendChan chan []byte
	GameID   string

	ctx       context.Context
	cancel    context.CancelFunc
	closeOnce sync.Once
	closed    bool
	mu        sync.RWMutex

	// Callbacks
	OnMessage    func(playerID string, env Envelope)
	OnDisconnect func(playerID string)
}

// NewConnection creates a new managed WebSocket connection.
func NewConnection(ctx context.Context, conn *websocket.Conn, playerID string) *Connection {
	ctx, cancel := context.WithCancel(ctx)
	c := &Connection{
		PlayerID: playerID,
		Conn:     conn,
		SendChan: make(chan []byte, SendChanSize),
		ctx:      ctx,
		cancel:   cancel,
	}
	return c
}

// Start launches the read and write goroutines.
func (c *Connection) Start() {
	go c.readLoop()
	go c.writeLoop()
}

// Send queues a message for sending. Returns false if the connection is closed or the buffer is full.
func (c *Connection) Send(data []byte) bool {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return false
	}
	c.mu.RUnlock()

	select {
	case c.SendChan <- data:
		return true
	default:
		// Buffer full — connection is stuck, close it
		slog.Warn("send buffer full, closing connection",
			"player_id", c.PlayerID,
		)
		c.Close()
		return false
	}
}

// SendMessage marshals and sends a typed message.
func (c *Connection) SendMessage(msgType string, data interface{}) error {
	bytes, err := NewEnvelope(msgType, data)
	if err != nil {
		return err
	}
	if !c.Send(bytes) {
		return fmt.Errorf("connection closed or buffer full")
	}
	return nil
}

// SendAck sends an ACK for a client action.
func (c *Connection) SendAck(seq int, actionType string) error {
	return c.SendMessage(MsgAck, AckData{Seq: seq, ActionType: actionType})
}

// SendNack sends a NACK for a rejected client action.
func (c *Connection) SendNack(seq int, actionType string, code string, message string) error {
	return c.SendMessage(MsgNack, NackData{
		Seq:        seq,
		ActionType: actionType,
		Error:      ErrorData{Code: model.ErrorCode(code), Message: message},
	})
}

// Close terminates the connection.
func (c *Connection) Close() {
	c.closeOnce.Do(func() {
		c.mu.Lock()
		c.closed = true
		c.mu.Unlock()
		c.cancel()
		c.Conn.Close(websocket.StatusNormalClosure, "connection closed")
	})
}

// readLoop reads messages from the WebSocket and dispatches them.
func (c *Connection) readLoop() {
	defer func() {
		c.Close()
		if c.OnDisconnect != nil {
			c.OnDisconnect(c.PlayerID)
		}
	}()

	c.Conn.SetReadLimit(MaxMessageSize)

	for {
		_, data, err := c.Conn.Read(c.ctx)
		if err != nil {
			if c.ctx.Err() != nil {
				return // context cancelled, graceful shutdown
			}
			slog.Debug("websocket read error",
				"player_id", c.PlayerID,
				"error", err,
			)
			return
		}

		slog.Debug("websocket raw message",
			"player_id", c.PlayerID,
			"raw", string(data),
		)

		var env Envelope
		if err := json.Unmarshal(data, &env); err != nil {
			slog.Warn("malformed websocket message",
				"player_id", c.PlayerID,
				"error", err,
			)
			continue
		}

		if c.OnMessage != nil {
			c.OnMessage(c.PlayerID, env)
		}
	}
}

// writeLoop writes queued messages to the WebSocket.
func (c *Connection) writeLoop() {
	defer c.Close()

	for {
		select {
		case data, ok := <-c.SendChan:
			if !ok {
				return
			}
			ctx, cancel := context.WithTimeout(c.ctx, 10*time.Second)
			err := c.Conn.Write(ctx, websocket.MessageText, data)
			cancel()
			if err != nil {
				slog.Debug("websocket write error",
					"player_id", c.PlayerID,
					"error", err,
				)
				return
			}
		case <-c.ctx.Done():
			return
		}
	}
}
