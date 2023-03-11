package dto

type SignupRequest struct {
	Email     string `json:"email,omitempty" binding:"required"`
	Password  string `json:"password,omitempty" binding:"required"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenDetails struct {
	Token     *string
	TokenUuid string
	UserId    int
	ExpiresIn *int64
}
