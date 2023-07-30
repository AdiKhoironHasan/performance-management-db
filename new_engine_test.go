package main

import (
	"engine-db/database"
	"engine-db/entity"
	"fmt"
	"strings"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestEngineDBNew(t *testing.T) {
	db := database.New()

	scaleAnswers := []entity.FormScaleAnswer{}
	f := faker.New()
	reviewee := 1

	// for d := 1; d <= 3; d++ {

	loop := 0
	for a := 1; a <= 10; a++ {

		reviewer := 1
		for b := 1; b <= 10; b++ {

			question := 1
			for c := 1; c <= 10; c++ {
				scaleAnswers = append(scaleAnswers, entity.FormScaleAnswer{
					SessionID:  1,
					RevieweeID: reviewee,
					ReviewerID: reviewer,
					QuestionID: question,
					Sort:       question,
					ScaleValue: f.RandomIntElement([]int{1, 2, 3, 4}),
					Max:        4,
				})

				loop++
				question++
			}

			reviewer++
		}

		reviewee++
	}
	fmt.Println("loooooop : ", loop)
	// }

	scaleAnswerValues := []string{}

	for _, v := range scaleAnswers {
		scaleAnswerValues = append(scaleAnswerValues, fmt.Sprintf("(%v, %v, %v, %v, %v, %v, %v)", v.SessionID, v.RevieweeID, v.ReviewerID, v.QuestionID, v.Sort, v.ScaleValue, v.Max))
	}

	query := fmt.Sprintf("INSERT INTO form_scale_answers (session_id, reviewee_id, reviewer_id, question_id, sort, scale_value, max) VALUES %s", strings.Join(scaleAnswerValues, ","))

	err := db.Exec(query).Error
	assert.NoError(t, err)

}

func TestReport(t *testing.T) {
	db := database.New()

	type ReportData struct {
		ID            int64     `json:"id" gorm:"column:reviewee_id"`
		Name          string    `json:"name" gorm:"column:name"`
		SelfAssess    float64   `json:"self_assess" gorm:"column:self_assess"`
		Average       float64   `json:"average" gorm:"column:avg"`
		MemberAsses   float64   `json:"member_asses" gorm:"column:leader_asses"`
		PeersAsses    float64   `json:"peers_asses" gorm:"column:peers_asses"`
		PeersAsseses  []float64 `json:"peers_asseses" gorm:"-"`
		TotalReviewer float64   `json:"total_reviewer" gorm:"column:total_reviewer"`
		FinalScore    float64   `json:"final_score" gorm:"column:member_asses"`
	}

	// TODO: kurang leader assess
	reports := []ReportData{}
	rows, err := db.Raw(`
	SELECT distinct 
		fsa.reviewee_id                                        AS reviewee_id
		     
		       ,COUNT(DISTINCT fsa.reviewer_id)                        AS total_reviewer
		       ,SUM(fsa.scale_value) / COUNT(DISTINCT fsa.reviewer_id) AS avg

		       ,self_assess.avg                                  AS self_assess
		       ,member_assess.avg as member_assess
		        ,peers_assess.avg as peers_assess
		FROM form_scale_answers fsa
		JOIN users
		ON fsa.reviewee_id = users.id
		left join (
			SELECT  distinct
			reviewee_id
			,reviewer_id
			       ,SUM(scale_value)
			       ,COUNT(DISTINCT reviewer_id)
			       ,SUM(scale_value) / COUNT(DISTINCT question_id) AS avg
			FROM form_scale_answers
			WHERE reviewer_id != 1
			AND session_id IN (1, 2, 3)
			GROUP BY  reviewee_id, reviewer_id
		) AS peers_assess on peers_assess.reviewee_id = fsa.reviewee_id
		LEFT JOIN
		(
			SELECT  distinct reviewee_id
			       ,SUM(scale_value)
			       ,COUNT(DISTINCT reviewer_id)
			       ,SUM(scale_value) / COUNT(DISTINCT question_id) AS avg
			FROM form_scale_answers
			WHERE reviewer_id = reviewee_id
			AND session_id IN (1, 2, 3)
			GROUP BY  reviewee_id
		) AS self_assess on self_assess.reviewee_id = fsa.reviewee_id
		left join 
		(
			SELECT  distinct reviewee_id
			       ,SUM(scale_value)
			       ,COUNT(DISTINCT reviewer_id)
			       ,SUM(scale_value) / COUNT(DISTINCT question_id) AS avg
			FROM form_scale_answers
			WHERE reviewer_id = 1
			AND session_id IN (1, 2, 3)
			GROUP BY  reviewee_id
		) AS member_assess ON member_assess.reviewee_id = fsa.reviewee_id
		WHERE fsa.session_id IN (1, 2, 3) and fsa.reviewee_id in (1,2)
		GROUP BY  fsa.reviewee_id
		         ,users.name
		         , self_assess.avg  
		         ,member_assess.avg
		          ,peers_assess.avg
	`).Rows()
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		report := ReportData{}
		err := rows.Scan(
			&report.ID,
			&report.TotalReviewer,
			&report.Average,
			&report.SelfAssess,
			&report.FinalScore,
			&report.PeersAsses,
		)
		assert.NoError(t, err)

		reports = append(reports, report)
	}

	fmt.Println(reports)

	reportList := []ReportData{}

	reportMaps := map[int64]*ReportData{}
	for _, val := range reports {
		if _, ok := reportMaps[val.ID]; !ok {
			reportMaps[val.ID] = &ReportData{
				ID:         val.ID,
				Name:       "Your Name",
				SelfAssess: val.SelfAssess,
				// MemberAsses: val.MemberAsses,
				Average: val.Average,
				PeersAsseses: []float64{
					val.PeersAsses,
				},
				TotalReviewer: val.TotalReviewer,
				FinalScore:    val.FinalScore,
			}

			continue
		}

		reportMaps[val.ID].PeersAsseses = append(reportMaps[val.ID].PeersAsseses, val.PeersAsses)
	}

	for _, val := range reportMaps {
		reportList = append(reportList, *val)
	}

	fmt.Println(reportList)

}
