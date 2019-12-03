package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"

	core "cipherassets.core"
	"cipherassets.core/gql/resolvers"
	"cipherassets.core/gql/schema"
	"cipherassets.core/model"
	"cipherassets.core/store"
)

// Rest is a rest access server
type REST struct {
	config      *core.Config
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
	// router.Get("/api", func(w http.ResponseWriter, r *http.Request) {
	// 	if _, err := w.Write([]byte("graphql")); err != nil {
	// 		log.Printf("can't write to response: %s", err.Error())
	// 	}
	// 	// users, err := s.dataService.UserStore.GetUsers()
	// 	// if err != nil {
	// 	// 	log.Printf(err.Error())
	// 	// }
	// 	// log.Printf("%+v", users)
	// 	// userStr, err := json.Marshal(users)
	// 	// if err != nil {
	// 	// 	log.Printf("can't marshal users list: %s", err.Error())
	// 	// }
	// 	// log.Printf("%s", userStr)
	// 	//
	// 	// _ = users
	// })
	router.Use(s.AuthMiddleware)

	router.Handle("/playground", handler.Playground("GraphQL playground", "/api"))
	router.Handle("/api", handler.GraphQL(gql.NewExecutableSchema(gql.Config{Resolvers: resolvers.NewResolver(s.dataService)})))

	log.Printf("REST server runs on http://localhost:7000")
	return http.ListenAndServe(":7000", router)
}

func (s REST) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		var err error
		ctx := r.Context()
		if len(tokenStr) > 0 {
			token := model.AuthJWT{}
			if err = s.dataService.JWTStore.ParseJWTString(&tokenStr, &token); err == nil {
				ctx = context.WithValue(ctx, resolvers.AuthIdentityIDContextKey, token.AuthIdentityID)
			}

		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
