package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"../secret"
)

type QuestionPost struct {
	SubjectID string `json:"subjectId"`
}

type QuestionResponse struct {
	Question Question `json:"question"`
}

type Question struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Grade    int    `json:"grade"`
}

var QuestionHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// POSTの解析
	var post QuestionPost
	log.Println("body: ", r.Body)
	json.NewDecoder(r.Body).Decode(&post)

	// subjectIdの取得
	subjectID := post.SubjectID
	log.Println("subjectID: ", subjectID)

	// DB 接続
	db, err := sql.Open("mysql", secret.HOSTNAME+":"+secret.PASSWORD+"@tcp("+secret.HOSTNAME+")/"+secret.DBNAME)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Subject の Solvable な問題をすべて抽出

	stmt, err := db.Prepare(`
	SELECT question, answer, grade FROM questions 
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
	;`)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := stmt.Query(subjectID)
	if err != nil {
		log.Fatal(err)
	}

	var questions []Question

	for rows.Next() {
		var question string
		var answer string
		var grade int
		err := rows.Scan(&question, &answer, &grade)
		if err != nil {
			log.Fatal(err)
		}
		questions = append(questions, Question{Question: question, Answer: answer, Grade: grade})
	}
	rows.Close()

	// ランダムな Question を抽出
	var response QuestionResponse
	rand.Seed(time.Now().Unix())
	log.Println("len: ", len(questions))
	log.Printf("%v, %T", len(questions), len(questions))
	response.Question = questions[rand.Intn(len(questions))]

	// レスポンスの返信
	json.NewEncoder(w).Encode(response)
})
