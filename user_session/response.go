package user_session

type SessionResponseDTO struct {
	Id string `json:"id"`
	Handle string `json:"handle"`
	SessionToken string `json:"session_token"`
	ActiveTill int64 `json:"active_till"`
}

type RegisterDeviceResponseDTO struct {
	Id string `json:"id"`
}

func NewSessionResponseDTO(id string, handle string, sessionToken string, activeTill int64) *SessionResponseDTO {
	return &SessionResponseDTO{Id: id, Handle: handle, SessionToken: sessionToken, ActiveTill: activeTill}
}
