package dto

type TodoDTO struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"required,min=5"`
	UserID      string `json:"user_id,omitempty"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"required,min=5"`
}

type DeleteTodoRequest struct {
	ID string `json:"id"`
}

type TodoResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
}
