package validators

import "tasius.my.id/todolistapi/internal/application/dto"

func ValidateCreateTodo(req *dto.TodoDTO) []string {
	var errors []string

	if req.Title == "" {
		errors = append(errors, "Title is required")
	}

	if req.Description == "" {
		errors = append(errors, "Description is required")
	}

	return errors
}

func ValidateUpdateTodo(req *dto.UpdateTodoRequest) []string {
	var errors []string

	if req.Title == "" {
		errors = append(errors, "Title is required")
	}

	if req.Description == "" {
		errors = append(errors, "Description is required")
	}

	return errors
}

func ValidateDeleteTodo(req *dto.DeleteTodoRequest) []string {
	var errors []string

	if req.ID == "" {
		errors = append(errors, "ID is required")
	}

	return errors
}