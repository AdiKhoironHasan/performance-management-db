package gather

import (
	"engine-db/entity"
	"fmt"
)

func (g *gather) GetSelfAssessment() error {
	type Result struct {
		Reviewee string
		Question string
		Total    float64
	}

	rows, err := g.db.Model(&entity.User{}).
		Select([]string{
			"users.name AS reviewee",
			"questions.name AS question",
			"SUM(question_answers.scale) AS total",
		}).
		Joins("JOIN question_answers ON question_answers.reviewee_id = users.id AND question_answers.deleted_at IS NULL").
		Joins("JOIN questions ON questions.id = question_answers.question_id AND questions.deleted_at IS NULL").
		Where("users.deleted_at IS NULL").
		Group("users.name, questions.name").
		Rows()

	if err != nil {
		logs.Err(err).Msg("failed to get self assessment")
		return err
	}

	results := []Result{}
	for rows.Next() {
		result := Result{}
		err = rows.Scan(
			&result.Reviewee,
			&result.Question,
			&result.Total,
		)
		if err != nil {
			logs.Err(err).Msg("failed to scan self assessment")
			return err
		}

		results = append(results, result)
	}

	for _, result := range results {
		fmt.Println(result)
	}

	return nil
}
