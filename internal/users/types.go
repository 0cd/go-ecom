package users

type CreateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Verified bool   `json:"verified"`
}

type UpdateUserPasswordParams struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
