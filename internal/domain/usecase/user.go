package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/samthehai/chat/internal/domain/entity"
	domainerrors "github.com/samthehai/chat/internal/domain/errors"
	"github.com/samthehai/chat/internal/domain/usecase/repository"
)

type UserUsecase struct {
	userRepository repository.UserRepository
}

func NewUserUsecase(
	userRepository repository.UserRepository,
) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
	}
}

func (c *UserUsecase) CurrentUser(ctx context.Context) (*entity.User, error) {
	token, err := c.userRepository.GetAuthTokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("[User Repository] get auth token from context")
	}

	user, err := c.userRepository.FindByFirebaseID(ctx, token.UserID)
	if err != nil && !errors.Is(err, domainerrors.ErrNotFound) {
		return nil, fmt.Errorf("[User Repository] find by firebase id: %w", err)
	}

	if user == nil || errors.Is(err, domainerrors.ErrNotFound) {
		newUser := entity.User{
			Name:          token.Name,
			PictureUrl:    token.PictureUrl,
			EmailAddress:  token.EmailAddress,
			EmailVerified: token.EmailVerified,
			Provider:      token.Provider,
			FirebaseID:    token.UserID,
		}

		createdUser, err := c.userRepository.AddUser(ctx, newUser)
		if err == nil {
			return nil, fmt.Errorf("[User Repository] add user: %w", err)
		}

		user = createdUser
	}

	return user, nil
}

func (c *UserUsecase) Friends(ctx context.Context, first int, after entity.ID, sortBy entity.FriendsSortByType, sortOrder entity.SortOrderType) (*entity.UserFriendsConnection, error) {
	users, err := c.userRepository.FindFriends(ctx, first, after, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("[User Repository] find all: %w", err)
	}

	return users, nil
}

func (c *UserUsecase) UserJoined(ctx context.Context) (<-chan *entity.User, error) {
	user, err := c.userRepository.GetUserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("[User Repository] get user from context: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("[User Repository] user is nil")
	}

	users, err := c.userRepository.UserJoined(ctx, *user)
	if err != nil {
		return nil, fmt.Errorf("[User Repository] user joined: %w", err)
	}

	return users, nil
}

func (c *UserUsecase) User(ctx context.Context, id entity.ID) (*entity.User, error) {
	user, err := c.userRepository.FindUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[User Repository] find user: %w", err)
	}

	return user, nil
}
