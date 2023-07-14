package entity

// new custom type entity
type ResultSelfAssessment struct {
	Reviewee string
	Question string
	Total    float64

	Point float64
	Count int64
}

type SessionAssessment struct {
	SelfReview   Session
	LeaderReview Session

	PeersReview Session

	MemberReview Session
}
