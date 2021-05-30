package auth

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/samthehai/chat/internal/domain/entity"
	"google.golang.org/api/option"
)

type FirebaseClient struct {
	authClient *auth.Client
}

func NewFirebaseClient(
	ctx context.Context,
	firebaseCredentials string,
) (*FirebaseClient, error) {
	app, err := firebase.NewApp(
		ctx,
		nil,
		option.WithCredentialsJSON([]byte(firebaseCredentials)),
	)
	if err != nil {
		return nil, fmt.Errorf("initialize firebase app: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("initialize firebase auth client: %w", err)
	}

	return &FirebaseClient{authClient: client}, nil
}

func (m *FirebaseClient) VerifyIDToken(
	ctx context.Context,
	idToken string,
) (*entity.AuthToken, error) {
	if idToken == "" {
		err := fmt.Errorf("verify id token: invalid idToken")

		return nil, err
	}

	token, err := m.authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		err := fmt.Errorf("verify id token: %w", err)

		return nil, err
	}

	return &entity.AuthToken{
		UserID:        token.UID,
		Name:          token.Claims["name"].(string),
		PictureUrl:    token.Claims["picture"].(string),
		Provider:      token.Firebase.SignInProvider,
		EmailAddress:  token.Claims["email"].(string),
		EmailVerified: token.Claims["email_verified"].(bool),
	}, nil
}
