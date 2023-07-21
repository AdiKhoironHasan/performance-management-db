package seeder

import (
	"engine-db/entity"
	"engine-db/logger"
	"engine-db/storage"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"gorm.io/gorm"
)

var (
	logs = logger.NewLogger().Logger.With().Str("pkg", "seeder").Logger()
)

type seeder struct {
	db    *gorm.DB
	faker *faker.Faker
}

func NewSeeder(db *gorm.DB) Seeder {
	f := faker.New()
	return &seeder{
		db:    db,
		faker: &f,
	}
}

func (s *seeder) SeedUsersFromCSV() {
	users := storage.ReadCsv()

	u := uuid.New().String()
	err := s.db.Create(entity.Versioning{
		ID:      u,
		Version: 1,
		Offset:  len(users),
	}).Error
	if err != nil {
		logs.Error().Err(err).Msg("failed to seed users")
		os.Exit(1)
	}

	for k, _ := range users {
		users[k].Version = u
	}

	err = s.db.Create(&users).Error
	if err != nil {
		logs.Error().Err(err).Msg("failed to seed users")
		os.Exit(1)
	}
}

func (s *seeder) SeedUsers() {
	for i := 1; i <= 20; i++ {
		role := "employee"

		switch i {
		case 1:
			role = "superadmin"
		case 2:
			role = "admin"
		case 3:
			role = "admin"
		case 4:
			role = "leader"
		case 5:
			role = "leader"
		case 6:
			role = "leader"
		case 7:
			role = "leader"
		}

		user := entity.User{
			PrivyID: s.faker.Internet().Email(),
			Name:    s.faker.Person().Name(),
			Role:    role,
		}

		err := s.db.Create(&user).Error
		if err != nil {
			logs.Error().Err(err).Msg("failed to seed users")
			os.Exit(1)
		}
	}
}

func (s *seeder) SeedEventCycle() {
	events := []entity.Event{
		{
			Name:        fmt.Sprintf("Event %v", 1),
			Description: fmt.Sprintf("Description for Event %v", 1),
		},
	}
	err := s.db.Create(&events).Error
	if err != nil {
		logs.Error().Err(err).Msg("failed to seed events")
		os.Exit(1)
	}

	sessions := []entity.Session{}
	for i := 1; i <= 2; i++ {
		sessionType := "self_assess"
		date := Date{
			StartDate: time.Now(),
			EndDate:   time.Now(),
		}

		switch i {
		// case 2:
		// 	sessionType = "choose_peers"
		// 	date.StartDate = date.StartDate.AddDate(0, 0, 10)
		// 	date.EndDate = date.EndDate.AddDate(0, 1, 10)
		case 2:
			sessionType = "peers_assess"
			date.StartDate = date.StartDate.AddDate(0, 1, 11)
			date.EndDate = date.EndDate.AddDate(0, 2, 11)
		case 3:
			sessionType = "member_assess"
			date.StartDate = date.StartDate.AddDate(0, 2, 12)
			date.EndDate = date.EndDate.AddDate(0, 2, 12)
		}

		sessions = append(sessions, entity.Session{
			EventID:     1,
			Type:        sessionType,
			StartDate:   date.StartDate,
			EndDate:     date.EndDate,
			Name:        fmt.Sprintf("Session %v", i),
			Description: fmt.Sprintf("Description for Session %v", 1),
		})
	}
	err = s.db.Create(&sessions).Error
	if err != nil {
		logs.Error().Err(err).Msg("failed to seed sessions")
		os.Exit(1)
	}

	questions := []entity.Question{}
	for _, v := range sessions {
		// sort alway reset by session
		sort := 1

		for i := 1; i <= 10; i++ {
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
				SessionID: v.ID,
				Sort:      sort,
				Name:      fmt.Sprintf("Question %v", i),
				Type:      qusetionType,
				Option:    "",
				Required:  s.faker.RandomStringElement([]string{"true", "false"}),
			})

			sort++
		}
	}
	err = s.db.Create(&questions).Error
	if err != nil {
		logs.Error().Err(err).Msg("failed to seed questions")
		os.Exit(1)
	}

}

func (s *seeder) SeedQuestionAnswer() {
	users := []entity.User{}

	err := s.db.Model(&entity.User{}).
		Where(`role = $1`, "admin").
		Or(`role = $2`, "leader").
		Or(`role = $3`, "employee").
		Find(&users).Error
	if err != nil {
		logs.Error().Err(err).Msg("failed to find users")
		os.Exit(1)
	}

	questions := []entity.Question{}
	// self review questions
	err = s.db.Model(&entity.Question{}).
		Where("session_id = $1", 1).
		Where("type = $2", "scale").
		Find(&questions).Error

	if err != nil {
		logs.Error().Err(err).Msg("failed to find questions")
		os.Exit(1)
	}

	loops := 1
	for _, user := range users {
		for _, question := range questions {

			switch user.Role {
			case "leader", "employee", "admin":
				answer := entity.QuestionAnswer{
					SessionID:  1,
					RevieweeID: user.ID,
					ReviewerID: user.ID,
					QuestionID: question.ID,
					Status:     "reviewed",
					Scale:      4,
				}
				err := s.db.Create(&answer).Error
				if err != nil {
					logs.Error().Err(err).Msg("failed to seed question answer")
					os.Exit(1)
				}
			}

			loops++
		}
	}

	fmt.Println("users : ", len(users))
	fmt.Println("loops : ", loops)
}
