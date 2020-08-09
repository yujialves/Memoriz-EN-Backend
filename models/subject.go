package models

type Subject struct {
	SubjectId      int       `json:"subjectId"`
	Name           string    `json:"name"`
	Grades         [13]Grade `json:"grades"`
	CorrectCount   int       `json:"correctCount"`
	InCorrectCount int       `json:"inCorrectCount"`
}
