package models

type QuestionResponse struct {
	Question Question `json:"question"`
	Rest     int      `json:"rest"`
}

type SubjectsResponse struct {
	Subjects []Subject `json:"subjects"`
}
