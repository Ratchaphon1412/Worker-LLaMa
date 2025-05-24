package entities

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

const (
	STATUS_PENDING  string = "pending"
	STATUS_COMPLETE string = "complete"
	STATUS_ERROR    string = "error"
	STATUS_PROCESS  string = "process"
)

type Chat struct {
	gorm.Model
	WorkflowID  string     `json:"workflow_id"`
	Prompt      string     `json:"prompt"`
	Answer      string     `json:"answer"`
	AnswerMedia string     `json:"answer_media"`
	Status      string     `json:"status"`
	Researches  []Research `gorm:"constraint:OnUpdate:CASCADE"`
}

type Research struct {
	gorm.Model
	Kind        string      `json:"kind"`
	Title       string      `json:"title"`
	DisplayLink string      `json:"display_link"`
	Link        string      `json:"link"`
	Image       string      `json:"image"`
	ChatID      uint        `json:"chat_id"`
	Thumbnails  []Thumbnail `gorm:"constraint:OnUpdate:CASCADE"`
}
type Thumbnail struct {
	gorm.Model
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Src        string `json:"src"`
	ResearchID uint   `json:"research_id"`
}

func (c *Chat) BeforeCreate(tx *gorm.DB) (err error) {
	/* Before create hook to set the WorkflowID and Status */
	c.WorkflowID = "workflow_id_" + uuid.NewV4().String()
	c.Status = STATUS_PENDING
	return
}
