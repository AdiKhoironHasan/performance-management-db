package seeder

import (
	"engine-db/entity"
	"fmt"
)

func (s *seeder) SeedQuestion(count int, session entity.Session) ([]entity.Question, error) {
	questions := []entity.Question{}

	sort := 1
	for i := 1; i <= count; i++ {
		qusetionType := "scale"

		switch i {
		case 1:
			qusetionType = "radio"
		case 2:
			qusetionType = "dropdown"
		case 3:
			qusetionType = "essay"
		}

		questions = append(questions, entity.Question{
			SessionID: session.ID,
			Sort:      sort,
			Name:      fmt.Sprintf("Question %v", i),
			Type:      qusetionType,
			Option:    "",
			Required:  s.faker.RandomStringElement([]string{"true", "false"}),
		})

		sort++
	}

	err := s.db.Create(&questions).Error
	if err != nil {
		logs.Error().Err(err).Msg("failed to seed questions")
	}

	return questions, nil
}
