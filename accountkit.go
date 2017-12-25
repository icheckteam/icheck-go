package icheck

// AccountKitLoginParams ...
type AccountKitLoginParams struct {
	Code     string
	Name     string
	Password string
	TTL      int64
}

// AccountKitResetPasswordParams ...
type AccountKitResetPasswordParams struct {
	Code     string
	Password string
}

// AccountKitResetPasswordResponse ...
type AccountKitResetPasswordResponse struct {
	Data map[string]interface{}
}
