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
	Todo(id string) (*Todo, error)
	Todos(id string) ([]Todo, error)
	Create(t *Todo) (*Todo, error)
	Update(t *Todo) error
	Delete(id string) error
	TodoPublic(id string) (string, error)
	TodoPublicGet(link string) (*Todo, error)
}
