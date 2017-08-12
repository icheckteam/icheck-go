package icheck

// User is the resource representing your icheck user.
type User struct {
	ID            int    `json:"id"`
	IcheckID      string `json:"icheck_id"`
	Name          string `json:"social_name"`
	Type          string `json:"social_type"`
	SocialID      string `json:"social_id"`
	Avatar        string `json:"avatar"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"string,omitempty"`
	Cover         string `json:"cover"`
	EmailVerified bool   `json:"email_verified"`
	PhoneVerified bool   `json:"phone_verified"`
}

// Owner is the structure for an account owner.
type Owner struct {
	ID     int    `json:"id"`
	Name   string `json:"icheck_id"`
	Avatar string `json:"social_name"`
	Cover  string `json:"cover"`
}

// AccessToken ...
type AccessToken struct {
	ID            string `json:"id"`
	User          *User  `json:"user"`
	TTL           int    `json:"ttl"`
	FirebaseToken string `json:"firebase_token"`
}

// UserListParams ...
type UserListParams struct {
	ListParams
	IcheckID []string
}

// UserList ...
type UserList struct {
	Data []User
}
