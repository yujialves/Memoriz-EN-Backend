package models

type LoginPost struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type CorrectPost struct {
	QuestionID int `json:"questionId"`
	SubjectID  int `json:"subjectId"`
}

type InCorrectPost struct {
	QuestionID int `json:"questionId"`
	SubjectID  int `json:"subjectId"`
}

type QuestionPost struct {
	SubjectID int `json:"subjectId"`
}

type QuestionListPost struct {
	SubjectID int `json:"subjectId"`
}

type BingPost struct {
	Word string `json:"word"`
	ID   int    `json:"id"`
}
