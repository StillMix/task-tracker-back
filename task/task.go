package task

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	Tasks   []Task
	Mu      sync.Mutex
	Clients = make(map[*websocket.Conn]bool)
)
