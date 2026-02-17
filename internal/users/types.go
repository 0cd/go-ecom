package users

type CreateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Verified bool   `json:"verified"`
	IsAdmin  bool   `json:"is_admin"`
}
