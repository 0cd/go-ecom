package auth

type userLoginAndRegisterParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
