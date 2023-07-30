package model

const CtxCredentialsKey = "credentials"

type CredentialsSecure struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     Role   `json:"role"`
}

type Credentials struct {
	CredentialsSecure
	Password string `json:"-"`
}
