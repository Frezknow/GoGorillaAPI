package entities

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type Event struct {
	gorm.Model
	ID          int `gorm:"primary_key, AUTO_INCREMENT"`
	Title       string
	Description string
	Comments    []Comment //`gorm:"foreignKey:EventID"`
}

func (event *Event) TableName() string {
	return "events"
}

func (event Event) ToString() string {
	return fmt.Sprintf("id: %d\ntitle: %s \ndescription: %s", event.ID, event.Title, event.Description)
}
