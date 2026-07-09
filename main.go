package main

import (
	"fmt"
	"net/http"
	"task-tracker/middleware"
	"task-tracker/task"
	"task-tracker/user"
	"task-tracker/ws"
)

func main() {
	// 1. Инициализация (Seeding)
	user.InitAdmin()

	// 2. Настройка ключевых параметров
	jwtKey := []byte("super_secret_key")

	// 3. Создаем "обертку" для Middleware
	// Теперь каждый пакет получает готовый auth-обработчик
	auth := func(next http.HandlerFunc) http.HandlerFunc {
		return middleware.AuthMiddleware(jwtKey, next)
	}

	// 4. Регистрация маршрутов модулей
	user.RegisterUserRoutes(auth, jwtKey)
	task.RegisterTaskRoutes(auth)

	// 5. Отдельный роут для WebSocket
	http.HandleFunc("/ws", ws.HandleWebSocket)

	// 6. Запуск сервера
	fmt.Println("Сервер успешно запущен на порту :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %v\n", err)
	}
}
