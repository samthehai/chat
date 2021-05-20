package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/samthehai/chat/internal/domain/entity"
)

type AuthManager interface {
	VerifyIDToken(ctx context.Context, idToken string) (*entity.AuthToken, error)
}

type Authenticator struct{}

func NewAuthenticator() *Authenticator {
	return &Authenticator{}
}

func (a *Authenticator) GetAuthTokenFromContext(
	ctx context.Context,
) (*entity.AuthToken, error) {
	token, ok := ctx.Value(accessKeyAuthToken).(*entity.AuthToken)

	if !ok || token == nil {
		return nil, fmt.Errorf("not authenticated")
	}

	return token, nil
}

func parseAuthorizationHeader(
	ctx context.Context,
	authManager AuthManager,
	tokenHeader string,
) (*entity.AuthToken, error) {
	idToken := strings.Replace(tokenHeader, "Bearer ", "", 1)

	if idToken == "" {
		return nil, fmt.Errorf("token not set")
	}

	token, err := authManager.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return token, nil
}

func NewAuthenticationHandler(
	authManager AuthManager,
	isDevelopment bool,
	debugUserID string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				token, err := parseAuthorizationHeader(
					r.Context(),
					authManager,
					r.Header.Get("Authorization"),
				)

				if isDevelopment && token == nil && len(debugUserID) > 0 {
					token = &entity.AuthToken{
						UserID: debugUserID,
					}
				} else if err != nil {
					http.Error(w, "invalid token", http.StatusUnauthorized)
					return
				}

				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), accessKeyAuthToken, token)))
			},
		)
	}
}
