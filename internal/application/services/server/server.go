package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/samthehai/chat/internal/interfaces/graph/generated"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver"
)

const (
	graphqlEndpoint = "/query"
	defaultPort     = "8080"
)

type Server interface {
	Serve() error
}

type server struct {
	resolvers  resolver.Resolver
	httpServer *http.Server
}

func NewServer(resolvers resolver.Resolver) (Server, func()) {
	svr := &server{resolvers: resolvers}
	cleaner := func() {
		if svr.httpServer != nil {
			_ = svr.httpServer.Shutdown(context.Background())
		}
	}

	return svr, cleaner
}

func (s *server) Serve() error {
	// TODO: use config package
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Println("runnning server...")
	router := chi.NewRouter()
	s.registerRoutes(router)
	s.httpServer = &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: router}

	if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (s *server) registerRoutes(router *chi.Mux) {
	router.Handle("/", playground.Handler("GraphQL playground", graphqlEndpoint))
	router.Handle(graphqlEndpoint, s.newGraphQLServer())
}

func (s *server) newGraphQLServer() *handler.Server {
	return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &s.resolvers}))
}
