package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"../models"
)

func (c Controller) QuestionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// POSTの解析
		var post models.QuestionPost
		json.NewDecoder(r.Body).Decode(&post)

		// subjectIdの取得
		subjectID := post.SubjectID

		// Subject の Solvable な問題をすべて抽出

		stmt, err := db.Prepare(`
		SELECT Q.id, Q.question, Q.answer, G.grade FROM questions AS Q
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
		;`)
		if err != nil {
			log.Fatal(err)
		}

		rows, err := stmt.Query(subjectID)
		if err != nil {
			log.Fatal(err)
		}

		var questions []models.Question

		for rows.Next() {
			var id int
			var question string
			var answer string
			var grade int
			err := rows.Scan(&id, &question, &answer, &grade)
			if err != nil {
				log.Fatal(err)
			}
			questions = append(questions, models.Question{ID: id, Question: question, Answer: answer, Grade: grade})
		}
		rows.Close()

		// ランダムな Question を抽出
		var response models.QuestionResponse
		response.Rest = len(questions)
		if len(questions) == 0 {
			// レスポンスの返信
			json.NewEncoder(w).Encode(response)
		} else {
			rand.Seed(time.Now().Unix())
			log.Printf("%v, %T", len(questions), len(questions))
			response.Question = questions[rand.Intn(len(questions))]
			// レスポンスの返信
			json.NewEncoder(w).Encode(response)
		}
	}
}
