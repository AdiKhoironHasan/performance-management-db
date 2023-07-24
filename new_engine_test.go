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

	scaleAnswers := []entity.FormAnswerScale{}
	f := faker.New()
	reviewee := 1

	// for d := 1; d <= 3; d++ {

	loop := 0
	for a := 1; a <= 10; a++ {

		reviewer := 1
		for b := 1; b <= 10; b++ {

			question := 1
			for c := 1; c <= 10; c++ {
				scaleAnswers = append(scaleAnswers, entity.FormAnswerScale{
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
