// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package wire

import (
	"context"
	"github.com/google/wire"
	"github.com/samthehai/chat/internal/application/config"
	"github.com/samthehai/chat/internal/application/services/server"
	"github.com/samthehai/chat/internal/application/services/server/middlewares"
	"github.com/samthehai/chat/internal/domain/usecase"
	repository2 "github.com/samthehai/chat/internal/domain/usecase/repository"
	"github.com/samthehai/chat/internal/infrastructure/external/auth"
	"github.com/samthehai/chat/internal/infrastructure/external/postgres"
	"github.com/samthehai/chat/internal/infrastructure/external/redis"
	"github.com/samthehai/chat/internal/infrastructure/repository"
	"github.com/samthehai/chat/internal/infrastructure/repository/external"
	"github.com/samthehai/chat/internal/interfaces/graph/loader"
	usecase3 "github.com/samthehai/chat/internal/interfaces/graph/loader/usecase"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver"
	loader2 "github.com/samthehai/chat/internal/interfaces/graph/resolver/loader"
	usecase2 "github.com/samthehai/chat/internal/interfaces/graph/resolver/usecase"
)

// Injectors from wire.go:

func InitializeServer() (server.Server, func(), error) {
	redisClientOption := proviveRedisClientOption()
	redisClient := redis.NewRedisClient(redisClientOption)
	authenticator := middlewares.NewAuthenticator()
	context := _wireContextValue
	connectionConfig := provivePostgresConnectionConfig()
	db := postgres.NewConnection(context, connectionConfig)
	userRepository := repository.NewUserRepository(redisClient, authenticator, db)
	messageRepository := repository.NewMessageRepository(redisClient, db)
	messageUsecase := usecase.NewMessageUsecase(userRepository, messageRepository)
	userUsecase := usecase.NewUserUsecase(userRepository)
	queryResolver := resolver.NewQueryResolver(messageUsecase, userUsecase)
	mutationResolver := resolver.NewMutationResolver(messageUsecase, userUsecase)
	subscriptionResolver := resolver.NewSubscriptionResolver(messageUsecase, userUsecase)
	userLoader := loader.NewUserLoader(userUsecase)
	conversationLoader := loader.NewConversationLoader(messageUsecase)
	messageResolver := resolver.NewMessageResolver(userLoader, conversationLoader)
	messageLoader := loader.NewMessageLoader(messageUsecase)
	conversationResolver := resolver.NewConversationResolver(messageLoader, userLoader, conversationLoader)
	userResolver := resolver.NewUserResolver(userLoader, conversationLoader)
	resolverResolver := resolver.NewResolver(queryResolver, mutationResolver, subscriptionResolver, messageResolver, conversationResolver, userResolver)
	string2 := proviveFirebaseCredentials()
	firebaseClient, err := auth.NewFirebaseClient(context, string2)
	if err != nil {
		return nil, nil, err
	}
	serverOption := proviveServerOption()
	serverServer, cleanup := server.NewServer(resolverResolver, firebaseClient, serverOption)
	return serverServer, func() {
		cleanup()
	}, nil
}

var (
	_wireContextValue = context.Background()
)

// wire.go:

var superSet = wire.NewSet(wire.InterfaceValue(new(context.Context), context.Background()), proviveRedisClientOption,
	provivePostgresConnectionConfig,
	proviveFirebaseCredentials,
	proviveServerOption, wire.NewSet(redis.NewRedisClient, postgres.NewConnection, auth.NewFirebaseClient, server.NewServer), wire.NewSet(resolver.NewSubscriptionResolver, resolver.NewMutationResolver, resolver.NewQueryResolver, resolver.NewMessageResolver, resolver.NewConversationResolver, resolver.NewUserResolver, resolver.NewResolver), wire.Bind(new(usecase2.MessageUsecase), new(*usecase.MessageUsecase)), wire.Bind(new(usecase2.UserUsecase), new(*usecase.UserUsecase)), wire.NewSet(usecase.NewMessageUsecase, usecase.NewUserUsecase), wire.Bind(new(repository2.UserRepository), new(*repository.UserRepository)), wire.Bind(new(repository2.MessageRepository), new(*repository.MessageRepository)), wire.NewSet(repository.NewMessageRepository, repository.NewUserRepository), wire.Bind(new(external.Cacher), new(*redis.RedisClient)), wire.Bind(new(external.Authenticator), new(*middlewares.Authenticator)), wire.NewSet(middlewares.NewAuthenticator), wire.Bind(new(loader2.MessageLoader), new(*loader.MessageLoader)), wire.Bind(new(loader2.ConversationLoader), new(*loader.ConversationLoader)), wire.Bind(new(loader2.UserLoader), new(*loader.UserLoader)), wire.NewSet(loader.NewMessageLoader, loader.NewConversationLoader, loader.NewUserLoader), wire.Bind(new(usecase3.MessageUsecase), new(*usecase.MessageUsecase)), wire.Bind(new(usecase3.UserUsecase), new(*usecase.UserUsecase)), wire.Bind(new(middlewares.AuthManager), new(*auth.FirebaseClient)),
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
