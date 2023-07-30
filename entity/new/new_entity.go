package new_entity

import (
	"database/sql"
	"time"
)

type User struct {
	ID                     int64        `db:"id,omitempty" json:"id"`
	Version                string       `db:"version,omitempty" json:"version"`
	EnterpriseToken        string       `db:"enterprise_token,omitempty" json:"enterprise_token"`
	Name                   string       `db:"name,omitempty" json:"name"`
	PrivyID                string       `db:"privy_id,omitempty" json:"privy_id"`
	Email                  string       `db:"email,omitempty" json:"email"`
	Status                 string       `db:"status,omitempty" json:"status"`
	JoinDate               time.Time    `db:"join_date,omitempty" json:"join_date"`
	JobTitle               string       `db:"job_title,omitempty" json:"job_title"`
	Level                  string       `db:"level,omitempty" json:"level"`
	Directorate            string       `db:"directorate,omitempty" json:"directorate"`
	Division               string       `db:"division,omitempty" json:"division"`
	Homebase               string       `db:"homebase,omitempty" json:"homebase"`
	DirectLeader           string       `db:"direct_leader,omitempty" json:"direct_leader"`
	DirectLeaderJobTitle   string       `db:"direct_leader_job_title,omitempty" json:"direct_leader_job_title"`
	DirectLeaderEmployeeID string       `db:"direct_leader_employee_id,omitempty" json:"direct_leader_employee_id"`
	PICHrbp                string       `db:"pic_hrbp,omitempty" json:"pic_hrbp"`
	HrbpPrivyID            string       `db:"hrbp_privy_id,omitempty" json:"hrbp_privy_id"`
	Role                   string       `db:"role,omitempty" json:"role"`
	LeadershipStatus       string       `db:"leadership_status,omitempty" json:"leadership_status"`
	CreatedAt              sql.NullTime `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt              sql.NullTime `db:"updated_at,omitempty" json:"updated_at"`
	DeletedAt              sql.NullTime `db:"deleted_at,omitempty" json:"deleted_at"`
}

type UserVersion struct {
	ID              int64  `db:"id,omitempty" json:"id"`
	Version         string `db:"version,omitempty" json:"version"`
	Ignore          int64  `db:"ignore,omitempty" json:"ignore"`
	UserCount       int64  `db:"user_count,omitempty" json:"user_count"`
	EnterpriseToken string `db:"enterprise_token,omitempty" json:"enterprise_token"`
}

type Enterprise struct {
	ID        int          `db:"id,omitempty" json:"id"`
	Name      string       `db:"name,omitempty" json:"name"`
	Token     string       `db:"token,omitempty" json:"token"`
	IsActive  bool         `db:"is_active,omitempty" json:"is_active"`
	PrivyID   string       `db:"privy_id,omitempty" json:"privy_id"`
	CreatedAt sql.NullTime `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at,omitempty" json:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at,omitempty" json:"deleted_at"`
}

type Event struct {
	ID          int64        `db:"id,omitempty" json:"id"`
	Name        string       `db:"name,omitempty" json:"name"`
	Description string       `db:"description,omitempty" json:"description"`
	UserVersion string       `db:"user_version,omitempty" json:"user_version"`
	Status      string       `db:"status,omitempty" json:"status"`
	Progress    float64      `gorm:"-"`
	StartDate   sql.NullTime `gorm:"-"`
	EndDate     sql.NullTime `gorm:"-"`
	CreatedAt   sql.NullTime `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at,omitempty" json:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at,omitempty" json:"deleted_at"`
}

type Session struct {
	ID          int64        `db:"id,omitempty" json:"id"`
	EventID     int64        `db:"event_id,omitempty" json:"event_id"`
	Type        string       `db:"type,omitempty" json:"type"`
	Name        string       `db:"name,omitempty" json:"name"`
	StartDate   time.Time    `db:"start_date,omitempty" json:"start_date"`
	EndDate     time.Time    `db:"end_date,omitempty" json:"end_date"`
	Description string       `db:"description,omitempty" json:"description"`
	Status      string       `db:"status,omitempty" json:"status"`
	CreatedAt   sql.NullTime `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at,omitempty" json:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at,omitempty" json:"deleted_at"`
}

type Question struct {
	ID        int64        `db:"id,omitempty" json:"id"`
	SessionID int64        `db:"session_id,omitempty" json:"session_id"`
	Sort      int          `db:"sort,omitempty" json:"sort"`
	Name      string       `db:"name,omitempty" json:"name"`
	Type      string       `db:"type,omitempty" json:"type"`
	IsDB      bool         `db:"is_db,omitempty" json:"is_db"`
	Option    string       `db:"option,omitempty" json:"option"`
	Required  bool         `db:"required,omitempty" json:"required"`
	Max       int          `db:"max,omitempty" json:"max"`
	CreatedAt sql.NullTime `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at,omitempty" json:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at,omitempty" json:"deleted_at"`
}

type FormTask struct {
	ID         int64        `db:"id,omitempty" json:"id"`
	SessionID  int64        `db:"session_id,omitempty" json:"session_id"`
	RevieweeID int64        `db:"reviewee_id,omitempty" json:"reviewee_id"`
	ReviewerID int64        `db:"reviewer_id,omitempty" json:"reviewer_id"`
	Status     string       `db:"status,omitempty" json:"status"`
	CreatedAt  sql.NullTime `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt  sql.NullTime `db:"updated_at,omitempty" json:"updated_at"`
	DeletedAt  sql.NullTime `db:"deleted_at,omitempty" json:"deleted_at"`
}

