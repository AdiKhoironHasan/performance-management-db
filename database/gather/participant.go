package gather

import (
	"engine-db/entity"
)

func (g *gather) GetParticipants() ([]entity.User, error) {
	users := []entity.User{}

	err := g.db.Model(&entity.User{}).
		Where(`role = $1`, "admin").
		Or(`role = $2`, "leader").
		Or(`role = $3`, "employee").
		Find(&users).Error
	if err != nil {
		logs.Error().Err(err).Msg("failed to find users")
		return users, err
	}

	return users, nil
}

func (g *gather) GetLeaderParticipants() ([]entity.User, error) {
	users := []entity.User{}

	err := g.db.Model(&entity.User{}).
		Where(`role = $1`, "leader").
		Find(&users).Error
	if err != nil {
		logs.Error().Err(err).Msg("failed to find users")
		return users, err
	}

	return users, nil
}

func (g *gather) GetEmployeeParticipants() ([]entity.User, error) {
	users := []entity.User{}

	err := g.db.Model(&entity.User{}).
		Where(`role = $1`, "employee").
		Find(&users).Error
	if err != nil {
		logs.Error().Err(err).Msg("failed to find users")
		return users, err
	}

	return users, nil
}
