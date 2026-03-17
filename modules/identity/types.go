package identity

type submitLoginRequest struct {
	Identifier     string `json:"identifier"`
	Password       string `json:"password"`
	LoginChallenge string `json:"login_challenge"`
}

type submitLoginResponse struct {
	RedirectTo string `json:"redirect_to"`
}
