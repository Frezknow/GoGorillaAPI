package entities

import "fmt"

type Comment struct {
	Id      int `gorm:"primary_key, AUTO_INCREMENT"`
	Body    string
	EventID int `gorm:"column:event_id"`
	Event   Event
}

func (comment *Comment) TableName() string {
	return "comments"
}

func (comment Comment) ToString() string {
	return fmt.Sprintf("id: %d\nbody: %s\neventId: %d", comment.Id, comment.Body, comment.EventID)
}
