// go:generate wire
// +build wireinject

package wire

import (
	"context"

	"github.com/google/wire"
	"github.com/samthehai/chat/internal/application/config"
	"github.com/samthehai/chat/internal/application/services/server"
	"github.com/samthehai/chat/internal/application/services/server/middlewares"
	usecase "github.com/samthehai/chat/internal/domain/usecase"
	usecaserepository "github.com/samthehai/chat/internal/domain/usecase/repository"
	"github.com/samthehai/chat/internal/infrastructure/external/auth"
	"github.com/samthehai/chat/internal/infrastructure/external/postgres"
	"github.com/samthehai/chat/internal/infrastructure/external/redis"
	"github.com/samthehai/chat/internal/infrastructure/repository"
	"github.com/samthehai/chat/internal/infrastructure/repository/external"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver"
	resolverusecase "github.com/samthehai/chat/internal/interfaces/graph/resolver/usecase"
)

var superSet = wire.NewSet(
	wire.InterfaceValue(new(context.Context), context.Background()),

	proviveRedisClientOption,
	provivePostgresConnectionConfig,
	proviveFirebaseCredentials,
	proviveServerOption,

	wire.NewSet(
		redis.NewRedisClient,
		postgres.NewConnection,
		auth.NewFirebaseClient,
		server.NewServer,
	),
	wire.NewSet(
		resolver.NewSubscriptionResolver,
		resolver.NewMutationResolver,
		resolver.NewQueryResolver,
		resolver.NewMessageResolver,
		resolver.NewResolver,
	),

	wire.Bind(new(resolverusecase.MessageUsecase), new(*usecase.MessageUsecase)),
	wire.Bind(new(resolverusecase.UserUsecase), new(*usecase.UserUsecase)),
	wire.NewSet(
		usecase.NewMessageUsecase,
		usecase.NewUserUsecase,
	),

	wire.Bind(new(usecaserepository.UserRepository), new(*repository.UserRepository)),
	wire.Bind(new(usecaserepository.MessageRepository), new(*repository.MessageRepository)),
	wire.NewSet(
		repository.NewMessageRepository,
		repository.NewUserRepository,
	),

	wire.Bind(new(external.Cacher), new(*redis.RedisClient)),
	wire.Bind(new(external.Authenticator), new(*middlewares.Authenticator)),
	wire.NewSet(
		middlewares.NewAuthenticator,
	),

	wire.Bind(new(middlewares.AuthManager), new(*auth.FirebaseClient)),
)

var configObj = config.NewConfigFromEnv()

func proviveRedisClientOption() redis.RedisClientOption {
	return redis.RedisClientOption{
		Addr:     configObj.Redis.Addr,
		Password: configObj.Redis.Password,
	}
}

func provivePostgresConnectionConfig() postgres.ConnectionConfig {
	return postgres.ConnectionConfig{
		Host:            configObj.Postgres.Host,
		Port:            configObj.Postgres.Port,
		User:            configObj.Postgres.User,
		Pass:            configObj.Postgres.Pass,
		Database:        configObj.Postgres.Database,
		ConnMaxLifetime: configObj.Postgres.ConnMaxLifetime,
		MaxIdleConns:    configObj.Postgres.MaxIdleConns,
		MaxOpenConns:    configObj.Postgres.MaxOpenConns,
	}
}

func proviveServerOption() server.ServerOption {
	return server.ServerOption{
		Port:               configObj.HTTP.Port,
		CORSAllowedOrigins: configObj.HTTP.CORSAllowedOrigins,
		Environment:        configObj.App.Environment,
		DebugUserID:        configObj.Debug.UserID,
	}
}

func proviveFirebaseCredentials() string {
	return configObj.Firebase.Credentials
}

func InitializeServer() (server.Server, func(), error) {
	panic(wire.Build(superSet))
}
