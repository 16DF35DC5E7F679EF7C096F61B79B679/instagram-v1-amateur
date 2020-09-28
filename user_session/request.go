package user_session

type CreateSessionRequestDTO struct {
	Handle string `json:"handle"`
	Password string `json:"password"`
	DeviceId string `json:"device_id"`
	DeviceIP string `json:"device_ip"`
	BrowserType string `json:"browser_type"`
	Timestamp int64 `json:"timestamp"`
}

type RegisterDeviceRequestDTO struct {
	DeviceName string `json:"device_name"`
	DeviceIP string `json:"device_ip"`
}