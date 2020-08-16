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
		SET G.correct_count = 0, G.incorrect_count = 0
		WHERE G.last_updated < CURDATE()
		AND U.user = ?
		;`)
		_, err = stmt.Exec(user)
		if err != nil {
			log.Fatal(err)
		}
		stmt.Close()

		// レスポンスとして返すsubject
		var response models.SubjectsResponse

		// 各 Subject の解ける問題の数を抽出
		stmt, err = db.Prepare(`
		SELECT id, subject, SUM(grade0), SUM(grade1), SUM(grade2), SUM(grade3), SUM(grade4), SUM(grade5), SUM(grade6), SUM(grade7), SUM(grade8), SUM(grade9), SUM(grade10), SUM(grade11), SUM(grade12)
		FROM (
		SELECT S.id, S.subject,
		CASE G.grade WHEN 0 THEN COUNT(Q.id) ELSE 0 END AS grade0,
		CASE G.grade WHEN 1 THEN COUNT(Q.id) ELSE 0 END AS grade1,
		CASE G.grade WHEN 2 THEN COUNT(Q.id) ELSE 0 END AS grade2,
		CASE G.grade WHEN 3 THEN COUNT(Q.id) ELSE 0 END AS grade3,
		CASE G.grade WHEN 4 THEN COUNT(Q.id) ELSE 0 END AS grade4,
		CASE G.grade WHEN 5 THEN COUNT(Q.id) ELSE 0 END AS grade5,
		CASE G.grade WHEN 6 THEN COUNT(Q.id) ELSE 0 END AS grade6,
		CASE G.grade WHEN 7 THEN COUNT(Q.id) ELSE 0 END AS grade7,
		CASE G.grade WHEN 8 THEN COUNT(Q.id) ELSE 0 END AS grade8,
		CASE G.grade WHEN 9 THEN COUNT(Q.id) ELSE 0 END AS grade9,
		CASE G.grade WHEN 10 THEN COUNT(Q.id) ELSE 0 END AS grade10,
		CASE G.grade WHEN 11 THEN COUNT(Q.id) ELSE 0 END AS grade11,
		CASE G.grade WHEN 12 THEN COUNT(Q.id) ELSE 0 END AS grade12
		FROM ((((
		SELECT id, subject
		FROM subjects) AS S
		INNER JOIN (
		SELECT id, subject_id
		FROM questions) AS Q ON S.id = Q.subject_id)
		INNER JOIN (
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
		) AS G ON Q.id = G.question_id)
		INNER JOIN (
		SELECT id
		FROM users
		WHERE user = ?) AS U ON G.user_id = U.id)
		GROUP BY S.id, G.grade
		ORDER BY S.id, G.grade) AS N
		GROUP BY subject
		ORDER BY subject
		;`)
		if err != nil {
			log.Fatal()
		}
		rows, err := stmt.Query(user)
		if err != nil {
			log.Fatal()
		}

		for rows.Next() {
			var subject models.Subject
			err = rows.Scan(&subject.SubjectId, &subject.Name, &subject.Grades[0].Solvable, &subject.Grades[1].Solvable, &subject.Grades[2].Solvable, &subject.Grades[3].Solvable, &subject.Grades[4].Solvable, &subject.Grades[5].Solvable, &subject.Grades[6].Solvable, &subject.Grades[7].Solvable, &subject.Grades[8].Solvable, &subject.Grades[9].Solvable, &subject.Grades[10].Solvable, &subject.Grades[11].Solvable, &subject.Grades[12].Solvable)
			response.Subjects = append(response.Subjects, subject)
		}
		rows.Close()
		stmt.Close()

		// 各 Subject の問題の数を抽出
		stmt, err = db.Prepare(`
		SELECT SUM(grade0), SUM(grade1), SUM(grade2), SUM(grade3), SUM(grade4), SUM(grade5), SUM(grade6), SUM(grade7), SUM(grade8), SUM(grade9), SUM(grade10), SUM(grade11), SUM(grade12)
		FROM (
		SELECT S.id, S.subject,
		CASE G.grade WHEN 0 THEN COUNT(Q.id) ELSE 0 END AS grade0,
		CASE G.grade WHEN 1 THEN COUNT(Q.id) ELSE 0 END AS grade1,
		CASE G.grade WHEN 2 THEN COUNT(Q.id) ELSE 0 END AS grade2,
		CASE G.grade WHEN 3 THEN COUNT(Q.id) ELSE 0 END AS grade3,
		CASE G.grade WHEN 4 THEN COUNT(Q.id) ELSE 0 END AS grade4,
		CASE G.grade WHEN 5 THEN COUNT(Q.id) ELSE 0 END AS grade5,
		CASE G.grade WHEN 6 THEN COUNT(Q.id) ELSE 0 END AS grade6,
		CASE G.grade WHEN 7 THEN COUNT(Q.id) ELSE 0 END AS grade7,
		CASE G.grade WHEN 8 THEN COUNT(Q.id) ELSE 0 END AS grade8,
		CASE G.grade WHEN 9 THEN COUNT(Q.id) ELSE 0 END AS grade9,
		CASE G.grade WHEN 10 THEN COUNT(Q.id) ELSE 0 END AS grade10,
		CASE G.grade WHEN 11 THEN COUNT(Q.id) ELSE 0 END AS grade11,
		CASE G.grade WHEN 12 THEN COUNT(Q.id) ELSE 0 END AS grade12
		FROM ((((
		SELECT id, subject
		FROM subjects) AS S
		INNER JOIN (
		SELECT id, subject_id
		FROM questions) AS Q ON S.id = Q.subject_id)
		INNER JOIN (
		SELECT user_id, question_id, grade
		FROM grades) AS G ON Q.id = G.question_id)
		INNER JOIN (
		SELECT id
		FROM users
		WHERE user = ?) AS U ON G.user_id = U.id)
		GROUP BY S.id, G.grade
		ORDER BY S.id, G.grade) AS N
		GROUP BY subject
		ORDER BY subject
		;`)
		if err != nil {
			log.Fatal()
		}
		rows, err = stmt.Query(user)
		if err != nil {
			log.Fatal()
		}

		var i int
		for rows.Next() {
			err = rows.Scan(&response.Subjects[i].Grades[0].All, &response.Subjects[i].Grades[1].All, &response.Subjects[i].Grades[2].All, &response.Subjects[i].Grades[3].All, &response.Subjects[i].Grades[4].All, &response.Subjects[i].Grades[5].All, &response.Subjects[i].Grades[6].All, &response.Subjects[i].Grades[7].All, &response.Subjects[i].Grades[8].All, &response.Subjects[i].Grades[9].All, &response.Subjects[i].Grades[10].All, &response.Subjects[i].Grades[11].All, &response.Subjects[i].Grades[12].All)
			i++
		}
		rows.Close()
		stmt.Close()

		// 各 Subject の正解した数、不正解した数などを抽出
		stmt, err = db.Prepare(`
		SELECT SUM(G.correct_count), SUM(G.incorrect_count), SUM(G.total_correct_count), SUM(G.total_incorrect_count)
		FROM ((((
		SELECT id, subject
		FROM subjects) AS S
		INNER JOIN (
		SELECT id, subject_id
		FROM questions) AS Q ON S.id = Q.subject_id)
		INNER JOIN (
		SELECT question_id, user_id, correct_count, incorrect_count, total_correct_count, total_incorrect_count
		FROM grades) AS G ON Q.id = G.question_id)
		INNER JOIN (
		SELECT id
		FROM users
		WHERE user = ?) AS U ON G.user_id = U.id)
		GROUP BY subject
		ORDER BY subject
		;`)
		if err != nil {
			log.Fatal()
		}
		rows, err = stmt.Query(user)
		if err != nil {
			log.Fatal()
		}

		i = 0
		for rows.Next() {
			err = rows.Scan(&response.Subjects[i].CorrectCount, &response.Subjects[i].InCorrectCount, &response.Subjects[i].TotalCorrectCount, &response.Subjects[i].TotalInCorrectCount)
			i++
		}
		rows.Close()
		stmt.Close()

		response.Exp = utils.CalculateExp(response.Subjects)
		utils.ResponseJSON(w, response)
	}
}
