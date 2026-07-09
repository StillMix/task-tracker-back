package task

import (
	"encoding/json"
	"net/http"
	"strconv"
	"task-tracker/ws"
)

// RegisterTaskRoutes регистрирует все маршруты задач
func RegisterTaskRoutes(authMiddleware func(http.HandlerFunc) http.HandlerFunc) {

	http.HandleFunc("GET /tasks", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		Mu.Lock()
		defer Mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Tasks)
	}))

	http.HandleFunc("POST /tasks", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		var t Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		Mu.Lock()
		t.ID = len(Tasks) + 1
		t.CreatorUserId = r.Context().Value("userID").(int)
		Tasks = append(Tasks, t)
		Mu.Unlock()

		_, _ = w.Write([]byte("Task created successfully"))

		// Используем Broadcast вместо ручного цикла по клиентам
		ws.Broadcast(t)
	}))

	http.HandleFunc("PUT /tasks/{id}", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		var updatedData Task
		_ = json.NewDecoder(r.Body).Decode(&updatedData)

		id, _ := strconv.Atoi(r.PathValue("id"))

		Mu.Lock()
		defer Mu.Unlock()

		for i, task := range Tasks {
			if task.ID == id {
				Tasks[i].Title = updatedData.Title
				Tasks[i].Description = updatedData.Description
				Tasks[i].ProjectID = updatedData.ProjectID
				Tasks[i].UserID = updatedData.UserID
				Tasks[i].PostUserId = updatedData.PostUserId

				// Используем Broadcast для оповещения об обновлении
				ws.Broadcast(Tasks[i])
				break
			}
		}
	}))
}
