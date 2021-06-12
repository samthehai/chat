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

func (u *UserUsecase) GetUserFromContext(ctx context.Context) (*entity.User, error) {
	user, err := u.userRepository.GetUserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user from context: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	return user, nil
}

func (u *UserUsecase) Login(ctx context.Context) (*entity.User, error) {
	token, err := u.userRepository.GetAuthTokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get auth token from context")
	}

	user, err := u.userRepository.FindByFirebaseID(ctx, token.UserID)
	if err != nil && !errors.Is(err, domainerrors.ErrNotFound) {
		return nil, fmt.Errorf("find by firebase id: %w", err)
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

		createdUser, err := u.userRepository.AddUser(ctx, newUser)
		if err != nil {
			return nil, fmt.Errorf("add user: %w", err)
		}

		user = createdUser
	}

	return user, nil
}

func (u *UserUsecase) Friends(ctx context.Context, first int, after entity.ID, sortBy entity.FriendsSortByType, sortOrder entity.SortOrderType) (*entity.UserFriendsConnection, error) {
	users, err := u.userRepository.FindFriends(ctx, first, after, sortBy, sortOrder)
	if err != nil {
		return nil, fmt.Errorf("find all: %w", err)
	}

	return users, nil
}

func (u *UserUsecase) UserJoined(ctx context.Context) (<-chan *entity.User, error) {
	user, err := u.userRepository.GetUserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get user from context: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	users, err := u.userRepository.UserJoined(ctx, *user)
	if err != nil {
		return nil, fmt.Errorf("user joined: %w", err)
	}

	return users, nil
}

func (u *UserUsecase) Me(ctx context.Context) (*entity.User, error) {
	token, err := u.userRepository.GetAuthTokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get auth token from context: %w", err)
	}

	user, err := u.userRepository.FindByFirebaseID(ctx, token.UserID)
	if err != nil && !errors.Is(err, domainerrors.ErrNotFound) {
		return nil, fmt.Errorf("find by firebase id: %w", err)
	}

	return user, nil
}

func (u *UserUsecase) Users(ctx context.Context, ids []entity.ID) ([]*entity.User, error) {
	users, err := u.userRepository.FindUsers(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	return users, nil
}
