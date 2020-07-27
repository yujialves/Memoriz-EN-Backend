package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"../secret"
)

type SubjectsResponse struct {
	Subjects []Subject `json:"subjects"`
}

type Subject struct {
	SubjectId      int       `json:"subjectId"`
	Name           string    `json:"name"`
	Grades         [13]Grade `json:"grades"`
	CorrectCount   int       `json:"correctCount"`
	InCorrectCount int       `json:"inCorrectCount"`
}

type Grade struct {
	Solvable int `json:"solvable"`
	All      int `json:"all"`
}

var SubjectsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// DB 接続
	db, err := sql.Open("mysql", secret.HOSTNAME+":"+secret.PASSWORD+"@tcp("+secret.HOSTNAME+")/"+secret.DBNAME)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// レスポンスとして返すsubject
	var response SubjectsResponse
	// 抽出したsubjectを一時的に格納する変数
	var tmpSubject SubjectsResponse

	// subjects の抽出
	rows, err := db.Query(`
	SELECT id, subject FROM subjects
	ORDER BY id
	;`)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		tmpSubject.Subjects = append(tmpSubject.Subjects, Subject{SubjectId: id, Name: name})
	}
	rows.Close()

	for _, subject := range tmpSubject.Subjects {
		// グレードの初期化
		var grades [13]Grade

		// 各グレードの Solvable の個数を抽出
		rows, err := db.Query(`
		SELECT G.grade, COUNT(Q.id) FROM questions AS Q 
		INNER JOIN grades AS G
		ON Q.grade_id = G.id 
		WHERE Q.subject_id = ?
		AND (
		G.grade = 0 
		OR (G.grade = 1 AND (G.last_updated < (NOW() - INTERVAL 1 DAY)))
		OR (G.grade = 2 AND (G.last_updated < (NOW() - INTERVAL 2 DAY)))
		OR (G.grade = 3 AND (G.last_updated < (NOW() - INTERVAL 4 DAY)))
		OR (G.grade = 4 AND (G.last_updated < (NOW() - INTERVAL 1 WEEK)))
		OR (G.grade = 5 AND (G.last_updated < (NOW() - INTERVAL 2 WEEK)))
		OR (G.grade = 6 AND (G.last_updated < (NOW() - INTERVAL 1 MONTH)))
		OR (G.grade = 7 AND (G.last_updated < (NOW() - INTERVAL 2 MONTH)))
		OR (G.grade = 8 AND (G.last_updated < (NOW() - INTERVAL 3 MONTH)))
		OR (G.grade = 9 AND (G.last_updated < (NOW() - INTERVAL 4 MONTH)))
		OR (G.grade = 10 AND (G.last_updated < (NOW() - INTERVAL 6 MONTH)))
		OR (G.grade = 11 AND (G.last_updated < (NOW() - INTERVAL 9 MONTH)))
		OR (G.grade = 12 AND (G.last_updated < (NOW() - INTERVAL 1 YEAR)))
		)
		GROUP BY G.grade
		ORDER BY G.grade
		;`, subject.SubjectId)
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var grade int
			var cnt int
			err := rows.Scan(&grade, &cnt)
			if err != nil {
				log.Fatal(err)
			}
			grades[grade].Solvable = cnt
		}
		rows.Close()

		// 各グレードの個数を抽出
		rows, err = db.Query(`
		SELECT G.grade, COUNT(Q.id) FROM questions AS Q
		INNER JOIN grades AS G
		ON Q.grade_id = G.id 
		WHERE Q.subject_id = ?
		GROUP BY G.grade
		ORDER BY G.grade
		;`, subject.SubjectId)
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			var grade int
			var cnt int
			err := rows.Scan(&grade, &cnt)
			if err != nil {
				log.Fatal(err)
			}
			grades[grade].All = cnt
		}
		rows.Close()

		// 各グレードの正解した数、不正解した数を抽出
		rows, err = db.Query(`
		SELECT G.correct_count, G.incorrect_count FROM questions AS Q
		INNER JOIN grades AS G
		ON Q.grade_id = G.id 
		WHERE Q.subject_id = ?
		;`, subject.SubjectId)
		if err != nil {
			log.Fatal(err)
		}

		var totalCorrectCount int
		var totalInCorrectCount int
		for rows.Next() {
			var correctCount int
			var inCorrectCount int
			err := rows.Scan(&correctCount, &inCorrectCount)
			if err != nil {
				log.Fatal(err)
			}
			totalCorrectCount += correctCount
			totalInCorrectCount += inCorrectCount
		}
		subject.CorrectCount = totalCorrectCount
		subject.InCorrectCount = totalInCorrectCount
		rows.Close()

		subject.Grades = grades
		response.Subjects = append(response.Subjects, subject)

	}

	json.NewEncoder(w).Encode(response)

})
