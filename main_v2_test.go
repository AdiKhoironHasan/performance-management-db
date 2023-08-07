package main

import (
	"database/sql"
	"encoding/csv"
	"engine-db/database"
	entity "engine-db/entity/new"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type DataSession struct {
	Session entity.Session

	Question []entity.Question

	FormTask []entity.FormTask

	UserAnswerText  []entity.FormScaleAnswer
	UserAnswerScale []entity.FormTextAnswer
}

func FindLeaderAndMember(db *gorm.DB, userVersion string, leadershipStatus string) ([][]int64, error) {
	query := `
						SELECT
						u1.id AS leader_id, u2.id AS member_id
						FROM users u1, users u2
						WHERE u1.privy_id = u2.direct_leader_employee_id
						AND u1.version = $1
						AND u2.version = $1
						`

	if leadershipStatus != "" {
		query += fmt.Sprintf("AND u2.leadership_status = '%s'", leadershipStatus)
	}

	rows, err := db.Raw(query, userVersion).Rows()
	if err != nil {
		return nil, err
	}

	var result [][]int64

	if rows != nil {
		for rows.Next() {
			var leaderID, memberID int64
			err := rows.Scan(&leaderID, &memberID)
			if err != nil {
				return nil, err
			}
			result = append(result, []int64{leaderID, memberID})
		}
	}

	// the result is [[leaderID, memberID]...]
	return result, nil
}

func FindMemberByLeaderID(db *gorm.DB, userVersion string, leaderID int64) ([]entity.User, error) {
	var (
		result []entity.User
		err    error
	)

	query := `
					SELECT u2.*
					FROM users u1, users u2
					WHERE
					u1.privy_id = u2.direct_leader_employee_id
					AND
					u1.version = $1 and u2.version = $1
					AND
					u1.id = $2;
					`

	err = db.Raw(query, userVersion, leaderID).Scan(&result).Error
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return result, err
}

func TestMainV2(t *testing.T) {
	const (
		EventStatusToDo       = "TODO"
		EventStatusInProgress = "IN_PROGRESS"
		EventStatusPassed     = "PASSED"

		SessionChoosePeers                 = "choose_peers"           //done
		SessionSelfAssessmentNonLeadership = "self_assess_non_lead"   //done
		SessionSelfAssessmentLeadership    = "self_assess_lead"       //done
		SessionPeerReviewNonLeadership     = "peers_assess_non_lead"  //done
		SessionPeerReviewLeadership        = "peers_assess_lead"      //done
		SessionTeamMemberNonLeadership     = "member_assess_non_lead" //done
		SessionTeamMemberLeadership        = "member_assess_lead"     //done
		SessionTeamMemberReviewLeader      = "leader_assess"          //done

		SessionStatusToDo       = "TODO"
		SessionStatusProcessed  = "PROCESSED"
		SessionStatusInProgress = "IN_PROGRESS"
		SessionStatusPassed     = "PASSED"

		QuestionTypeText  = "text"
		QuestionTypeScale = "scale"

		FormTasksStatusToDo     = "TODO"
		FormTasksStatusReviewed = "REVIEWED"
		FormTasksStatusPassed   = "PASSED"

		LeadershipStatusLeader    = "leader"
		LeadershipStatusNonLeader = "non leader"
	)

	var (
		dataEnterprise = entity.Enterprise{}

		dataUsers         = []entity.User{}
		dataUserLeader    = []entity.User{}
		dataUserNonLeader = []entity.User{}

		// dataUserVersionDefault = entity.UserVersion{}
		dataUserVersion = entity.UserVersion{}

		dataEvent = entity.Event{}

		f = faker.New()
	)

	// make connection db
	db, err := database.New()
	assert.NoError(t, err)

	// enterprise token
	enterpriseToken := "YOUR_ENTERPRISE_TOKEN"

	// create new enterprise
	err = db.Create(&entity.Enterprise{
		Name:     "PT TEST STAGING LUCU",
		Token:    enterpriseToken,
		IsActive: true,
		PrivyID:  "AGD8504",
	}).Scan(&dataEnterprise).Error
	assert.NoError(t, err)

	// super_admin
	defaultSuperAdminVersion := uuid.New().String()
	assert.NoError(t, err)
	err = db.Create(&entity.UserVersion{
		Version:         "default-" + defaultSuperAdminVersion,
		Ignore:          0,
		UserCount:       1,
		EnterpriseToken: dataEnterprise.Token,
	}).Scan(&dataEnterprise).Error
	assert.NoError(t, err)
	err = db.Create(&entity.User{
		EnterpriseToken:        enterpriseToken,
		Version:                "default-" + defaultSuperAdminVersion,
		Name:                   "Adnan Ganteng",
		PrivyID:                "AGD8504",
		Email:                  "default",
		Status:                 "default",
		JobTitle:               "default",
		Level:                  "default",
		Directorate:            "default",
		Division:               "default",
		Homebase:               "default",
		DirectLeader:           "default",
		DirectLeaderJobTitle:   "default",
		DirectLeaderEmployeeID: "default",
		PICHrbp:                "default",
		HrbpPrivyID:            "default",
		Role:                   "super_admin",
		LeadershipStatus:       "default",
	}).Scan(&dataEnterprise).Error
	assert.NoError(t, err)

	// 2. create new user version (default)
	// err = db.Create(&entity.UserVersion{
	// 	Version:         "default_seed",
	// 	EnterpriseToken: dataEnterprise.Token,
	// 	Ignore:          0,
	// 	UserCount:       0,
	// }).Scan(&dataUserVersionDefault).Error
	// assert.NoError(t, err)

	// Open the CSV file
	filePath := "storage/data-user-v6.csv"
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// create a CSV reader
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	num := 1

	user_version := "USR-" + uuid.New().String()
	// iterate over CSV rows
	for {
		// read 1 per 1 of rows
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)

		// skip row number 1 (using for header table)
		if num != 1 {
			leadershipStatus := strings.ToLower((row[16]))

			dataUsers = append(dataUsers, entity.User{
				EnterpriseToken:        dataEnterprise.Token,
				Version:                user_version,
				Name:                   row[0],
				PrivyID:                row[1],
				Email:                  row[2],
				Status:                 row[3],
				JobTitle:               row[5],
				Level:                  row[6],
				Directorate:            row[7],
				Division:               row[8],
				Homebase:               row[9],
				DirectLeader:           row[10],
				DirectLeaderJobTitle:   row[11],
				DirectLeaderEmployeeID: row[12],
				PICHrbp:                row[13],
				HrbpPrivyID:            row[14],
				Role:                   "employee",
				LeadershipStatus:       leadershipStatus,
			})
		}

		num++
	}

	err = db.Create(&entity.UserVersion{
		Version:         user_version,
		Ignore:          0,
		UserCount:       int64(len(dataUsers)),
		EnterpriseToken: dataEnterprise.Token,
	}).Scan(&dataUserVersion).Error
	assert.NoError(t, err)

	err = db.Create(&dataUsers).Error
	assert.NoError(t, err)

	err = db.Model(&entity.User{}).Where("id != ANY(ARRAY[1, 2])").Updates(map[string]interface{}{
		"version": dataUserVersion.Version,
	}).Error
	assert.NoError(t, err)

	err = db.Create(&entity.Event{
		Name:        "Event PA 1",
		Description: "Desc for Event PA 1",
		Status:      EventStatusToDo,
		UserVersion: dataUserVersion.Version,
	}).Scan(&dataEvent).Error
	assert.NoError(t, err)

	// preparation data user
	// find user non leader
	err = db.Model(&entity.User{}).
		Where("leadership_status = ?", LeadershipStatusNonLeader).
		Find(&dataUserNonLeader).Error
	assert.NoError(t, err)

	// find user leader
	err = db.Model(&entity.User{}).
		Where("leadership_status = ?", LeadershipStatusLeader).
		Find(&dataUserLeader).Error
	assert.NoError(t, err)

	// ====== 1. create self assessment non leadership ======
	var (
		selfAssesNonLeader = DataSession{}
	)

	err = db.Create(&entity.Session{
		EventID:     dataEvent.ID,
		Type:        SessionSelfAssessmentNonLeadership,
		Name:        "Self Assessment Non Leadership",
		Description: "Desc for Self Assessment Non Leadership",
		Status:      SessionStatusToDo,
		StartDate:   time.Now().AddDate(0, 0, 4),
		EndDate:     time.Now().AddDate(0, 0, 20),
	}).Scan(&selfAssesNonLeader.Session).Error
	assert.NoError(t, err)

	// create scale question
	for i := 1; i <= 10; i++ {
		questionType := QuestionTypeScale
		questionMax := 4

		if i > 7 {
			questionType = QuestionTypeText
			questionMax = 0
		}

		selfAssesNonLeader.Question = append(selfAssesNonLeader.Question, entity.Question{
			SessionID: selfAssesNonLeader.Session.ID,
			Sort:      i,
			Name:      fmt.Sprintf("Question %d", i),
			Type:      questionType,
			Max:       questionMax,
		})
	}

	err = db.Create(selfAssesNonLeader.Question).Error
	assert.NoError(t, err)

	// create form task
	for _, user := range dataUserNonLeader {
		selfAssesNonLeader.FormTask = append(selfAssesNonLeader.FormTask, entity.FormTask{
			SessionID:  selfAssesNonLeader.Session.ID,
			RevieweeID: user.ID,
			ReviewerID: user.ID,
			Status:     FormTasksStatusToDo,
		})
	}

	err = db.Create(selfAssesNonLeader.FormTask).Error
	assert.NoError(t, err)

	// user answer form task self assessment non leadership
	for _, formTask := range selfAssesNonLeader.FormTask {
		questions := []entity.Question{}
		err = db.Model(&entity.Question{}).Where("session_id = ?", selfAssesNonLeader.Session.ID).Find(&questions).Error
		assert.NoError(t, err)

		for _, question := range questions {
			if question.Type == QuestionTypeText {
				selfAssesNonLeader.UserAnswerScale = append(selfAssesNonLeader.UserAnswerScale, entity.FormTextAnswer{
					SessionID:  formTask.SessionID,
					QuestionID: question.ID,
					RevieweeID: formTask.RevieweeID,
					ReviewerID: formTask.ReviewerID,
					Sort:       question.Sort,
					TextValue:  fmt.Sprintf("Answer for question %d", question.Sort),
				})
			}

			if question.Type == QuestionTypeScale {
				selfAssesNonLeader.UserAnswerText = append(selfAssesNonLeader.UserAnswerText, entity.FormScaleAnswer{
					SessionID:  formTask.SessionID,
					QuestionID: question.ID,
					RevieweeID: formTask.RevieweeID,
					ReviewerID: formTask.ReviewerID,
					Sort:       question.Sort,
					ScaleValue: 4,
					Max:        question.Max,
				})
			}
		}
	}

	// create user answer for scale
	err = db.Create(selfAssesNonLeader.UserAnswerScale).Error
	assert.NoError(t, err)

	// create user answer for text
	err = db.Create(selfAssesNonLeader.UserAnswerText).Error
	assert.NoError(t, err)

	// ====== 2. create self assessment leadership ======
	var (
		selfAssesLeader = DataSession{}
	)

	err = db.Create(&entity.Session{
		EventID:     dataEvent.ID,
		Type:        SessionSelfAssessmentNonLeadership,
		Name:        "Self Assessment Non Leadership",
		Description: "Desc for Self Assessment Non Leadership",
		Status:      SessionStatusToDo,
		StartDate:   time.Now().AddDate(0, 0, 4),
		EndDate:     time.Now().AddDate(0, 0, 20),
	}).Scan(&selfAssesLeader.Session).Error
	assert.NoError(t, err)

	// create scale question
	for i := 1; i <= 10; i++ {
		questionType := QuestionTypeScale
		questionMax := 4

		if i > 7 {
			questionType = QuestionTypeText
			questionMax = 0
		}

		selfAssesLeader.Question = append(selfAssesLeader.Question, entity.Question{
			SessionID: selfAssesLeader.Session.ID,
			Sort:      i,
			Name:      fmt.Sprintf("Question %d", i),
			Type:      questionType,
			Max:       questionMax,
		})
	}

	err = db.Create(selfAssesLeader.Question).Error
	assert.NoError(t, err)

	// create form task
	for _, user := range dataUserNonLeader {
		selfAssesLeader.FormTask = append(selfAssesLeader.FormTask, entity.FormTask{
			SessionID:  selfAssesLeader.Session.ID,
			RevieweeID: user.ID,
			ReviewerID: user.ID,
			Status:     FormTasksStatusToDo,
		})
	}

	err = db.Create(selfAssesLeader.FormTask).Error
	assert.NoError(t, err)

	// user answer form task self assessment non leadership
	for _, formTask := range selfAssesLeader.FormTask {
		questions := []entity.Question{}
		err = db.Model(&entity.Question{}).Where("session_id = ?", selfAssesLeader.Session.ID).Find(&questions).Error
		assert.NoError(t, err)

		for _, question := range questions {
			if question.Type == QuestionTypeText {
				selfAssesLeader.UserAnswerScale = append(selfAssesLeader.UserAnswerScale, entity.FormTextAnswer{
					SessionID:  formTask.SessionID,
					QuestionID: question.ID,
					RevieweeID: formTask.RevieweeID,
					ReviewerID: formTask.ReviewerID,
					Sort:       question.Sort,
					TextValue:  fmt.Sprintf("Answer for question %d", question.Sort),
				})
			}

			if question.Type == QuestionTypeScale {
				selfAssesLeader.UserAnswerText = append(selfAssesLeader.UserAnswerText, entity.FormScaleAnswer{
					SessionID:  formTask.SessionID,
					QuestionID: question.ID,
					RevieweeID: formTask.RevieweeID,
					ReviewerID: formTask.ReviewerID,
					Sort:       question.Sort,
					ScaleValue: 4,
					Max:        question.Max,
				})
			}
		}
	}

	// create user answer for scale
	err = db.Create(selfAssesLeader.UserAnswerScale).Error
	assert.NoError(t, err)

	// create user answer for text
	err = db.Create(selfAssesLeader.UserAnswerText).Error
	assert.NoError(t, err)

	// ====== 3. create leader assesment ======
	var (
		leaderAssess = DataSession{}
	)

	err = db.Create(&entity.Session{
		EventID:     dataEvent.ID,
		Type:        SessionTeamMemberReviewLeader,
		Name:        "Leader Assessment",
		Description: "Desc for Leader Assessment",
		Status:      SessionStatusToDo,
		StartDate:   time.Now().AddDate(0, 0, 4),
		EndDate:     time.Now().AddDate(0, 0, 20),
	}).Scan(&leaderAssess.Session).Error
	assert.NoError(t, err)

	// create questions
	for i := 1; i <= 10; i++ {
		questionType := QuestionTypeScale
		questionMax := 4

		if i > 7 {
			questionType = QuestionTypeText
			questionMax = 0
		}

		leaderAssess.Question = append(leaderAssess.Question, entity.Question{
			SessionID: leaderAssess.Session.ID,
			Sort:      i,
			Name:      fmt.Sprintf("Question %d", i),
			Type:      questionType,
			Max:       questionMax,
		})
	}

	err = db.Create(leaderAssess.Question).Error
	assert.NoError(t, err)

	leaderMember, err := FindLeaderAndMember(db, dataEvent.UserVersion, "")
	assert.NoError(t, err)

	// iterate through leader and members
	for _, val := range leaderMember {
		leaderAssess.FormTask = append(leaderAssess.FormTask, entity.FormTask{
			SessionID:  leaderAssess.Session.ID,
			RevieweeID: val[0],
			ReviewerID: val[1],
			Status:     FormTasksStatusToDo,
		})
	}

	err = db.Create(leaderAssess.FormTask).Error
	assert.NoError(t, err)

	// user answer form task self assessment non leadership
	for _, formTask := range leaderAssess.FormTask {
		questions := []entity.Question{}
		err = db.Model(&entity.Question{}).Where("session_id = ?", leaderAssess.Session.ID).Find(&questions).Error
		assert.NoError(t, err)

		for _, question := range questions {
			if question.Type == QuestionTypeText {
				leaderAssess.UserAnswerScale = append(leaderAssess.UserAnswerScale, entity.FormTextAnswer{
					SessionID:  formTask.SessionID,
					QuestionID: question.ID,
					RevieweeID: formTask.RevieweeID,
					ReviewerID: formTask.ReviewerID,
					Sort:       question.Sort,
					TextValue:  fmt.Sprintf("Answer for question %d", question.Sort),
				})
			}

			if question.Type == QuestionTypeScale {
				leaderAssess.UserAnswerText = append(leaderAssess.UserAnswerText, entity.FormScaleAnswer{
					SessionID:  formTask.SessionID,
					QuestionID: question.ID,
					RevieweeID: formTask.RevieweeID,
					ReviewerID: formTask.ReviewerID,
					Sort:       question.Sort,
					ScaleValue: 4,
					Max:        question.Max,
				})
			}
		}
	}

	// create user answer for scale
	err = db.Create(leaderAssess.UserAnswerScale).Error
	assert.NoError(t, err)

	// create user answer for text
	err = db.Create(leaderAssess.UserAnswerText).Error
	assert.NoError(t, err)

	// ====== 4. create member assess leadership ======
	var (
		memberAssessLeader = DataSession{}
	)

	err = db.Create(&entity.Session{
		EventID:     dataEvent.ID,
		Type:        SessionTeamMemberLeadership,
		Name:        "Leader Assessment Leadership",
		Description: "Desc for Leader Assessment Leadership",
		Status:      SessionStatusToDo,
		StartDate:   time.Now().AddDate(0, 0, 4),
		EndDate:     time.Now().AddDate(0, 0, 20),
	}).Scan(&memberAssessLeader.Session).Error
	assert.NoError(t, err)

	// create questions
	for i := 1; i <= 10; i++ {
		questionType := QuestionTypeScale
		questionMax := 4

		if i > 7 {
			questionType = QuestionTypeText
			questionMax = 0
		}

		memberAssessLeader.Question = append(memberAssessLeader.Question, entity.Question{
			SessionID: memberAssessLeader.Session.ID,
			Sort:      i,
			Name:      fmt.Sprintf("Question %d", i),
			Type:      questionType,
			Max:       questionMax,
		})
	}

	err = db.Create(memberAssessLeader.Question).Error
	assert.NoError(t, err)

	leaderMember, err = FindLeaderAndMember(db, dataEvent.UserVersion, LeadershipStatusLeader)
	assert.NoError(t, err)

	// iterate through leader and members
	for _, val := range leaderMember {
		memberAssessLeader.FormTask = append(memberAssessLeader.FormTask, entity.FormTask{
			SessionID:  memberAssessLeader.Session.ID,
			RevieweeID: val[1],
			ReviewerID: val[0],
			Status:     FormTasksStatusToDo,
		})
	}

	// store form tasks
	err = db.Create(memberAssessLeader.FormTask).Error
	assert.NoError(t, err)

	// user answer form task self assessment non leadership
	for _, formTask := range memberAssessLeader.FormTask {
		questions := []entity.Question{}
		err = db.Model(&entity.Question{}).Where("session_id = ?", memberAssessLeader.Session.ID).Find(&questions).Error
		assert.NoError(t, err)

		for _, question := range questions {
			if question.Type == QuestionTypeText {
				memberAssessLeader.UserAnswerScale = append(memberAssessLeader.UserAnswerScale, entity.FormTextAnswer{
					SessionID:  formTask.SessionID,
					QuestionID: question.ID,
					RevieweeID: formTask.RevieweeID,
					ReviewerID: formTask.ReviewerID,
					Sort:       question.Sort,
					TextValue:  fmt.Sprintf("Answer for question %d", question.Sort),
				})
			}

			if question.Type == QuestionTypeScale {
				memberAssessLeader.UserAnswerText = append(memberAssessLeader.UserAnswerText, entity.FormScaleAnswer{
					SessionID:  formTask.SessionID,
					QuestionID: question.ID,
					RevieweeID: formTask.RevieweeID,
					ReviewerID: formTask.ReviewerID,
					Sort:       question.Sort,
					ScaleValue: 4,
					Max:        question.Max,
				})
			}
		}
	}

	// create user answer for scale
	err = db.Create(memberAssessLeader.UserAnswerScale).Error
	assert.NoError(t, err)

	// create user answer for text
	err = db.Create(memberAssessLeader.UserAnswerText).Error
	assert.NoError(t, err)

	// ====== 5. create member assess leadership ======
	var (
		memberAssessLeaderNonLeader = DataSession{}
	)

	err = db.Create(&entity.Session{
		EventID:     dataEvent.ID,
		Type:        SessionTeamMemberLeadership,
		Name:        "Leader Assessment Non Leadership",
		Description: "Desc for Leader Assessment Non Leadership",
		Status:      SessionStatusToDo,
		StartDate:   time.Now().AddDate(0, 0, 4),
		EndDate:     time.Now().AddDate(0, 0, 20),
	}).Scan(&memberAssessLeaderNonLeader.Session).Error
	assert.NoError(t, err)

	// create questions
	for i := 1; i <= 10; i++ {
		questionType := QuestionTypeScale
		questionMax := 4

		if i == 1 {
			questionMax = 5
		}
		if i > 7 {
			questionType = QuestionTypeText
			questionMax = 0
		}

		memberAssessLeaderNonLeader.Question = append(memberAssessLeaderNonLeader.Question, entity.Question{
			SessionID: memberAssessLeaderNonLeader.Session.ID,
			Sort:      i,
			Name:      fmt.Sprintf("Question %d", i),
			Type:      questionType,
			Max:       questionMax,
		})
	}

	err = db.Create(memberAssessLeaderNonLeader.Question).Error
	assert.NoError(t, err)

	leaderMember, err = FindLeaderAndMember(db, dataEvent.UserVersion, LeadershipStatusNonLeader)
	assert.NoError(t, err)

	// iterate through leader and members
	for _, val := range leaderMember {
		memberAssessLeaderNonLeader.FormTask = append(memberAssessLeaderNonLeader.FormTask, entity.FormTask{
			SessionID:  memberAssessLeaderNonLeader.Session.ID,
			RevieweeID: val[1],
			ReviewerID: val[0],
			Status:     FormTasksStatusToDo,
		})
	}

	// store form tasks
	err = db.Create(memberAssessLeaderNonLeader.FormTask).Error
	assert.NoError(t, err)

	// user answer form task self assessment non leadership
	for _, formTask := range memberAssessLeaderNonLeader.FormTask {
		questions := []entity.Question{}
		err = db.Model(&entity.Question{}).Where("session_id = ?", memberAssessLeaderNonLeader.Session.ID).Find(&questions).Error
		assert.NoError(t, err)

		for _, question := range questions {
			poin := 4

			if question.Max == 5 {
				poin = 5
			}

			if question.Type == QuestionTypeText {
				memberAssessLeaderNonLeader.UserAnswerScale = append(memberAssessLeaderNonLeader.UserAnswerScale, entity.FormTextAnswer{
					SessionID:  formTask.SessionID,
					QuestionID: question.ID,
					RevieweeID: formTask.RevieweeID,
					ReviewerID: formTask.ReviewerID,
					Sort:       question.Sort,
					TextValue:  fmt.Sprintf("Answer for question %d", question.Sort),
				})
			}

			if question.Type == QuestionTypeScale {
				memberAssessLeaderNonLeader.UserAnswerText = append(memberAssessLeaderNonLeader.UserAnswerText, entity.FormScaleAnswer{
					SessionID:  formTask.SessionID,
					QuestionID: question.ID,
					RevieweeID: formTask.RevieweeID,
					ReviewerID: formTask.ReviewerID,
					Sort:       question.Sort,
					ScaleValue: float64(poin),
					Max:        question.Max,
				})
			}
		}
	}

	// create user answer for scale
	err = db.Create(memberAssessLeaderNonLeader.UserAnswerScale).Error
	assert.NoError(t, err)

	// create user answer for text
	err = db.Create(memberAssessLeaderNonLeader.UserAnswerText).Error
	assert.NoError(t, err)

	// ====== 6. create peers assess leadership ======
	var (
		peersAssessLeader = DataSession{}
	)

	err = db.Create(&entity.Session{
		EventID:     dataEvent.ID,
		Type:        SessionPeerReviewLeadership,
		Name:        "Peers Assess Leadership",
		Description: "Desc for Peers Assess Leadership",
		Status:      SessionStatusToDo,
		StartDate:   time.Now().AddDate(0, 0, 4),
		EndDate:     time.Now().AddDate(0, 0, 20),
	}).Scan(&peersAssessLeader.Session).Error
	assert.NoError(t, err)

	// create questions
	for i := 1; i <= 10; i++ {
		questionType := QuestionTypeScale
		questionMax := 4

		if i > 7 {
			questionType = QuestionTypeText
			questionMax = 0
		}

		peersAssessLeader.Question = append(peersAssessLeader.Question, entity.Question{
			SessionID: peersAssessLeader.Session.ID,
			Sort:      i,
			Name:      fmt.Sprintf("Question %d", i),
			Type:      questionType,
			Max:       questionMax,
		})
	}

	err = db.Create(peersAssessLeader.Question).Error
	assert.NoError(t, err)

	// ====== 7. create peers assess non leadership ======
	var (
		peersAssessNonLeader = DataSession{}
	)

	err = db.Create(&entity.Session{
		EventID:     dataEvent.ID,
		Type:        SessionPeerReviewNonLeadership,
		Name:        "Peers Assess Non Leadership",
		Description: "Desc for Peers Assess Non Leadership",
		Status:      SessionStatusToDo,
		StartDate:   time.Now().AddDate(0, 0, 4),
		EndDate:     time.Now().AddDate(0, 0, 20),
	}).Scan(&peersAssessNonLeader.Session).Error
	assert.NoError(t, err)

	// create questions
	for i := 1; i <= 10; i++ {
		questionType := QuestionTypeScale
		questionMax := 4

		if i > 7 {
			questionType = QuestionTypeText
			questionMax = 0
		}

		peersAssessNonLeader.Question = append(peersAssessNonLeader.Question, entity.Question{
			SessionID: peersAssessNonLeader.Session.ID,
			Sort:      i,
			Name:      fmt.Sprintf("Question %d", i),
			Type:      questionType,
			Max:       questionMax,
		})
	}

	err = db.Create(peersAssessNonLeader.Question).Error
	assert.NoError(t, err)

	// ====== 8. create choose peers ======
	var (
		choosePeers = DataSession{}
	)

	err = db.Create(&entity.Session{
		EventID:     dataEvent.ID,
		Type:        SessionChoosePeers,
		Name:        "Choose Peers",
		Description: "Desc for Choose Peers",
		Status:      SessionStatusToDo,
		StartDate:   time.Now().AddDate(0, 0, 4),
		EndDate:     time.Now().AddDate(0, 0, 20),
	}).Scan(&choosePeers.Session).Error
	assert.NoError(t, err)

	userIDs := []int{}
	for _, val := range dataUsers {
		userIDs = append(userIDs, int(val.ID))
	}

	// iterate leader user
	for _, val := range dataUserLeader {
		// get members detail
		members, err := FindMemberByLeaderID(db, val.Version, val.ID)
		assert.NoError(t, err)

		for _, member := range members {
			if member.LeadershipStatus == LeadershipStatusLeader {
				peersAssessLeader.FormTask = append(peersAssessLeader.FormTask, entity.FormTask{
					SessionID:  peersAssessLeader.Session.ID,
					RevieweeID: member.ID,
					ReviewerID: int64(f.RandomIntElement(userIDs)),
					Status:     FormTasksStatusToDo,
				})
			} else {
				peersAssessNonLeader.FormTask = append(peersAssessNonLeader.FormTask, entity.FormTask{
					SessionID:  peersAssessNonLeader.Session.ID,
					RevieweeID: member.ID,
					ReviewerID: int64(f.RandomIntElement(userIDs)),
					Status:     FormTasksStatusToDo,
				})
			}
		}
	}

	// store form tasks answer for peers assess leadership
	for _, formTask := range peersAssessNonLeader.FormTask {
		for _, question := range peersAssessNonLeader.Question {
			if question.Type == QuestionTypeText {
				peersAssessLeader.UserAnswerScale = append(peersAssessLeader.UserAnswerScale, entity.FormTextAnswer{
					SessionID:  formTask.SessionID,
					QuestionID: question.ID,
					RevieweeID: formTask.RevieweeID,
					ReviewerID: formTask.ReviewerID,
					Sort:       question.Sort,
					TextValue:  fmt.Sprintf("Answer for question %d", question.Sort),
				})
			}

			if question.Type == QuestionTypeScale {
				peersAssessLeader.UserAnswerText = append(peersAssessLeader.UserAnswerText, entity.FormScaleAnswer{
					SessionID:  formTask.SessionID,
					QuestionID: question.ID,
					RevieweeID: formTask.RevieweeID,
					ReviewerID: formTask.ReviewerID,
					Sort:       question.Sort,
					ScaleValue: 4,
					Max:        question.Max,
				})
			}
		}
	}

	// create user answer for scale
	err = db.Create(peersAssessLeader.UserAnswerScale).Error
	assert.NoError(t, err)

	// create user answer for text
	err = db.Create(peersAssessLeader.UserAnswerText).Error
	assert.NoError(t, err)

	// store form tasks answer for peers assess non leadership
	for _, formTask := range peersAssessNonLeader.FormTask {
		for _, question := range peersAssessNonLeader.Question {
			if question.Type == QuestionTypeText {
				peersAssessNonLeader.UserAnswerScale = append(peersAssessNonLeader.UserAnswerScale, entity.FormTextAnswer{
					SessionID:  formTask.SessionID,
					QuestionID: question.ID,
					RevieweeID: formTask.RevieweeID,
					ReviewerID: formTask.ReviewerID,
					Sort:       question.Sort,
					TextValue:  fmt.Sprintf("Answer for question %d", question.Sort),
				})
			}

			if question.Type == QuestionTypeScale {
				peersAssessNonLeader.UserAnswerText = append(peersAssessNonLeader.UserAnswerText, entity.FormScaleAnswer{
					SessionID:  formTask.SessionID,
					QuestionID: question.ID,
					RevieweeID: formTask.RevieweeID,
					ReviewerID: formTask.ReviewerID,
					Sort:       question.Sort,
					ScaleValue: 4,
					Max:        question.Max,
				})
			}
		}
	}

	// create user answer for scale
	err = db.Create(peersAssessNonLeader.UserAnswerScale).Error
	assert.NoError(t, err)

	// create user answer for text
	err = db.Create(peersAssessNonLeader.UserAnswerText).Error
	assert.NoError(t, err)
}
