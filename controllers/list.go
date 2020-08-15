package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"../models"
	"../utils"
)

func (c QuestionController) QuestionListHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		expired, user := utils.CheckTokenDate(w, r)
		if expired {
			return
		}

		// POSTの解析
		var post models.QuestionListPost
		json.NewDecoder(r.Body).Decode(&post)

		// Questoin の情報を抽出
		stmt, err := db.Prepare(`
		SELECT Q.id, Q.question, Q.answer, G.grade, G.total_correct_count, G.total_incorrect_count
		FROM (((
		SELECT id, question, answer
		FROM questions
		WHERE subject_id = ?) AS Q
		INNER JOIN (
		SELECT question_id, user_id, grade, total_correct_count, total_incorrect_count
		FROM grades) AS G
		ON Q.id = G.question_id)
		INNER JOIN (
		SELECT id
		FROM users
		WHERE user = ?) AS U 
		ON G.user_id = U.id) 
		ORDER BY Q.id
		`)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		rows, err := stmt.Query(post.SubjectID, user)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// レスポンスに Question 情報を格納
		var response models.QuestionListResponse
		for rows.Next() {
			var question models.Question
			err := rows.Scan(&question.ID, &question.Question, &question.Answer, &question.Grade, &question.CorrectCountSum, &question.InCorrectCountSum)
			if err != nil {
				log.Fatal(err)
			}
			response.QuestionList = append(response.QuestionList, question)
		}
		defer rows.Close()

		utils.ResponseJSON(w, response)
	}
}
