package seeder

// import (
// 	"engine-db/entity"
// 	"fmt"
// 	"time"
// )

// func (s *seeder) SeedSelfAssesmentTask(session entity.Session, participants []entity.User) ([]entity.Question, error) {
// 	questions, err := s.SeedQuestion(10, session)
// 	if err != nil {
// 		logs.Error().Err(err).Msg("failed to seed question")
// 		return nil, err
// 	}

// 	for _, p := range participants {
// 		for _, question := range questions {

// 			switch p.Role {
// 			case "leader", "employee", "admin":
// 				answer := entity.QuestionAnswer{
// 					SessionID:  1,
// 					RevieweeID: p.ID,
// 					ReviewerID: p.ID,
// 					QuestionID: question.ID,
// 					Status:     "todo",
// 				}
// 				err := s.db.Create(&answer).Error
// 				if err != nil {
// 					logs.Error().Err(err).Msg("failed to seed question answer")
// 					return questions, err
// 				}
// 			}

// 		}
// 	}

// 	return questions, nil
// }

// func (s *seeder) SeedAnswerSelfAssesmentTask(questions []entity.Question, participants []entity.User) error {
// 	for _, p := range participants {
// 		for _, question := range questions {

// 			switch p.Role {
// 			case "leader", "employee", "admin":

// 				err := s.db.Model(&entity.QuestionAnswer{}).
// 					Where("reviewee_id = ?", p.ID).
// 					Where("question_id = ?", question.ID).
// 					Updates(map[string]interface{}{
// 						"status": "reviewed",
// 						"scale":  s.faker.RandomIntElement([]int{1, 2, 3, 4}),
// 					}).Error
// 				if err != nil {
// 					logs.Error().Err(err).Msg("failed to seed question answer")
// 					return err
// 				}
// 			}

// 		}
// 	}

// 	return nil
// }

// func (s *seeder) SeedSessionSelfAssesment(eventID int) (entity.Session, error) {
// 	sessions := entity.Session{
// 		EventID:     eventID,
// 		Type:        "self_assess",
// 		StartDate:   time.Now().AddDate(0, 0, 1),
// 		EndDate:     time.Now().AddDate(0, 0, 1),
// 		Name:        fmt.Sprintf("Session for event %v", eventID),
// 		Description: fmt.Sprintf("Description Session for event %v", eventID),
// 	}

// 	err := s.db.Create(&sessions).Error
// 	if err != nil {
// 		logs.Error().Err(err).Msg("failed to seed sessions")
// 		return sessions, err
// 	}

// 	return sessions, nil
// }
