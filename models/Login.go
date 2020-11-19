package models


type loginResponse struct {
	Message string
	Status int32
	Token string
	LoginRowAffected int64
	// Data *ReqBody
}

func GetLoginResponse (msg string, status int32, token string, rowAffected int64) *loginResponse {
	return &loginResponse {
		msg ,status, token, rowAffected,
	}
}

