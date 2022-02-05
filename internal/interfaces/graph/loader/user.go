package loader

import (
	"context"
	"fmt"

	"github.com/graph-gophers/dataloader"
	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/interfaces/graph/loader/usecase"
)

type UserLoader struct {
	usersLoader     *dataloader.Loader
	friendIDsLoader *dataloader.Loader
}

func NewUserLoader(
	userFetcher usecase.UserUsecase,
) *UserLoader {
	return &UserLoader{
		usersLoader:     newUsersLoader(userFetcher.Users),
		friendIDsLoader: newFriendIDsLoader(userFetcher.GetFriendIDsFromUserIDs),
	}
}

func (l *UserLoader) LoadUser(
	ctx context.Context,
	userID entity.ID,
) (*entity.User, error) {
	raw, err := l.usersLoader.Load(ctx, userID)()
	if err != nil {
		return nil, fmt.Errorf("load user: id=%v, %w", userID, err)
	}

	return raw.(*entity.User), nil
}

func (l *UserLoader) LoadFriendIDs(ctx context.Context,
	input entity.RelayQueryInput) (*entity.IDsConnection, error) {
	raw, err := l.friendIDsLoader.Load(ctx, input)()
	if err != nil {
		return nil, fmt.Errorf("load friend ids: userid=%v, %w", input.KeyID, err)
	}

	return raw.(*entity.IDsConnection), nil
}

func newUsersLoader(
	fetchFunc func(ctx context.Context, userIDs []entity.ID) ([]*entity.User, error),
) *dataloader.Loader {
	return dataloader.NewBatchedLoader(
		func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
			ids := getIDsFromKeys(keys)

			users, err := fetchFunc(ctx, ids)
			if err != nil {
				return fillUpResultsWithError(len(keys), err)
			}

			m := make(map[entity.ID]*entity.User)
			for _, u := range users {
				m[u.ID] = u
			}

			results := make([]*dataloader.Result, 0, len(keys))
			for _, id := range ids {
				d := m[id]

				results = append(results, &dataloader.Result{
					Data:  d,
					Error: err,
				})
			}

			return results
		},
	)
}

func newFriendIDsLoader(
	fetchFunc func(ctx context.Context, inputs []entity.RelayQueryInput,
	) (map[entity.ID]*entity.IDsConnection, error),
) *dataloader.Loader {
	return dataloader.NewBatchedLoader(
		func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
			inputs := make([]entity.RelayQueryInput, 0, len(keys))

			for _, key := range keys {
				inputs = append(inputs, key.Raw().(entity.RelayQueryInput))
			}

			idsConnection, err := fetchFunc(ctx, inputs)
			if err != nil {
				return fillUpResultsWithError(len(keys), err)
			}

			results := make([]*dataloader.Result, 0, len(keys))
			for _, input := range inputs {
				d := idsConnection[input.KeyID]

				results = append(results, &dataloader.Result{
					Data:  d,
					Error: err,
				})
			}

			return results
		},
	)
}
