package util

import (
	"github.com/gorilla/websocket"
)

type WebSocketStreamHandler struct {
	Conn *websocket.Conn
}

func (h *WebSocketStreamHandler) Read(p []byte) (int, error) {
	_, msg, err := h.Conn.ReadMessage()
	if err != nil {
		return 0, err
	}
	copy(p, msg)
	return len(msg), nil
}

func (h *WebSocketStreamHandler) Write(p []byte) (int, error) {
	err := h.Conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
