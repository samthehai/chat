package entity

type AuthToken struct {
	UserID        string
	Provider      string
	EmailAddress  string
	EmailVerified bool
}
