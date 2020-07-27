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

type CorrectPost struct {
	QuestionID int `json:"questionId"`
	SubjectID  int `json:"subjectId"`
}

var CorrectHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	// POSTの解析
	var post CorrectPost
	json.NewDecoder(r.Body).Decode(&post)

	// questionID, subjectIDの取得
	questionID := post.QuestionID
	subjectID := post.SubjectID

	// DB 接続
	db, err := sql.Open("mysql", secret.HOSTNAME+":"+secret.PASSWORD+"@tcp("+secret.HOSTNAME+")/"+secret.DBNAME)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// question のグレードをアップ
	stmt, err := db.Prepare(`
	UPDATE questions
	SET
		grade = grade + 1,
		last_updated = NOW()
	WHERE 
		id = ?
	;`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(questionID)
	if err != nil {
		log.Fatal(err)
	}

	// Subject の Solvable な問題をすべて抽出

	stmt, err = db.Prepare(`
	SELECT id, question, answer, grade FROM questions 
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
		var id int
		var question string
		var answer string
		var grade int
		err := rows.Scan(&id, &question, &answer, &grade)
		if err != nil {
			log.Fatal(err)
		}
		questions = append(questions, Question{ID: id, Question: question, Answer: answer, Grade: grade})
	}
	rows.Close()

	// ランダムな Question を抽出
	var response QuestionResponse
	response.Rest = len(questions)
	if len(questions) == 0 {
		// レスポンスの返信
		json.NewEncoder(w).Encode(response)
	} else {
		rand.Seed(time.Now().Unix())
		response.Question = questions[rand.Intn(len(questions))]
		// レスポンスの返信
		json.NewEncoder(w).Encode(response)
	}
})
