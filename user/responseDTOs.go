package user

type UserResponseDTO struct {
	Id string `json:"id"`
	Name string `json:"name"`
	SessionToken string `json:"session_token"`
}
