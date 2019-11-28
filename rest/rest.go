package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	core "cipherassets.core"
	"cipherassets.core/store"
)

// Rest is a rest access server
type REST struct {
	dataService *store.Store
}

func NewREST(config *core.Config) (*REST, error) {
	s, err := store.NewStore(config)
	if err != nil {
		return nil, fmt.Errorf("can't create store: %w", err)
	}

	return &REST{dataService: s}, nil
}

func (s REST) Serve() error {
	router := chi.NewRouter()
	router.Get("/api", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("graphql")); err != nil {
			log.Printf("can't write to response: %s", err.Error())
		}
		// users, err := s.dataService.UserStore.GetUsers()
		// if err != nil {
		// 	log.Printf(err.Error())
		// }
		// log.Printf("%+v", users)
		// userStr, err := json.Marshal(users)
		// if err != nil {
		// 	log.Printf("can't marshal users list: %s", err.Error())
		// }
		// log.Printf("%s", userStr)
		//
		// _ = users
	})
	log.Printf("REST server runs on http://localhost:7000")
	return http.ListenAndServe(":7000", router)
}
