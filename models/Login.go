package models


type LoginResponse struct {
	Message string
	Status int32
	Token string
	// Data *ReqBody
}

func GetLoginResponse (msg string, status int32, token string) *LoginResponse {
	return &LoginResponse {
		msg ,status, token,
	}
}