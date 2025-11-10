package dto

type TokenReq struct {
	AccessToken string `json:"access_token" binding:"required"`
}
