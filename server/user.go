package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Email       string        `json:"email"`
	CreatedAt   time.Time     `json:"createdAt"`
	AccessToken string        `json:"accessToken"`
}

func (u *User) CreateAccessToken() {
	if u.AccessToken == "" {
		u.AccessToken = uuid.New().String()
	}
}

type ContextKey string

const userContextKey ContextKey = "User"

func authenticateUser(db *store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := r.Header.Get("Authorization")
			token = strings.TrimPrefix(token, "Bearer ")

			user, err := db.FindUserByAccessToken(token)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})

	}
}

func userFromContext(ctx context.Context) *User {
	return ctx.Value(userContextKey).(*User)
}
