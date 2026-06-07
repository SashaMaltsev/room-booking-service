package dto

type DummyLoginRequest struct {
	Role string `json:"role"`
	User string `json:"user"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
