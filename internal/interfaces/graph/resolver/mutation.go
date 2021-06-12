package resolver

import (
	"context"
	"fmt"

	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/interfaces/graph/model"
	"github.com/samthehai/chat/internal/interfaces/graph/resolver/usecase"
)

type MutationResolver struct {
	messageUsecase usecase.MessageUsecase
	userUsecase    usecase.UserUsecase
}

func NewMutationResolver(
	messageUsecase usecase.MessageUsecase,
	userUsecase usecase.UserUsecase,
) *MutationResolver {
	return &MutationResolver{
		messageUsecase: messageUsecase,
		userUsecase:    userUsecase,
	}
}

func (r *MutationResolver) CreateNewConversation(ctx context.Context, input model.CreateNewConversationInput) (*model.CreateNewConversationPayload, error) {
	user, err := r.userUsecase.GetUserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from context: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("failed to get user from context: user is nil")
	}

	conversationType := entity.ConversationTypeSingle
	if len(input.RecipentIDList) > 1 {
		conversationType = entity.ConversationTypeGroup
	}

	createdConversation, err := r.messageUsecase.CreateNewConversation(ctx, user.ID, input.Title, conversationType, input.RecipentIDList, input.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to create new conversation and post message: %v", err)
	}

	return &model.CreateNewConversationPayload{Conversation: createdConversation}, nil
}

func (r *MutationResolver) PostMessage(ctx context.Context, input model.PostMessageInput) (*model.PostMessagePayload, error) {
	user, err := r.userUsecase.GetUserFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from context: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("failed to get user from context: user is nil")
	}

	message, err := r.messageUsecase.PostMessage(ctx, input.ConversationID, entity.MessageTypeText, user.ID, input.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to post message: %w", err)
	}

	return &model.PostMessagePayload{
		Message: message,
	}, nil
}

func (r *MutationResolver) Login(ctx context.Context) (*entity.User, error) {
	return r.userUsecase.Login(ctx)
}
