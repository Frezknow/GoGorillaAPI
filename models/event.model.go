package models

import (
	"gorilla_api/config"
	"gorilla_api/entities"
)

type EventModel struct {
}

func (eventModel EventModel) FindAll() ([]entities.Event, error) {
	db, err := config.GetDB()
	if err != nil {
		return nil, err
	} else {
		var events []entities.Event
		db.Preload("Comments").Find(&events)
		return events, nil
	}

}
