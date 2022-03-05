// go:generate wire
//go:build wireinject
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
	"github.com/samthehai/chat/internal/infrastructure/repository/transactor"
	"github.com/samthehai/chat/internal/interfaces/graph/loader"
	loaderusecase "github.com/samthehai/chat/internal/interfaces/graph/loader/usecase"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver"
	resolverloader "github.com/samthehai/chat/internal/interfaces/graph/resolver/loader"
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
		resolver.NewConversationResolver,
		resolver.NewUserResolver,
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
	wire.Bind(new(usecaserepository.Transactor), new(*transactor.DBTransactor)),
	wire.NewSet(
		repository.NewMessageRepository,
		repository.NewUserRepository,
		transactor.NewDBTransactor,
	),

	wire.Bind(new(external.Cacher), new(*redis.RedisClient)),
	wire.Bind(new(external.Authenticator), new(*middlewares.Authenticator)),
	wire.Bind(new(external.Transactor), new(*transactor.DBTransactor)),
	wire.NewSet(
		middlewares.NewAuthenticator,
	),

	wire.Bind(new(resolverloader.MessageLoader), new(*loader.MessageLoader)),
	wire.Bind(new(resolverloader.ConversationLoader), new(*loader.ConversationLoader)),
	wire.Bind(new(resolverloader.UserLoader), new(*loader.UserLoader)),
	wire.NewSet(
		loader.NewMessageLoader,
		loader.NewConversationLoader,
		loader.NewUserLoader,
	),

	wire.Bind(new(loaderusecase.MessageUsecase), new(*usecase.MessageUsecase)),
	wire.Bind(new(loaderusecase.UserUsecase), new(*usecase.UserUsecase)),

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
		DebugUser:          configObj.Debug.User,
	}
}

func proviveFirebaseCredentials() string {
	return configObj.Firebase.Credentials
}

func InitializeServer() (server.Server, func(), error) {
	panic(wire.Build(superSet))
}
