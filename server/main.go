package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	mgo "gopkg.in/mgo.v2"
)

func main() {

	port := flag.Int("http.port", 2000, "Port to run HTTP server on")
	dsn := flag.String("store.dsn", "localhost:27017/calendar_app", "DSN to connect to MongoDB")

	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Could not load .env file ... %v", err)
	}

	sess, err := mgo.Dial(*dsn)
	if err != nil {
		log.Fatal(err)
	}

	db := &store{sess}

	mux := chi.NewMux()

	mux.Post("/login", login(db))

	mux.Group(func(rr chi.Router) {

		rr.Route("/events", func(r chi.Router) {

			r.Use(authenticateUser(db))

			r.Post("/add", addEvent(db))
			r.Get("/{id}", viewEvent(db))
			r.Delete("/{id}", deleteEvent(db))
		})

	})

	log.Printf("Running HTTP server on %d", *port)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func encode(w http.ResponseWriter, v interface{}) {
	json.NewEncoder(w).Encode(v)
}

func login(s *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var body struct {
			Email string `json:"email"`
		}

		type response struct {
			Timestamp int64
			Message   string
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encode(w, response{Message: "Invalid request body", Timestamp: time.Now().Unix()})

			return
		}

		if _, err := mail.ParseAddress(body.Email); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encode(w, response{Message: "Please provide your email address", Timestamp: time.Now().Unix()})

			return
		}

		user, err := s.FindOrCreateUser(body.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			encode(w, response{Message: "An error occurred while authenticating the user", Timestamp: time.Now().Unix()})
			return
		}

		var res = struct {
			*User
			Message   string
			Timestamp int64
		}{
			User:      user,
			Message:   "Successful authenticated user",
			Timestamp: time.Now().Unix(),
		}

		w.WriteHeader(http.StatusOK)
		encode(w, res)
	}
}

func addEvent(s *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var e Event
		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encode(w, Response{
				Message:   "Invalid request body",
				Timestamp: time.Now().Unix(),
			})
			return
		}

		if err := IsValidEvent(e); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			encode(w, Response{
				Message:   err.Error(),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		if err := s.StoreEvent(userFromContext(r.Context()), e); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			encode(w, Response{
				Message:   "An error occurred while saving your event",
				Timestamp: time.Now().Unix(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		encode(w, Response{
			Message:   "Your event was successfully added",
			Timestamp: time.Now().Unix(),
		})

	}
}

func deleteEvent(s *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func viewEvent(s *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")

		event, err := s.FindEvent(id)
		if err != nil {
			var status = http.StatusInternalServerError
			var message = "An error occurred while fetching this event"

			if err == mgo.ErrNotFound {
				status = http.StatusNotFound
				message = "Event does not exists"
			}

			w.WriteHeader(status)
			encode(w, Response{
				Message:   message,
				Timestamp: time.Now().Unix(),
			})
			return
		}

		if event.UserID != userFromContext(r.Context()).ID.Hex() {
			w.WriteHeader(http.StatusUnauthorized)
			encode(w, Response{
				Message:   "You don't have permission to view this",
				Timestamp: time.Now().Unix(),
			})
			return
		}

		encode(w, event)
	}
}
