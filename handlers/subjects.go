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
	Subject_id int       `json:"subject_id"`
	Name       string    `json:"name"`
	Grades     [13]Grade `json:"grades"`
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

	var response SubjectsResponse

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
		response.Subjects = append(response.Subjects, Subject{Subject_id: id, Name: name})
	}
	rows.Close()

	for _, subject := range response.Subjects {
		// グレードの初期化
		var grades [13]Grade

		log.Printf("name: %v, id: %v ", subject.Name, subject.Subject_id)

		// 各グレードの Solvable の個数を抽出
		rows, err := db.Query(`
		SELECT grade, COUNT(id) FROM questions 
		WHERE subject_id = ?
		AND (
		grade = 0 
		OR (grade = 1 AND (last_updated < (NOW() - INTERVAL 1 DAY)))
		OR (grade = 2 AND (last_updated < (NOW() - INTERVAL 2 DAY)))
		OR (grade = 3 AND (last_updated < (NOW() - INTERVAL 4 DAY)))
		OR (grade = 4 AND (last_updated < (NOW() - INTERVAL 1 WEEK)))
		OR (grade = 5 AND (last_updated < (NOW() - INTERVAL 2 WEEK)))
		OR (grade = 6 AND (last_updated < (NOW() - INTERVAL 1 MONTH)))
		OR (grade = 7 AND (last_updated < (NOW() - INTERVAL 2 MONTH)))
		OR (grade = 8 AND (last_updated < (NOW() - INTERVAL 3 MONTH)))
		OR (grade = 9 AND (last_updated < (NOW() - INTERVAL 4 MONTH)))
		OR (grade = 10 AND (last_updated < (NOW() - INTERVAL 6 MONTH)))
		OR (grade = 11 AND (last_updated < (NOW() - INTERVAL 9 MONTH)))
		OR (grade = 12 AND (last_updated < (NOW() - INTERVAL 1 YEAR)))
		)
		GROUP BY grade
		ORDER BY grade
		;`, subject.Subject_id)
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
			log.Printf("sgrade%d: %d", grade, cnt)
			grades[grade].Solvable = cnt
		}
		rows.Close()

		// 各グレードの個数を抽出
		rows, err = db.Query(`
		SELECT grade, COUNT(id) FROM questions
		WHERE subject_id = ?
		GROUP BY grade
		ORDER BY grade
		;`, subject.Subject_id)
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
			log.Printf("agrade%d: %d", grade, cnt)
			grades[grade].All = cnt
		}
		rows.Close()

		subject.Grades = grades

	}

	json.NewEncoder(w).Encode(response)

})
