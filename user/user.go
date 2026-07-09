package user

import (
	"fmt"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

var (
	Users  []User // Сделали большой буквы, чтобы main видел
	UserMu sync.Mutex
)

// InitAdmin создает начального админа
func InitAdmin() {
	adminPass := "super_secret_admin"
	adminHash, _ := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)

	Users = append(Users, User{
		ID:           1,
		Username:     "admin",
		PasswordHash: string(adminHash),
		IsAdmin:      true,
	})
	fmt.Println("Системный администратор успешно инициализирован.")
}
