package ws

import (
	"net/http"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var (
	Mu       sync.Mutex
	Clients  = make(map[*websocket.Conn]bool)
	Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	JwtKey = []byte("super_secret_key")
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 1. Проверяем наличие куки "jwt"
	cookie, err := r.Cookie("jwt")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. Валидируем токен
	token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 3. Если токен валиден, делаем Upgrade
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	Mu.Lock()
	Clients[conn] = true
	Mu.Unlock()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			Mu.Lock()
			delete(Clients, conn)
			Mu.Unlock()
			conn.Close()
			break
		}
	}
}

func Broadcast(message interface{}) {
	Mu.Lock()
	defer Mu.Unlock()
	for conn := range Clients {
		if err := conn.WriteJSON(message); err != nil {
			conn.Close()
			delete(Clients, conn)
		}
	}
}
