package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type NotificationService struct {
	clients map[string]*websocket.Conn // Store clients by user ID
	lock    sync.Mutex
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		clients: make(map[string]*websocket.Conn),
	}
}

func (ns *NotificationService) AddClient(userID string, conn *websocket.Conn) {
	ns.lock.Lock()
	defer ns.lock.Unlock()
	ns.clients[userID] = conn
}

func (ns *NotificationService) Broadcast(userIDs []string, message string) {
	ns.lock.Lock()
	defer ns.lock.Unlock()

	for _, userID := range userIDs {
		if conn, exists := ns.clients[userID]; exists {
			conn.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}
}
