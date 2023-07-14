package seeder

// import (
// 	"engine-db/entity"
// 	"fmt"
// 	"time"
// )

// func (s *seeder) SeedPeersAssesmentTask(session entity.Session, questions []entity.Question, reviewer, peer entity.User) error {
// 	for _, question := range questions {
// 		answer := entity.QuestionAnswer{
// 			SessionID:  session.ID,
// 			RevieweeID: peer.ID,
// 			ReviewerID: reviewer.ID,
// 			QuestionID: question.ID,
// 			Status:     "todo",
// 		}
// 		err := s.db.Create(&answer).Error
// 		if err != nil {
// 			logs.Error().Err(err).Msg("failed to seed question answer")
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (s *seeder) SeedPeersAssesment(eventID int) (entity.Session, error) {
// 	sessions := entity.Session{
// 		EventID:     eventID,
// 		Type:        "peers_assess",
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

// func (s *seeder) SeedAnswerPeersAssesmentTask(questions []entity.Question, participants []entity.User) error {
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
