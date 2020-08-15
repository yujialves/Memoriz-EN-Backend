package models

type LoginResponse struct {
	User         string `json:"user"`
	Password     string `json:"password"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	ExpireDate   int64  `json:"expireDate"`
}

type RefreshResponse struct {
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
	Exp      int       `json:"exp"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type QuestionListResponse struct {
	QuestionList []Question `json:"questionList"`
}
