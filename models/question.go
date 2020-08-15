package models

type Question struct {
	ID                int    `json:"id"`
	Question          string `json:"question"`
	Answer            string `json:"answer"`
	Grade             int    `json:"grade"`
	CorrectCountSum   int    `json:"correctCountSum"`
	InCorrectCountSum int    `json:"inCorrectCountSum"`
}
