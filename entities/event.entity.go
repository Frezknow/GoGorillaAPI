package entities

import "fmt"

type Event struct {
	ID          int `gorm:"primary_key, AUTO_INCREMENT"`
	Title       string
	Description string
	Comments    []Comment `gorm:"ForeignKey:EventID"`
}

func (event *Event) TableName() string {
	return "events"
}

func (event Event) ToString() string {
	return fmt.Sprintf("id: %d\ntitle: %s \ndescription: %s", event.ID, event.Title, event.Description)
}
