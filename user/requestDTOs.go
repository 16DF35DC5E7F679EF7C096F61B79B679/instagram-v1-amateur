package user

type UserCreationRequestDTO struct {
	Name string `json:"name"`
	DOB int64 `json:"dob"`
	Email string `json:"email"`
	Username string `json:"username"`
	Handle string `json:"handle"`
	Password string `json:"password"`
	DeviceId string `json:"device_id"`
	BrowserType string `json:"browser_type"`
}

type SignInRequestDTO struct {
	Handle string `json:"handle"`
	Password string `json:"password"`
	DeviceId string `json:"device_id"`
	BrowserType string `json:"browser_type"`
}
