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
	user.InitAdmin()
	jwtKey := []byte("super_secret_key")

	auth := func(next http.HandlerFunc) http.HandlerFunc {
		return middleware.AuthMiddleware(jwtKey, next)
	}

	user.RegisterUserRoutes(auth, jwtKey)
	task.RegisterTaskRoutes(auth)

	http.HandleFunc("/ws", ws.HandleWebSocket)

	fmt.Println("Сервер успешно запущен на порту :8080")

	// ОБОРАЧИВАЕМ ВЕСЬ СЕРВЕР В CORS
	// http.DefaultServeMux содержит все зарегистрированные маршруты
	globalHandler := middleware.CorsMiddleware(http.DefaultServeMux)

	// Запускаем сервер с нашим globalHandler вместо nil
	if err := http.ListenAndServe(":8080", globalHandler); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %v\n", err)
	}
}
