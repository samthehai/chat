package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"github.com/samthehai/chat/internal/application/services/server/middlewares"
	"github.com/samthehai/chat/internal/interfaces/graph/generated"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver"
)

const (
	graphqlEndpoint = "/query"
)

type Server interface {
	Serve() error
}

type ServerOption struct {
	Port               int
	CORSAllowedOrigins []string
	Environment        string
	DebugUserID        string
}

type server struct {
	resolvers   resolver.Resolver
	authManager middlewares.AuthManager
	httpServer  *http.Server
	options     ServerOption
}

func NewServer(
	resolvers resolver.Resolver,
	authManager middlewares.AuthManager,
	options ServerOption,
) (Server, func()) {
	svr := &server{
		resolvers:   resolvers,
		authManager: authManager,
		options:     options,
	}
	cleaner := func() {
		if svr.httpServer != nil {
			_ = svr.httpServer.Shutdown(context.Background())
		}
	}

	return svr, cleaner
}

func (s *server) Serve() error {
	log.Printf("runnning server at port: %v ...\n", s.options.Port)

	router := chi.NewRouter()
	s.registerMiddlewares(router)
	s.registerRoutes(router)
	s.httpServer = &http.Server{Addr: fmt.Sprintf(":%v", s.options.Port), Handler: router}

	if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (s *server) registerMiddlewares(router *chi.Mux) {
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   s.options.CORSAllowedOrigins,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowCredentials: true,
	}).Handler)
	router.Use(middlewares.NewAuthenticationHandler(s.authManager, s.options.Environment == "development", s.options.DebugUserID))
}

func (s *server) registerRoutes(router *chi.Mux) {
	router.Handle("/", playground.Handler("GraphQL playground", graphqlEndpoint))
	router.Handle(graphqlEndpoint, s.newGraphQLServer())
}

func (s *server) newGraphQLServer() *handler.Server {
	return handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &s.resolvers}))
}
