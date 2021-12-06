package model

type Todo struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Body        string `json:"body" db:"body"`
	UserId      string `json:"user_id" db:"user_id"`
	IsDone      bool   `json:"is_done" db:"isdone"`
	IsFavourite bool   `json:"is_favourite" db:"isfavourite"`
}

type TodoStore interface {
	Todo(todoId string) (*Todo, error)
	Todos(userId string) ([]Todo, error)
	Create(todo *Todo) (*Todo, error)
	Update(todo *Todo) error
	Delete(todoId string) error
	TodoPublic(todoId string) (string, error)
	TodoPublicGet(link string) (*Todo, error)
}
