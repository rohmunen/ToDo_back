package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testmod/internal/app/model"
	"testmod/internal/app/store"
	"testmod/pkg/auth"
	"testmod/pkg/middleware"
	"testmod/pkg/response"

	"github.com/gorilla/mux"
)

func AddUserHandlers(router *mux.Router, store *store.Store, auth *auth.Manager) {
	router.HandleFunc("/api/registration", middleware.SetupCORS(HandleUserCreate(store, auth))).Methods("POST", "OPTIONS")
}

func HandleUserCreate(store *store.Store, auth *auth.Manager) http.HandlerFunc {
	type jwtToken struct {
		Token string `json:"token"`
	}
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if (*r).Method == "OPTIONS" {
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			response.ErrorFunc(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
			Name:     req.Name,
		}
		if err := store.UserStore.Create(u); err != nil {
			response.ErrorFunc(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()

		jwt, err := auth.NewJWT(strconv.Itoa(u.Id), 10000)
		if err != nil {
			response.ErrorFunc(w, r, http.StatusInternalServerError, err)
		}

		response.Respond(w, r, http.StatusCreated, jwtToken{
			Token: jwt,
		})
	}
}
