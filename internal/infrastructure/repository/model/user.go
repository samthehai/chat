package model

import (
	"github.com/samthehai/chat/internal/domain/entity"
)

// User model
type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	FirebaseID string `json:"firebase_id"`
	Provider   string `json:"provider"`
}

func ConvertModelUser(u *User) *entity.User {
	if u == nil {
		return nil
	}

	return &entity.User{
		ID:         u.ID,
		Name:       u.Name,
		FirebaseID: u.FirebaseID,
		Provider:   u.Provider,
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
