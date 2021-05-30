package entity

type User struct {
	ID            ID     `json:"id"`
	Name          string `json:"name"`
	PictureUrl    string `json:"pictureUrl"`
	FirebaseID    string `json:"firebaseId"`
	Provider      string `json:"provider"`
	EmailAddress  string `json:"emailAddress"`
	EmailVerified bool   `json:"emailVerified"`
}

func (u *User) GetID() ID {
	return u.ID
}
