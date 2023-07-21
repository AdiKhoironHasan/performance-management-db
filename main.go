package main

import (
	"engine-db/database"
	"engine-db/database/gather"
	"engine-db/database/seeder"
	"engine-db/entity"
	"fmt"
)

func main() {
	db := database.New()

	seeder := seeder.NewSeeder(db)
	gather := gather.NewGather(db)

	seeder.SeedUsersFromCSV()
	event, err := seeder.SeedOneEvent()
	if err != nil {
		panic(err)
	}

	var (
		sessions entity.SessionAssessment
		results  []entity.ResultSelfAssessment
	)

	participants, err := gather.GetParticipants()
	if err != nil {
		panic(err)
	}

	// get leader participants
	leaders, err := gather.GetLeaderParticipants()
	if err != nil {
		panic(err)
	}

	// get employee participants
	employee, err := gather.GetEmployeeParticipants()
	if err != nil {
		panic(err)
	}

	// ============ SESSION 1 - SELF REVIEW & LEADER REVIEW =============
	// SELF REVIEW

	// seed new self assesment session
	sessions.SelfReview, err = seeder.SeedSession(event.ID, "self_assess")
	if err != nil {
		panic(err)
	}

	// seed new self assesment questions
	sessions.SelfReview.Questions, err = seeder.SeedQuestion(10, sessions.SelfReview)
	if err != nil {
		panic(err)
	}

	// Task : self assessment
	for _, reviewee := range participants {
		for _, reviewer := range participants {
			err = seeder.SeedAssesmentTask(sessions.SelfReview, sessions.SelfReview.Questions, reviewee, reviewer)
			if err != nil {
				panic(err)
			}
		}
	}

	// Answer : self assessment
	for _, reviewee := range participants {
		for _, reviewer := range participants {
			err = seeder.SeedAnswerAssesmentTask(sessions.SelfReview, sessions.SelfReview.Questions, reviewee, reviewer)
			if err != nil {
				panic(err)
			}
		}
	}

	results, err = gather.GetSelfAssessmentTask(sessions.SelfReview)
	if err != nil {
		panic(err)
	}

	// show self assessment task
	for _, result := range results {
		fmt.Println("self : ", map[string]interface{}{
			"reviewee": result.Reviewee,
			"question": result.Question,
			"total":    result.Total,
		})
	}

	// // LEADER REVIEW
	// // Task : leader assessment
	// sessions.LeaderReview = sessions.SelfReview

	// sessions.LeaderReview.Questions, err = seeder.SeedSelfAssesmentTask(sessions.LeaderReview, participants)
	// if err != nil {
	// 	panic(err)
	// }

	// results, err = gather.GetSelfAssessmentTask(sessions.LeaderReview)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, result := range results {
	// 	fmt.Println("task : ", result)
	// }

	// // Answer : leader assessment
	// err = seeder.SeedAnswerSelfAssesmentTask(sessions.LeaderReview.Questions, participants)
	// if err != nil {
	// 	panic(err)
	// }

	// results, err = gather.GetSelfAssessmentTask(sessions.SelfReview)
	// if err != nil {
	// 	panic(err)
	// }

	// // show self assessment task
	// for _, result := range results {
	// 	fmt.Println("self : ", map[string]interface{}{
	// 		"reviewee": result.Reviewee,
	// 		"question": result.Question,
	// 		"total":    result.Total,
	// 	})
	// }

	// ========== SESSION 2 - CHOOSE PEERS ============
	// TODO: choose peers from leader
	sessions.PeersReview, err = seeder.SeedSession(event.ID, "peers_assess")
	if err != nil {
		panic(err)
	}

	sessions.PeersReview.Questions, err = seeder.SeedQuestion(10, sessions.PeersReview)
	if err != nil {
		panic(err)
	}

	for _, reviewer := range participants {
		for _, peer := range participants {
			if reviewer.ID != peer.ID {
				err = seeder.SeedAssesmentTask(sessions.PeersReview, sessions.PeersReview.Questions, peer, reviewer)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	// ========== SESSION 3 - PEERS REVIEW ============
	// Seed : new peers assesment task
	// TODO: get peers from leader

	// seed new peer assesment session
	// sessions.PeersReview, err = seeder.SeedSession(event.ID, "peers_assess")
	// if err != nil {
	// 	panic(err)
	// }

	// sessions.PeersReview.Questions, err = seeder.SeedQuestion(10, sessions.PeersReview)
	// if err != nil {
	// 	panic(err)
	// }

	// // Seed : new peers assesment task
	// // TODO: get peers from leader
	// for _, reviewer := range participants {
	// 	for _, peer := range participants {
	// 		if reviewer.ID != peer.ID {
	// 			err = seeder.SeedAssesmentTask(sessions.PeersReview, sessions.PeersReview.Questions, peer, reviewer)
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 		}
	// 	}
	// }

	// Answer : peers assesment task
	for _, reviewer := range participants {
		for _, peer := range participants {
			if reviewer.ID != peer.ID {
				err = seeder.SeedAnswerAssesmentTask(sessions.PeersReview, sessions.PeersReview.Questions, peer, reviewer)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	results, err = gather.GetPeersAssessmentTask(sessions.PeersReview)
	if err != nil {
		panic(err)
	}

	// show peers assessment task
	for _, result := range results {
		fmt.Println("peer : ", map[string]interface{}{
			"reviewee": result.Reviewee,
			"question": result.Question,
			"total":    result.Total,
		})
	}

	// ============ SESSION 4 - MEMBER REVIEW ============
	// TODO : get member from team

	sessions.MemberReview, err = seeder.SeedSession(event.ID, "member_assess")
	if err != nil {
		panic(err)
	}

	sessions.MemberReview.Questions, err = seeder.SeedQuestion(10, sessions.PeersReview)
	if err != nil {
		panic(err)
	}

	// Seed : new member assesment task
	for _, lead := range leaders {
		for _, emp := range employee {
			if lead.ID != emp.ID {
				err = seeder.SeedAssesmentTask(sessions.MemberReview, sessions.MemberReview.Questions, emp, lead)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	// Answer : member assesment task
	for _, lead := range leaders {
		for _, emp := range employee {
			if lead.ID != emp.ID {
				err = seeder.SeedAnswerAssesmentTask(sessions.MemberReview, sessions.MemberReview.Questions, emp, lead)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	reports, err := gather.GetAssesmentReportAverage([]entity.Session{
		sessions.SelfReview,
		sessions.PeersReview,
	})
	if err != nil {
		panic(err)
	}

	// show assesment report
	for _, report := range reports {
		fmt.Println(map[string]interface{}{
			"reviewee": report.Reviewee,
			"question": report.Question,
			"avg":      fmt.Sprintf("%.1f", report.Total),
		})
	}
}