type FormScaleAnswer struct {
	ID         int64        `db:"id,omitempty" json:"id"`
	SessionID  int64        `db:"session_id,omitempty" json:"session_id"`
	RevieweeID int64        `db:"reviewee_id,omitempty" json:"reviewee_id"`
	ReviewerID int64        `db:"reviewer_id,omitempty" json:"reviewer_id"`
	QuestionID int64        `db:"question_id,omitempty" json:"question_id"`
	Sort       int          `db:"sort,omitempty" json:"sort"`
	ScaleValue float64      `db:"scale_value,omitempty" json:"scale_value"`
	Max        int          `db:"max,omitempty" json:"max"`
	CreatedAt  sql.NullTime `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt  sql.NullTime `db:"updated_at,omitempty" json:"updated_at"`
	DeletedAt  sql.NullTime `db:"deleted_at,omitempty" json:"deleted_at"`
}

type FormTextAnswer struct {
	ID         int64        `db:"id,omitempty" json:"id"`
	SessionID  int64        `db:"session_id,omitempty" json:"session_id"`
	RevieweeID int64        `db:"reviewee_id,omitempty" json:"reviewee_id"`
	ReviewerID int64        `db:"reviewer_id,omitempty" json:"reviewer_id"`
	QuestionID int64        `db:"question_id,omitempty" json:"question_id"`
	Sort       int          `db:"sort,omitempty" json:"sort"`
	TextValue  string       `db:"text_value,omitempty" json:"text_value"`
	CreatedAt  sql.NullTime `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt  sql.NullTime `db:"updated_at,omitempty" json:"updated_at"`
	DeletedAt  sql.NullTime `db:"deleted_at,omitempty" json:"deleted_at"`
}

type Report struct {
	RevieweeID    int64           `db:"reviewee_id" json:"reviewee_id"`
	Name          sql.NullString  `db:"name" json:"name"`
	SelfAssess    sql.NullFloat64 `db:"self_assess" json:"self_assess"`
	MemberAssess  sql.NullFloat64 `db:"member_assess" json:"member_assess"`
	PeersAssess   []float64       `db:"peer_assess" json:"peer_assess"`
	TotalReviewer int64           `db:"total_reviewer" json:"total_reviewer"`
	FinalScore    float64         `db:"final_score" json:"final_score"`
}

type Document struct {
	ID        int          `db:"id,omitempty" json:"id"`
	Title     string       `db:"title,omitempty" json:"title"`
	DocToken  string       `db:"doc_token,omitempty" json:"doc_token"`
	Owner     string       `db:"owner,omitempty" json:"owner"`
	CreatedAt sql.NullTime `db:"created_at,omitempty" json:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at,omitempty" json:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at,omitempty" json:"deleted_at"`
}
