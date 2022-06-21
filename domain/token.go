package domain

// Token struct define
type Token struct {
	UserId       string `json:"userid" form:"userid" query:"userid"`
	Token        string `json:"token" form:"token" query:"token"`
	RefreshToken string `json:"refresh_token" form:"refresh_token" query:"refresh_token"`
}
