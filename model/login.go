package model

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
