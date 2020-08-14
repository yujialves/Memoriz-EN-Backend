package models

type LoginResponse struct {
	User         string `json:"user"`
	Password     string `json:"password"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	ExpireDate   int64  `json:"expireDate"`
}

type QuestionResponse struct {
	Question Question `json:"question"`
	Rest     int      `json:"rest"`
}

type SubjectsResponse struct {
	Subjects []Subject `json:"subjects"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
