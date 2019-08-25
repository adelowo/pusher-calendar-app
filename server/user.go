package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"createdAt"`
	AccessToken string    `json:"accessToken"`
}

func (u *User) CreateAccessToken() {
	if u.AccessToken == "" {
		u.AccessToken = uuid.New().String()
	}
}
