package websocket

import (
	"cloud/internal/service"
	"sync"

	"github.com/gorilla/websocket"
)

func init() {
	service.RegisterWebSocket(New())
}

type sWebSocket struct {
	mu                  sync.RWMutex
	pictureConnections  map[int64]map[*websocket.Conn]*WebSocketSession
	pictureEditingUsers map[int64]int64 // pictureId -> userId
}

type WebSocketSession struct {
	UserID    int64
	UserName  string
	PictureID int64
}

func New() *sWebSocket {
	return &sWebSocket{
		pictureConnections:  make(map[int64]map[*websocket.Conn]*WebSocketSession),
		pictureEditingUsers: make(map[int64]int64),
	}
}
