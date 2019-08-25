package main

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Event represents a saved event in the database
type Event struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Date        string        `json:"date"`
	Time        string        `json:"time"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	UserID      string        `json:"userID"`
}

type Events []Event

func IsValidEvent(e Event) error {
	if len(e.Title) == 0 {
		return errors.New("title is empty")
	}

	if len(e.Title) < 4 {
		return errors.New("title should be more than 4 characters")
	}

	if len(e.Description) == 0 {
		return errors.New("description is empty")
	}

	if len(e.Description) < 4 {
		return errors.New("description should be more than 4 characters")
	}

	t, err := time.Parse("2006-01-01", e.Date)
	if err != nil {
		return err
	}

	if time.Now().After(t) {
		return errors.New("You can only add an event going forward not backwards")
	}

	t, err = time.Parse(time.Kitchen, e.Time)
	if err != nil {
		return err
	}

	e.Time = t.String()
	e.Date = t.String()
	return nil
}

type Response struct {
	Message   string
	Timestamp int64
}
