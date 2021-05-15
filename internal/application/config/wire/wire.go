//go:generate wire
//+build wireinject

package wire

import (
	"context"

	"github.com/google/wire"
	"github.com/samthehai/chat/internal/application/config"
	"github.com/samthehai/chat/internal/application/services/server"
	"github.com/samthehai/chat/internal/domain/message"
	"github.com/samthehai/chat/internal/domain/user"
	"github.com/samthehai/chat/internal/infrastructure/external/redis"
	"github.com/samthehai/chat/internal/infrastructure/repository"
	"github.com/samthehai/chat/internal/infrastructure/repository/external"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver/commander"
)

var superSet = wire.NewSet(
	wire.InterfaceValue(new(context.Context), context.Background()),

	proviveRedisClientOption,
	proviveServerOption,

	wire.NewSet(
		redis.NewRedisClient,
		server.NewServer,
	),
	wire.NewSet(
		resolver.NewSubscriptionResolver,
		resolver.NewMutationResolver,
		resolver.NewQueryResolver,
		resolver.NewResolver,
	),

	wire.Bind(new(commander.MessageCommander), new(*message.MessageCommander)),
	wire.Bind(new(commander.UserCommander), new(*user.UserCommander)),
	wire.NewSet(
		message.NewMessageCommander,
		user.NewUserCommander,
	),

	wire.Bind(new(message.UserRepository), new(*repository.UserRepository)),
	wire.Bind(new(message.MessageRepository), new(*repository.MessageRepository)),
	wire.Bind(new(user.UserRepository), new(*repository.UserRepository)),
	wire.NewSet(
		repository.NewMessageRepository,
		repository.NewUserRepository,
	),

	wire.Bind(new(external.Cacher), new(*redis.RedisClient)),
)

var configObj = config.NewConfigFromEnv()

func proviveRedisClientOption() redis.RedisClientOption {
	return redis.RedisClientOption{
		Addr:     configObj.Redis.Addr,
		Password: configObj.Redis.Password,
	}
}

func proviveServerOption() server.ServerOption {
	return server.ServerOption{Port: configObj.HTTP.Port}
}

func InitializeServer() (server.Server, func(), error) {
	panic(wire.Build(superSet))
}
