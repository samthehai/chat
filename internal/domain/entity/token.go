package entity

type AuthToken struct {
	UserID        string `json:"user_id"`
	Name          string `json:"name"`
	PictureUrl    string `json:"picture_url"`
	Provider      string `json:"provider"`
	EmailAddress  string `json:"email_address"`
	EmailVerified bool   `json:"email_verified"`
}
