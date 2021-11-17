package apiserver

import (
	"net/http"
	"testmod/internal/app/controllers"
	"testmod/internal/app/store"
	"testmod/pkg/auth"
	"testmod/pkg/middleware"

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
	s.router.HandleFunc("/user", controllers.HandleUserCreate(s.store, s.auth)).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/login", controllers.HandleLogin(s.store, s.auth)).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/recoverpw", controllers.HandlePasswordRecovery(s.store)).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/password/", controllers.HandleValidatePasswordRecovery(s.store)).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/changepw", middleware.CheckAuth(controllers.HandleChangePassword(s.store, s.auth), s.auth)).Methods("PUT", "OPTIONS")
	s.router.HandleFunc("/todo", middleware.CheckAuth(controllers.HandleTodoCreate(s.store, s.auth), s.auth)).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/todo/{id}", middleware.CheckAuth(controllers.HandleTodoGet(s.store), s.auth)).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/todo/{id}", middleware.CheckAuth(controllers.HandleTodoDelete(s.store), s.auth)).Methods("DELETE", "OPTIONS")
	s.router.HandleFunc("/todo/{id}", middleware.CheckAuth(controllers.HandleTodoPut(s.store), s.auth)).Methods("PUT", "OPTIONS")
	s.router.HandleFunc("/todo/{id}", middleware.CheckAuth(controllers.HandleTodoCreatePublic(s.store), s.auth)).Methods("POST", "OPTIONS")
	s.router.HandleFunc("/todos", middleware.CheckAuth(controllers.HandleTodos(s.store, s.auth), s.auth)).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/todo/public/{link}", controllers.HandleTodoGetPublic(s.store)).Methods("GET", "OPTIONS")
	s.router.HandleFunc("/private/auth", middleware.CheckAuth(controllers.HandleAuth(), s.auth)).Methods("GET", "OPTIONS")
}

/*
func setupCORS(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}
*/

/*
func (s *server) handleTodoCreatePublic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		params := mux.Vars(r)
		id := params["id"]
		linkString, err := s.store.TodoStore.TodoPublic(id)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		s.respond(w, r, http.StatusOK, linkString)
	}
}
*/

/*
func (s *server) handleTodoGetPublic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		params := mux.Vars(r)
		link := params["link"]
		todo, err := s.store.TodoStore.TodoPublicGet(link)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusOK, todo)
	}
}
*/

/*
func (s *server) handleTodoPut() http.HandlerFunc {
	type request struct {
		Kind  string `json:"type"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		req := request{}
		params := mux.Vars(r)
		id := params["id"]
		fmt.Println(r.Body)
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		fmt.Println("kind", req.Kind)
		if req.Kind == "favourite" {
			todo, err := s.store.TodoStore.Todo(id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
			todo.IsFavourite = !todo.IsFavourite
			s.store.TodoStore.Update(todo)
		}
		if req.Kind == "done" {
			todo, err := s.store.TodoStore.Todo(id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
			todo.IsDone = !todo.IsDone
			s.store.TodoStore.Update(todo)
		}
		if req.Kind == "upd" {
			todo, err := s.store.TodoStore.Todo(id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
			todo.Title = req.Title
			todo.Body = req.Body
			s.store.TodoStore.Update(todo)
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}
*/

/*
func (s *server) handleTodoDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		params := mux.Vars(r)
		id := params["id"]
		s.store.TodoStore.Delete(id)
	}
}
*/

/*
func (s *server) handleTodoGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		params := mux.Vars(r)
		id := params["id"]
		todo, err := s.store.TodoStore.Todo(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusOK, todo)
	}
}
*/

/*
func (s *server) handleTodos() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		jwt := r.Header.Get("Authorization")
		fmt.Println(jwt)
		ok, iss := s.auth.Parse(jwt)
		if ok != "ok" {
			s.error(w, r, http.StatusUnauthorized, errors.New("error parsing the jwt"))
			return
		}
		strIss := fmt.Sprintf("%v", iss)
		todos, err := s.store.TodoStore.Todos(strIss)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.respond(w, r, http.StatusOK, todos)
	}
}
*/

