package task

type Task struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	ProjectID     int    `json:"project_id"`
	UserID        int    `json:"user_id"`
	PostUserId    int    `json:"post_user_id"`
	CreatorUserId int    `json:"creator_user_id"`
}
