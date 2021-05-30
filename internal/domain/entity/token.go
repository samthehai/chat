package entity

type AuthToken struct {
	UserID        string
	Name          string
	PictureUrl    string
	Provider      string
	EmailAddress  string
	EmailVerified bool
}
