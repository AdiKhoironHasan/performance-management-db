package entity

import (
	"time"

	"gorm.io/gorm"
)

type Versioning struct {
	ID      string `gorm:"primary_key"`
	Version int    `gorm:"not null"`
	Offset  int    `gorm:"not null"`
}

type User struct {
	ID                     int        `gorm:"primaryKey"`
	Version                string     `gorm:"not null"`
	Name                   string     `gorm:"not null"`
	PrivyID                string     `gorm:"not null"`
	Email                  string     `gorm:"not null"`
	Status                 string     `gorm:"not null"`
	JoinDate               time.Time  `gorm:"not null"`
	JobTitle               string     `gorm:"not null"`
	Level                  string     `gorm:"not null"`
	Directorate            string     `gorm:"not null"`
	Division               string     `gorm:"not null"`
	Homebase               string     `gorm:"not null"`
	DirectLeader           string     `gorm:"not null"`
	DirectLeaderJobTitle   string     `gorm:"not null"`
	DirectLeaderEmployeeID string     `gorm:"not null"`
	PICHrbp                string     `gorm:"not null"`
	HrbpPrivyID            string     `gorm:"not null"`
	Role                   string     `gorm:"not null"`
	CreatedAt              time.Time  `gorm:"column:created_at"`
	UpdatedAt              time.Time  `gorm:"column:updated_at"`
	DeletedAt              *time.Time `gorm:"column:deleted_at"`
}

type Event struct {
	ID          int    `gorm:"primary_key"`
	Name        string `gorm:"not null"`
	Description string
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
	Sessions    []Session  `gorm:"foreignkey:EventID"`
}

type Session struct {
	ID          int       `gorm:"primary_key"`
	EventID     int       `gorm:"not null"`
	Type        string    `gorm:"not null"`
	Name        string    `gorm:"not null"`
	StartDate   time.Time `gorm:"not null"`
	EndDate     time.Time `gorm:"not null"`
	Description string
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
	Questions   []Question `gorm:"foreignkey:SessionID"`
}

type Question struct {
	ID        int    `gorm:"primary_key"`
	SessionID int    `gorm:"not null"`
	Sort      int    `gorm:"not null"`
	Name      string `gorm:"not null"`
	Type      string `gorm:"not null"`
	Option    string `gorm:"not null"`
	Required  string `gorm:"not null"`
	Max       int
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

type QuestionAnswer struct {
	ID         int    `gorm:"primary_key"`
	SessionID  int    `gorm:"not null"`
	RevieweeID int    `gorm:"not null"`
	ReviewerID int    `gorm:"not null"`
	QuestionID int    `gorm:"not null"`
	Status     string `gorm:"not null"`
	Scale      int    `gorm:"type:integer"`
	Essay      string `gorm:"type:text"`
	Dropdown   string
	Radio      int
	CreatedAt  time.Time  `gorm:"column:created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at"`

	Type string `gorm:"column:type"`
}

type FormScaleAnswer struct {
	ID         uint `gorm:"primaryKey"`
	SessionID  int
	RevieweeID int
	ReviewerID int
	QuestionID int
	Sort       int
	ScaleValue int `gorm:"column:scale_value;note:maximum value by max"`
	Max        int `gorm:"note:maximum value of type scale"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// Struct for the "form_answer_text" table
type FormAnswerText struct {
	ID         uint `gorm:"primaryKey"`
	SessionID  int
	RevieweeID int
	ReviewerID int
	QuestionID int
	Sort       int
	TextValue  string `gorm:"column:text_value"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
