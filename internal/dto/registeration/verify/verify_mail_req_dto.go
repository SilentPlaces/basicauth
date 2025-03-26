package verify_mail_req_dto

type VerifyMailReqDTO struct {
	Token string `json:"token"`
	Mail  string `json:"email"`
}
