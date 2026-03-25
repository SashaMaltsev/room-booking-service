package dto

type DummyLoginRequest struct {
	Role string `json:"role"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
