package users

// type ActivateUserPayload struct {
// 	Username string `json:"username"`
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

type UpdateUserPayload struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}
