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

func (c QuestionController) InCorrectHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		expired, userID := utils.CheckTokenDate(w, r)
		if expired {
			return
		}

		// POSTの解析
		var post models.InCorrectPost
		json.NewDecoder(r.Body).Decode(&post)

		// questionID, subjectIDの取得
		questionID := post.QuestionID
		subjectID := post.SubjectID

		// question のグレードをダウン
		stmt, err := db.Prepare(`
		UPDATE questions AS Q
		INNER JOIN grades AS G
		ON Q.id = G.question_id
		INNER JOIN users AS U
		ON G.user_id = U.id
		SET G.grade = 0,
		G.last_updated = NOW(),
		G.incorrect_count = G.incorrect_count + 1,
		G.total_incorrect_count = G.total_incorrect_count + 1
		WHERE Q.id = ? AND U.user = ?
		;`)
		if err != nil {
			log.Fatal(err)
		}

		_, err = stmt.Exec(questionID, userID)
		if err != nil {
			log.Fatal(err)
		}

		// Subject の Solvable な問題をすべて抽出

		stmt, err = db.Prepare(`
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
			response.Question = questions[rand.Intn(len(questions))]
			// レスポンスの返信
			utils.ResponseJSON(w, response)
		}
	}
}
