package entity

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	FirebaseID string `json:"firebaseId"`
	Provider   string `json:"provider"`
}
