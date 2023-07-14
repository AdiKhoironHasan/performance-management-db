package seeder

import "engine-db/entity"

type Seeder interface {
	SeedUsers()
	SeedUsersFromCSV()
	SeedEventCycle()
	SeedQuestionAnswer()

	// event
	SeedOneEvent() (entity.Event, error)

	// seesion
	SeedSession(eventID int, sessionType string) (entity.Session, error)
	// SeedSessionSelfAssesment(eventID int) (entity.Session, error)
	// SeedSessionPeerAssesment(eventID int) (entity.Session, error)

	// question
	SeedQuestion(count int, session entity.Session) ([]entity.Question, error)

	// // self assesment task
	// SeedSelfAssesmentTask(session entity.Session, participants []entity.User) ([]entity.Question, error)
	// SeedAnswerSelfAssesmentTask(questions []entity.Question, participants []entity.User) error

	// // peers assesment task
	// SeedPeersAssesmentTask(session entity.Session, questions []entity.Question, reviewer, peer entity.User) error
	// SeedAnswerPeersAssesmentTask(questions []entity.Question, participants []entity.User) error

	// assessments
	SeedAssesmentTask(session entity.Session, questions []entity.Question, reviewe, reviewer entity.User) error
	SeedAnswerAssesmentTask(session entity.Session, questions []entity.Question, reviewe, reviewer entity.User) error
}
