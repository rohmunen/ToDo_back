package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"testmod/pkg/auth"
	"testmod/pkg/cors"
	"testmod/pkg/response"
)

func CheckAuth(next http.HandlerFunc, auth *auth.Manager) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cors.SetupCORS(&w, req)
		if (*req).Method == "OPTIONS" {
			return
		}
		tokenstr := req.Header.Get("Authorization")
		fmt.Println(tokenstr)
		ok, _ := auth.Parse(tokenstr)
		if ok {
			next.ServeHTTP(w, req)
		} else {
			response.ErrorFunc(w, req, http.StatusUnauthorized, errors.New("token is not valid"))
		}
	})
}
