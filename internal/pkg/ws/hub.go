package ws

import (
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Event 通知事件
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	Time    string          `json:"time"`
}

// Client WebSocket 客户端连接
type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

// Hub 管理所有 WebSocket 连接
type Hub struct {
	clients    map[string]map[*Client]bool // userID → clients
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

var DefaultHub = NewHub()

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
	}
}

// Run 启动 Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true
			h.mu.Unlock()
			slog.Debug("ws client connected", "user_id", client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.UserID]; ok {
				delete(clients, client)
				if len(clients) == 0 {
					delete(h.clients, client.UserID)
				}
			}
			close(client.Send)
			h.mu.Unlock()
			slog.Debug("ws client disconnected", "user_id", client.UserID)
		}
	}
}

// SendToUser 向指定用户发送事件
func (h *Hub) SendToUser(userID string, eventType string, payload interface{}) {
	data, _ := json.Marshal(payload)
	event := Event{
		Type:    eventType,
		Payload: data,
		Time:    time.Now().Format(time.RFC3339),
	}
	msg, _ := json.Marshal(event)

	h.mu.RLock()
	defer h.mu.RUnlock()
	if clients, ok := h.clients[userID]; ok {
		for client := range clients {
			select {
			case client.Send <- msg:
			default:
				close(client.Send)
				delete(clients, client)
			}
		}
	}
}

// Broadcast 向所有用户广播
func (h *Hub) Broadcast(eventType string, payload interface{}) {
	data, _ := json.Marshal(payload)
	event := Event{
		Type:    eventType,
		Payload: data,
		Time:    time.Now().Format(time.RFC3339),
	}
	msg, _ := json.Marshal(event)

	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, clients := range h.clients {
		for client := range clients {
			select {
			case client.Send <- msg:
			default:
				close(client.Send)
				delete(clients, client)
			}
		}
	}
}

// NotifyApplicationStatus 投递状态变更通知
func NotifyApplicationStatus(candidateID, status string) {
	DefaultHub.SendToUser(candidateID, "application.updated", map[string]string{"status": status})
}

// NotifyNewInterview 新面试邀约通知
func NotifyNewInterview(candidateID string) {
	DefaultHub.SendToUser(candidateID, "interview.created", map[string]string{"message": "您收到了新的面试邀约"})
}

// NotifyInterviewStatus 面试状态变更通知
func NotifyInterviewStatus(hrID, status string) {
	DefaultHub.SendToUser(hrID, "interview.updated", map[string]string{"status": status})
}
