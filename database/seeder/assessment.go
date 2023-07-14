package seeder

import (
	"engine-db/entity"
)

func (s *seeder) SeedAssesmentTask(session entity.Session, questions []entity.Question, reviewe, reviewer entity.User) error {
	for _, question := range questions {
		answer := entity.QuestionAnswer{
			SessionID:  session.ID,
			RevieweeID: reviewe.ID,
			ReviewerID: reviewer.ID,
			QuestionID: question.ID,
			Status:     "todo",
		}
		err := s.db.Create(&answer).Error
		if err != nil {
			logs.Error().Err(err).Msg("failed to seed question answer")
			return err
		}
	}

	return nil
}

func (s *seeder) SeedAnswerAssesmentTask(session entity.Session, questions []entity.Question, reviewe, reviewer entity.User) error {
	for _, question := range questions {
		switch reviewer.Role {
		case "leader", "employee", "admin":

			err := s.db.Model(&entity.QuestionAnswer{}).
				Where("reviewee_id = ?", reviewe.ID).
				Where("reviewer_id = ?", reviewer.ID).
				Where("question_id = ?", question.ID).
				Updates(map[string]interface{}{
					"status": "reviewed",
					"scale":  s.faker.RandomIntElement([]int{1, 2, 3, 4}),
				}).Error
			if err != nil {
				logs.Error().Err(err).Msg("failed to seed question answer")
				return err
			}
		}

	}

	return nil
}
