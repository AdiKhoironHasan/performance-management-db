package seeder

import (
	"engine-db/entity"
	"fmt"
	"time"
)

func (s *seeder) SeedSessionPeerAssesment(eventID int) (entity.Session, error) {
	session := entity.Session{
		EventID: eventID,
		Name:    "Peer Assesment",
	}

	err := s.db.Create(&session).Error
	if err != nil {
		logs.Error().Err(err).Msg("failed to seed session")
		return session, err
	}

	return session, nil
}

func (s *seeder) SeedSession(eventID int, sessionType string) (entity.Session, error) {
	var sessionName string

	switch sessionType {
	case "self_assess":
		sessionName = "Self Assesment"
	case "peers_assess":
		sessionName = "Peer Assesment"
	case "member_assess":
		sessionName = "Member Assesment"
	default:
		return entity.Session{}, fmt.Errorf("invalid session type")
	}

	session := entity.Session{
		EventID:     eventID,
		Name:        sessionName,
		Type:        sessionType,
		Description: fmt.Sprintf("This is %s session", sessionName),
		StartDate:   time.Now(),
		EndDate:     time.Now(),
	}

	err := s.db.Create(&session).Error
	if err != nil {
		logs.Error().Err(err).Msg("failed to seed session")
		return session, err
	}

	return session, nil
}
