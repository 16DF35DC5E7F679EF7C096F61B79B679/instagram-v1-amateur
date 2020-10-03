package user_session

type CreateSessionRequestDTO struct {
	Handle string `json:"handle"`
	Password string `json:"password"`
	DeviceId string `json:"device_id"`
	BrowserType string `json:"browser_type"`
	Timestamp int64 `json:"timestamp"`
}

type RegisterDeviceRequestDTO struct {
	DeviceName string `json:"device_name"`
	DeviceIP string `json:"device_ip"`
}

type InvalidateSessionRequestDTO struct {
	Handle string `json:"handle"`
	DeviceId string `json:"device_id"`
	BrowserType string `json:"browser_type"`
	Timestamp string `json:"timestamp"`
}

type ActiveSessionResponseDTO struct {
	SessionId string `json:"session_id"`
	DeviceId string `json:"device_id"`
	BrowserType string `json:"browser_type"`
	LoggedInSince string `json:"logged_in_since"`
	ExpiresOn string `json:"expires_on"`
}

type AllActiveSessionsResponseDTO struct {
	Handle string `json:"handle"`
	ActiveSessions []*ActiveSessionResponseDTO `json:"active_sessions"`
}