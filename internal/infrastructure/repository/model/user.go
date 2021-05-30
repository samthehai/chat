package model

import (
	"github.com/samthehai/chat/internal/domain/entity"
)

// User model
type User struct {
	ID            entity.ID `json:"id"`
	Name          string    `json:"name"`
	PictureUrl    string    `json:"picture_url"`
	FirebaseID    string    `json:"firebase_id"`
	Provider      string    `json:"provider"`
	EmailAddress  string    `json:"email_address"`
	EmailVerified bool      `json:"email_verified"`
}

func ConvertModelUser(u *User) *entity.User {
	if u == nil {
		return nil
	}

	return &entity.User{
		ID:            u.ID,
		Name:          u.Name,
		PictureUrl:    u.PictureUrl,
		FirebaseID:    u.FirebaseID,
		Provider:      u.Provider,
		EmailAddress:  u.EmailAddress,
		EmailVerified: u.EmailVerified,
	}
}

func ConvertModelUsers(users []*User) []*entity.User {
	if users == nil {
		return nil
	}

	uu := make([]*entity.User, 0, len(users))
	for _, u := range users {
		uu = append(uu, ConvertModelUser(u))
	}

	return uu
}

func GetColumnNameByFriendsSortByType(t entity.FriendsSortByType) string {
	switch t {
	case entity.FriendsSortByTypeName:
		return "name"
	default:
		return ""
	}
}
