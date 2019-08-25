package main

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type store struct {
	inner *mgo.Session
}

func (db *store) StoreEvent(u *User, e Event) error {

	sess := db.inner.Copy()
	defer sess.Close()

	return nil
}

func (db *store) DeleteEvent(e Event) error {
	return nil
}

func (db *store) FindOrCreateUser(email string) (*User, error) {

	u := new(User)

	sess := db.inner.Copy()
	defer sess.Close()

	err := sess.DB("").C("users").Find(bson.M{"email": email}).One(&u)
	if err != nil {
		if err == mgo.ErrNotFound {
			u.Email = email
			u.CreatedAt = time.Now()
			u.CreateAccessToken()

			if err := sess.DB("").C("users").Insert(u); err != nil {
				return nil, err
			}

			return u, nil
		}

		return nil, err
	}

	return u, nil
}

func (db *store) FindUserByAccessToken(token string) (*User, error) {
	u := new(User)

	sess := db.inner.Copy()
	defer sess.Close()

	return u, sess.DB("").C("users").Find(bson.M{"accesstoken": token}).One(&u)
}
