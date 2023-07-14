package seeder

import (
	"engine-db/entity"
	"fmt"
)

func (s *seeder) SeedOneEvent() (entity.Event, error) {
	event := entity.Event{
		Name:        fmt.Sprintf("Event %v", 1),
		Description: fmt.Sprintf("Description for Event %v", 1),
	}
	err := s.db.Create(&event).Error
	if err != nil {
		logs.Error().Err(err).Msg("failed to seed events")
		return event, err
	}

	return event, nil
}
