package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"../models"
	"../utils"
)

type QuestionController struct{}

func (c QuestionController) QuestionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		expired, userID := utils.CheckTokenDate(w, r)
		if expired {
			return
		}

		// POSTの解析
		var post models.QuestionPost
		json.NewDecoder(r.Body).Decode(&post)

		// subjectIdの取得
		subjectID := post.SubjectID

		// Subject の Solvable な問題をすべて抽出

		stmt, err := db.Prepare(`
		SELECT Q.id, Q.question, Q.answer, G.grade
		FROM (((
		SELECT id, question, answer
		FROM questions
		WHERE subject_id = ?) AS Q
		INNER JOIN (
		SELECT grade, user_id, question_id
		FROM grades
		WHERE grade = 0 
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
		) AS G
		ON Q.id = G.question_id)
		INNER JOIN (
		SELECT id
		FROM users
		WHERE user = ?) AS U 
		ON G.user_id = U.id)
		;`)
		if err != nil {
			log.Fatal(err)
		}

		rows, err := stmt.Query(subjectID, userID)
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
			utils.ResponseJSON(w, response)
		} else {
			rand.Seed(time.Now().Unix())
			log.Printf("%v, %T", len(questions), len(questions))
			response.Question = questions[rand.Intn(len(questions))]
			// レスポンスの返信
			utils.ResponseJSON(w, response)
		}
	}
}
