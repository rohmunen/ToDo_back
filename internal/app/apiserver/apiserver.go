package apiserver

import (
	"net/http"
	"testmod/internal/app/store"
)

func Start(config *Config) error {
	store, err := store.NewStore(config.DbURL)
	if err != nil {
		return err
	}
	srv := newServer(store)
	return http.ListenAndServe(config.BindAddr, srv)
}
