package gather

import "engine-db/entity"

type Gather interface {
	GetSelfAssessment() error
	GetParticipants() ([]entity.User, error)
	GetSelfAssessmentTask(session entity.Session) ([]entity.ResultSelfAssessment, error)

	GetPeersAssessmentTask(session entity.Session) ([]entity.ResultSelfAssessment, error)

	// assesment report
	GetAssesmentReport(session []entity.Session) ([]entity.ResultSelfAssessment, error)
	GetAssesmentReportAverage(session []entity.Session) ([]entity.ResultSelfAssessment, error)

	GetLeaderParticipants() ([]entity.User, error)
	GetEmployeeParticipants() ([]entity.User, error)
}
