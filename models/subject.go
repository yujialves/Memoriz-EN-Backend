package models

type Subject struct {
	SubjectId           int64     `json:"subjectId"`
	Name                string    `json:"name"`
	Grades              [13]Grade `json:"grades"`
	CorrectCount        int       `json:"correctCount"`
	InCorrectCount      int       `json:"inCorrectCount"`
	TotalCorrectCount   int       `json:"totalCorrectCount"`
	TotalInCorrectCount int       `json:"totalInCorrectCount"`
}
