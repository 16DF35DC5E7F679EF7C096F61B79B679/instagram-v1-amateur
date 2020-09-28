package user_session

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

func generateToken(handle, deviceId, browserType string) string {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["handle"] = handle
	atClaims["deviceId"] = deviceId
	atClaims["browserType"] = browserType
	atClaims["timestamp"] = time.Now().Unix()
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return ""
	}
	return token
}
