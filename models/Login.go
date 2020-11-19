package models


type loginResponse struct {
	Message string
	Status int32
	Token string
	// Data *ReqBody
}

func GetLoginResponse (msg string, status int32, token string) *loginResponse {
	return &loginResponse {
		msg ,status, token,
	}
}

