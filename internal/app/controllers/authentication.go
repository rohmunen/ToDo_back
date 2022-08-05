package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testmod/internal/app/model"
	"testmod/internal/app/store"
	"testmod/pkg/auth"
	"testmod/pkg/middleware"
	"testmod/pkg/response"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func AddAuthenticationHandlers(router *mux.Router, store *store.Store, auth *auth.Manager) {
	router.HandleFunc("/api/private/auth", middleware.SetupCORS(middleware.CheckAuth(HandleAuth(), auth))).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/login", middleware.SetupCORS(HandleLogin(store, auth))).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/forgot-password", middleware.SetupCORS(middleware.CheckAuth(HandleChangePassword(store, auth), auth))).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/recover-password", middleware.SetupCORS(HandlePasswordRecovery(store))).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/validate-password-recovery/", middleware.SetupCORS(HandleValidatePasswordRecovery(store))).Methods("POST", "OPTIONS")
}

func HandleAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if (*r).Method == "OPTIONS" {
			return
		}
		response.Respond(w, r, http.StatusOK, nil)
	}
}

func HandleChangePassword(store *store.Store, auth *auth.Manager) http.HandlerFunc {
	type request struct {
		Password    string `json:"password"`
		OldPassword string `json:"old_password"`
	}
	type jwtToken struct {
		Token string `json:"token"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
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
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			response.ErrorFunc(w, r, http.StatusBadRequest, err)
			return
		}
		u, err := store.UserStore.FindById(strIss)
		if err != nil || !u.ComparePassword(req.OldPassword) {
			response.ErrorFunc(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}
		jwt, err = auth.NewJWT(strconv.Itoa(u.Id), 10000)
		if err != nil {
			response.ErrorFunc(w, r, http.StatusInternalServerError, err)
		}
		store.UserStore.UpdatePassword(strIss, req.Password)
		response.Respond(w, r, http.StatusOK, jwtToken{Token: jwt})
	}
}

func HandlePasswordRecovery(store *store.Store) http.HandlerFunc {
	type request struct {
		Email string `json:"email"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if (*r).Method == "OPTIONS" {
			return
		}
		expiration := auth.CreateTimeout()
		pswr, hash, err := auth.CreateEmailVerHash()
		if err != nil {
			response.ErrorFunc(w, r, http.StatusInternalServerError, err)
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			response.ErrorFunc(w, r, http.StatusBadRequest, err)
			return
		}
		a := &model.VerificationRow{
			Email:      req.Email,
			Hash:       hash,
			Expiration: expiration,
		}
		id, err := store.RecoveryStore.Create(a)
		if err != nil {
			response.ErrorFunc(w, r, http.StatusInternalServerError, err)
			return
		}
		auth.Email(req.Email, "https://www.mysite.com/forgotpwchange/"+id+"/"+pswr)
	}
}

func HandleValidatePasswordRecovery(store *store.Store) http.HandlerFunc {
	type request struct {
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if (*r).Method == "OPTIONS" {
			return
		}
		req := &request{}
		id := r.URL.Query().Get("id")
		evpw := r.URL.Query().Get("evpw")
		hash, email, err := store.RecoveryStore.Get(id)
		if err != nil {
			response.ErrorFunc(w, r, http.StatusInternalServerError, err)
		}
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(evpw)); err != nil {
			response.ErrorFunc(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			response.ErrorFunc(w, r, http.StatusBadRequest, err)
			return
		}
		user, err := store.UserStore.FindByEmail(email)
		if err != nil {
			response.ErrorFunc(w, r, http.StatusBadRequest, err)
		}
		store.UserStore.UpdatePassword(strconv.Itoa(user.Id), req.Password)

		response.Respond(w, r, http.StatusOK, nil)
	}
}

func HandleLogin(store *store.Store, auth *auth.Manager) http.HandlerFunc {
	type jwtToken struct {
		Token string `json:"token"`
	}
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

		u, err := store.UserStore.FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			response.ErrorFunc(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		jwt, err := auth.NewJWT(strconv.Itoa(u.Id), 10000)
		if err != nil {
			response.ErrorFunc(w, r, http.StatusInternalServerError, err)
		}

		response.Respond(w, r, http.StatusOK, jwtToken{Token: jwt})
	}
}