/*
func (s *server) handleTodoCreate() http.HandlerFunc {
	type request struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		jwt := r.Header.Get("Authorization")
		_, iss := s.auth.Parse(jwt)
		strIss := fmt.Sprintf("%v", iss)
		req := request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		fmt.Println(req.Body)
		todo := model.Todo{
			Title:       req.Title,
			Body:        req.Body,
			UserId:      strIss,
			IsDone:      false,
			IsFavourite: false,
		}
		fmt.Println(todo)
		s.store.TodoStore.Create(&todo)
		s.respond(w, r, http.StatusCreated, todo)
	}
}
*/

/*func (s *server) handleAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}*/

/*func (s *server) HandleChangePassword() http.HandlerFunc {
	type request struct {
		Password    string `json:"password"`
		OldPassword string `json:"old_password"`
	}
	type jwtToken struct {
		Token string `json:"token"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		s.respond(w, r, http.StatusOK, nil)
		jwt := r.Header.Get("Authorization")
		ok, iss := s.auth.Parse(jwt)
		if ok != "ok" {
			s.error(w, r, http.StatusUnauthorized, errors.New("error parsing the jwt"))
			return
		}
		strIss := fmt.Sprintf("%v", iss)
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u, err := s.store.UserStore.FindById(strIss)
		if err != nil || !u.ComparePassword(req.OldPassword) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}
		jwt, err = s.auth.NewJWT(strconv.Itoa(u.Id), 10000)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		s.store.UserStore.UpdatePassword(strIss, req.Password)
		s.respond(w, r, http.StatusOK, jwtToken{Token: jwt})
	}
} */

/*func (s *server) handlePasswordRecovery() http.HandlerFunc {
	type request struct {
		Email string `json:"email"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		expiration := auth.CreateTimeout()
		fmt.Println(expiration)
		pswr, hash, err := auth.CreateEmailVerHash()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		fmt.Println(pswr, hash)
		fmt.Println(r.Body)
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		fmt.Println("email", req.Email)
		a := &model.AccountRecovery{
			Email:      req.Email,
			Hash:       hash,
			Expiration: expiration,
		}
		id, err := s.store.RecoveryStore.Create(a)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		auth.Email(req.Email, "https://www.mysite.com/forgotpwchange/"+id+"/"+pswr)
	}
}*/

/*func (s *server) handleValidatePasswordRecovery() http.HandlerFunc {
	type request struct {
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		req := &request{}
		query := r.URL.Path
		fmt.Println(query)
		id := r.URL.Query().Get("id")
		evpw := r.URL.Query().Get("evpw")
		hash, email, err := s.store.RecoveryStore.Get(id)
		fmt.Println(email)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(evpw)); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		user, err := s.store.UserStore.FindByEmail(email)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		fmt.Println(strconv.Itoa(user.Id))
		s.store.UserStore.UpdatePassword(strconv.Itoa(user.Id), req.Password)

		s.respond(w, r, http.StatusOK, nil)
	}
} */

/*func (s *server) handleUserCreate() http.HandlerFunc {
	type jwtToken struct {
		Token string `json:"token"`
	}
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
			Name:     req.Name,
		}
		if err := s.store.UserStore.Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()

		jwt, err := s.auth.NewJWT(strconv.Itoa(u.Id), 10000)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}

		s.respond(w, r, http.StatusCreated, jwtToken{
			Token: jwt,
		})
	}
} */

/*func (s *server) handleLogin() http.HandlerFunc {
	type jwtToken struct {
		Token string `json:"token"`
	}
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.UserStore.FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		jwt, err := s.auth.NewJWT(strconv.Itoa(u.Id), 10000)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}

		s.respond(w, r, http.StatusOK, jwtToken{Token: jwt})
	}
}*/

/*func (s *server) checkAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		setupCORS(&w, req)
		if (*req).Method == "OPTIONS" {
			return
		}
		tokenstr := req.Header.Get("Authorization")
		fmt.Println(tokenstr)
		ok, _ := s.auth.Parse(tokenstr)
		if ok == "ok" {
			next.ServeHTTP(w, req)
		}
		w.WriteHeader(http.StatusUnauthorized)
	})
} */

/*func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}*/
