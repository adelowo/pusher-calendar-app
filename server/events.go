package main

import (
	"errors"
	"time"
)

// Event represents a saved event in the database
type Event struct {
	ID          int64  `json:"id"`
	Date        string `json:"date"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

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

	t, err := time.Parse("01-01-2006", e.Date)
	if err != nil {
		return err
	}

	e.Date = t.String()
	return nil
}

type Response struct {
	Message   string
	Timestamp int64
}
