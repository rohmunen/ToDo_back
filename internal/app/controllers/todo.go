package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testmod/internal/app/model"
	"testmod/internal/app/store"
	"testmod/pkg/auth"
	"testmod/pkg/middleware"
	"testmod/pkg/response"

	"github.com/gorilla/mux"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
)

func AddTodoHandlers(router *mux.Router, store *store.Store, auth *auth.Manager) {
	router.HandleFunc("/api/todo", middleware.SetupCORS(middleware.CheckAuth(HandleTodoCreate(store, auth), auth))).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/todo/{id}", middleware.SetupCORS(middleware.CheckAuth(HandleTodoGet(store), auth))).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/todo/{id}", middleware.SetupCORS(middleware.CheckAuth(HandleTodoDelete(store), auth))).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/todo/{id}", middleware.SetupCORS(middleware.CheckAuth(HandleTodoPut(store), auth))).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/todo/{id}", middleware.SetupCORS(middleware.CheckAuth(HandleTodoCreatePublic(store), auth))).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/todos", middleware.SetupCORS(middleware.CheckAuth(HandleTodos(store, auth), auth))).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/todo/public/{link}", middleware.SetupCORS(HandleTodoGetPublic(store))).Methods("GET", "OPTIONS")
}

func HandleTodoCreatePublic(store *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//cors.SetupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		params := mux.Vars(r)
		id := params["id"]
		linkString, err := store.TodoStore.TodoPublic(id)
		if err != nil {
			response.ErrorFunc(w, r, http.StatusBadRequest, err)
		}
		response.Respond(w, r, http.StatusOK, linkString)
	}
}

func HandleTodoGetPublic(store *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//cors.SetupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		params := mux.Vars(r)
		link := params["link"]
		todo, err := store.TodoStore.TodoPublicGet(link)
		if err != nil {
			response.ErrorFunc(w, r, http.StatusInternalServerError, err)
		}
		response.Respond(w, r, http.StatusOK, todo)
	}
}

func HandleTodoPut(store *store.Store) http.HandlerFunc {
	type request struct {
		Kind  string `json:"type"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		//cors.SetupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		req := request{}
		params := mux.Vars(r)
		id := params["id"]
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			response.ErrorFunc(w, r, http.StatusInternalServerError, err)
		}
		if req.Kind == "favourite" {
			todo, err := store.TodoStore.Todo(id)
			if err != nil {
				response.ErrorFunc(w, r, http.StatusInternalServerError, err)
			}
			todo.IsFavourite = !todo.IsFavourite
			store.TodoStore.Update(todo)
		}
		if req.Kind == "done" {
			todo, err := store.TodoStore.Todo(id)
			if err != nil {
				response.ErrorFunc(w, r, http.StatusInternalServerError, err)
			}
			todo.IsDone = !todo.IsDone
			store.TodoStore.Update(todo)
		}
		if req.Kind == "upd" {
			todo, err := store.TodoStore.Todo(id)
			if err != nil {
				response.ErrorFunc(w, r, http.StatusInternalServerError, err)
			}
			todo.Title = req.Title
			todo.Body = req.Body
			store.TodoStore.Update(todo)
		}
		response.Respond(w, r, http.StatusOK, nil)
	}
}

func HandleTodoDelete(store *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//cors.SetupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		params := mux.Vars(r)
		id := params["id"]
		store.TodoStore.Delete(id)
	}
}

func HandleTodoGet(store *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//cors.SetupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		params := mux.Vars(r)
		id := params["id"]
		todo, err := store.TodoStore.Todo(id)
		if err != nil {
			response.ErrorFunc(w, r, http.StatusInternalServerError, err)
		}
		response.Respond(w, r, http.StatusOK, todo)
	}
}

func HandleTodos(store *store.Store, auth *auth.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//cors.SetupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		jwt := r.Header.Get("Authorization")
		ok, iss := auth.Parse(jwt)
		if !ok {
			response.ErrorFunc(w, r, http.StatusUnauthorized, errors.New("error parsing the jwt"))
			return
		}
		strIss := fmt.Sprintf("%v", iss)
		todos, err := store.TodoStore.Todos(strIss)
		if err != nil {
			response.ErrorFunc(w, r, http.StatusInternalServerError, err)
		}
		response.Respond(w, r, http.StatusOK, todos)
	}
}

func HandleTodoCreate(store *store.Store, auth *auth.Manager) http.HandlerFunc {
	type request struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		//cors.SetupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		jwt := r.Header.Get("Authorization")
		_, iss := auth.Parse(jwt)
		strIss := fmt.Sprintf("%v", iss)
		req := request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.ErrorFunc(w, r, http.StatusBadRequest, err)
			return
		}
		todo := model.Todo{
			Title:       req.Title,
			Body:        req.Body,
			UserId:      strIss,
			IsDone:      false,
			IsFavourite: false,
		}
		store.TodoStore.Create(&todo)
		response.Respond(w, r, http.StatusCreated, todo)
	}
}
