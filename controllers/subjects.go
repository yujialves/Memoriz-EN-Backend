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
		SELECT subjects.id, subjects.subject, MAIN.grade1, MAIN.grade2, MAIN.grade3, MAIN.grade4, MAIN.grade5, MAIN.grade6, MAIN.grade7, MAIN.grade8, MAIN.grade9, MAIN.grade10, MAIN.grade11, MAIN.grade12 FROM subjects LEFT OUTER JOIN (
		SELECT id, subject, SUM(grade0) AS grade0, SUM(grade1) AS grade1, SUM(grade2) AS grade2, SUM(grade3) AS grade3, SUM(grade4) AS grade4, SUM(grade5) AS grade5, SUM(grade6) AS grade6, SUM(grade7) AS grade7, SUM(grade8) AS grade8, SUM(grade9) AS grade9, SUM(grade10) AS grade10, SUM(grade11) AS grade11, SUM(grade12) AS grade12
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
		) AS MAIN USING(id)
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
			var grade0, grade1, grade2, grade3, grade4, grade5, grade6, grade7, grade8, grade9, grade10, grade11, grade12 sql.NullInt64
			err = rows.Scan(&subject.SubjectId, &subject.Name, &grade0, &grade1, &grade2, &grade3, &grade4, &grade5, &grade6, &grade7, &grade8, &grade9, &grade10, &grade11, &grade12)
			if err != nil {
				log.Fatal(err)
			}
			if grade0.Valid {
				subject.Grades[0].Solvable = grade0.Int64
				subject.Grades[1].Solvable = grade1.Int64
				subject.Grades[2].Solvable = grade2.Int64
				subject.Grades[3].Solvable = grade3.Int64
				subject.Grades[4].Solvable = grade4.Int64
				subject.Grades[5].Solvable = grade5.Int64
				subject.Grades[6].Solvable = grade6.Int64
				subject.Grades[7].Solvable = grade7.Int64
				subject.Grades[8].Solvable = grade8.Int64
				subject.Grades[9].Solvable = grade9.Int64
				subject.Grades[10].Solvable = grade10.Int64
				subject.Grades[11].Solvable = grade11.Int64
				subject.Grades[12].Solvable = grade12.Int64
			}
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
