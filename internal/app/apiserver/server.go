package apiserver

import (
	"net/http"
	"testmod/internal/app/controllers"
	"testmod/internal/app/store"
	"testmod/pkg/auth"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	store  *store.Store
	auth   *auth.Manager
}

func newServer(store *store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		store:  store,
		auth:   auth.NewManager("secret"),
	}
	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	controllers.AddAuthenticationHandlers(s.router, s.store, s.auth)
	controllers.AddTodoHandlers(s.router, s.store, s.auth)
	controllers.AddUserHandlers(s.router, s.store, s.auth)
}