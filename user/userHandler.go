package user

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUserRoutes регистрирует маршруты пакета user
func RegisterUserRoutes(authMiddleware func(http.HandlerFunc) http.HandlerFunc, jwtKey []byte) {

	http.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		var loginReq RegisterLoginRequest
		_ = json.NewDecoder(r.Body).Decode(&loginReq)

		var foundUser *User
		UserMu.Lock()
		for i := range Users {
			if Users[i].Username == loginReq.Username {
				foundUser = &Users[i]
				break
			}
		}
		UserMu.Unlock()

		if foundUser == nil || bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(loginReq.Password)) != nil {
			http.Error(w, "Не верный логин или пароль", http.StatusUnauthorized)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": foundUser.ID})
		tokenString, _ := token.SignedString(jwtKey)

		// Устанавливаем куку
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    tokenString,
			HttpOnly: true,  // Защита от XSS
			Secure:   false, // Поставь true, если будешь использовать HTTPS
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
		})

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("POST /users", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		requesterID := r.Context().Value("userID").(int)

		isAdmin := false
		UserMu.Lock()
		for _, u := range Users {
			if u.ID == requesterID && u.IsAdmin {
				isAdmin = true
				break
			}
		}
		UserMu.Unlock()

		if !isAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		var req RegisterLoginRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

		UserMu.Lock()
		Users = append(Users, User{ID: len(Users) + 1, Username: req.Username, PasswordHash: string(hash)})
		UserMu.Unlock()
		w.WriteHeader(http.StatusCreated)
	}))

	// Внутри RegisterUserRoutes добавь этот маршрут:

	http.HandleFunc("POST /logout", func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем куку с истекшим временем жизни
		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Logged out successfully"))
	})
}
