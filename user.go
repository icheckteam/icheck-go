package icheck

type User struct {
	ID   int
	Name string `json:"social_name"`
}

type AccessToken struct {
	ID            string
	User          User
	TTL           int
	FirebaseToken string
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
	IcheckID []string
}
