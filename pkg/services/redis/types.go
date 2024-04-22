package redis

type AuthSession struct {
	Email     string `json:"email"`
	UserID    string `json:"user_id"`
	SessionId string `json:"session_id"`
	// For of auth session - staff, admin or external entities
	For string
}

type ReferenceSession struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}
