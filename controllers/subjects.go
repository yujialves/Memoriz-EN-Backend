package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"../models"
	"../utils"
)

type SubjectsController struct{}

func (c SubjectsController) SubjectsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		expired, user := utils.CheckTokenDate(w, r)
		if expired {
			return
		}

		// 日付が変わっていたら正解数、不正解数をリセット
		stmt, err := db.Prepare(`
		UPDATE grades AS G
		INNER JOIN users AS U
		ON G.user_id = U.id
		SET 
			G.correct_count = 0,
			G.incorrect_count = 0
		WHERE G.last_updated < CURDATE()
		AND U.user = ?
		;`)
		_, err = stmt.Exec(user)
		if err != nil {
			log.Fatal(err)
		}

		// レスポンスとして返すsubject
		var response models.SubjectsResponse
		// 抽出したsubjectを一時的に格納する変数
		var tmpSubject models.SubjectsResponse

		// subjects の抽出
		rows, err := db.Query(`
		SELECT id, subject FROM subjects
		ORDER BY subject
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
			tmpSubject.Subjects = append(tmpSubject.Subjects, models.Subject{SubjectId: id, Name: name})
		}
		rows.Close()

		for _, subject := range tmpSubject.Subjects {
			// グレードの初期化
			var grades [13]models.Grade

			// 各グレードの Solvable の個数を抽出
			rows, err := db.Query(`
			SELECT 
				G.grade, COUNT(Q.id)
			FROM (
				SELECT id
				FROM questions
				WHERE subject_id = ?
			) AS Q INNER JOIN (
				SELECT user_id, question_id, grade
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
			) AS G INNER JOIN (
				SELECT id
				FROM users
				WHERE user = ?
			) AS U ON Q.id = G.question_id AND G.user_id = U.id
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
			SELECT G.grade, COUNT(Q.id)
			FROM (
				SELECT id
				FROM questions
				WHERE subject_id = ?
			) AS Q INNER JOIN (
				SELECT grade, question_id, user_id
				FROM grades
			) AS G INNER JOIN (
				SELECT id
				FROM users
				WHERE user = ?
			) AS U ON Q.id = G.question_id AND G.user_id = U.id
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
			SELECT G.correct_count, G.incorrect_count
			FROM (
				SELECT id
				FROM questions
				WHERE subject_id = ?
			) AS Q INNER JOIN (
				SELECT question_id, user_id, correct_count, incorrect_count
				FROM grades
			) AS G INNER JOIN (
				SELECT id
				FROM users
				WHERE user = ?
			) AS U ON Q.id = G.question_id AND G.user_id = U.id
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

		utils.ResponseJSON(w, response)
	}
}
