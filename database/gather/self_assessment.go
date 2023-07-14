package gather

import (
	"engine-db/entity"
)

func (g *gather) GetSelfAssessmentTask(session entity.Session) ([]entity.ResultSelfAssessment, error) {
	results := []entity.ResultSelfAssessment{}

	rows, err := g.db.Model(&entity.User{}).
		Select([]string{
			"users.name AS reviewee",
			"questions.name AS question",
			"SUM(question_answers.scale) AS total",
		}).
		Joins("JOIN question_answers ON question_answers.reviewee_id = users.id AND question_answers.deleted_at IS NULL").
		Joins("JOIN questions ON questions.id = question_answers.question_id AND questions.deleted_at IS NULL").
		Where("users.deleted_at IS NULL").
		Where("questions.type = $1", "scale").
		Where("question_answers.session_id = $2", session.ID).
		Group("users.name, questions.name").
		Rows()

	if err != nil {
		logs.Err(err).Msg("failed to get self assessment")
		return results, err
	}

	for rows.Next() {
		result := entity.ResultSelfAssessment{}
		err = rows.Scan(
			&result.Reviewee,
			&result.Question,
			&result.Total,
		)
		if err != nil {
			logs.Err(err).Msg("failed to scan self assessment")
			return results, err
		}

		results = append(results, result)
	}

	return results, nil
}
