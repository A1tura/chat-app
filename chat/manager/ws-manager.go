package manager

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Manager struct {
	connections map[int]*websocket.Conn
	mu          sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		connections: make(map[int]*websocket.Conn),
	}
}

func (manager *Manager) Store(id int, con *websocket.Conn) {
    if con == nil {
        panic("xdd2")
    }
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.connections[id] = con
}

func (manager *Manager) Get(id int) *websocket.Conn {
	manager.mu.Lock()
	defer manager.mu.Unlock()
    if manager.connections[id] == nil {
        panic("xdd")
    }
	return manager.connections[id]
}

func (manager *Manager) Delete(id int) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	delete(manager.connections, id)
}

func (manager *Manager) ConSaved(id int) bool {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	_, exist := manager.connections[id]

	return exist
}
