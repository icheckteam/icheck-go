package icheck

type User struct {
	ID            int    `json:"id"`
	IcheckID      string `json:"icheck_id"`
	Avatar        string `json:"avatar"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	PhoneVerified bool   `json:"phone_verified"`
	EmailVerified bool   `json:"email_verified"`
}

type AccessToken struct {
	ID            string
	User          User
	TTL           int
	FirebaseToken string `json:"firebase_token"`
}

// LoginResponse
type LoginResponse struct {
	Data *AccessToken
}

// UserResponse
type UserResponse struct {
	User *User `json:"data"`
}

// LoginParams ...
type LoginParams struct {
	Username string
	Password string
	TTL      int64
}

type UserListResponse struct {
	Users []User `json:"data"`
}

type UserListParams struct {
	Params
	IcheckID []string
}

// RegisterParams
type RegisterParams struct {
	Username string
	Password string
	Name     string
}

type LoginSocialParams struct {
	Provider string
	Code     string
	TTL      int64
}
