package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

//Raver is user
type Raver struct {
	Name     string `firebase:"name"`
	LastName string `firebase:"lastname"`
	Uni      string `firebase:"uni"`
	ID       int64  `firebase:"id"`
	Age      int    `firebase:"age"`
}

//FreeList is Event with freelist
type FreeList struct {
	Ravers map[Raver][]Raver
}

//Event is event
type Event struct {
	N     int    `firebase:"n"`
	Name  string `firebase:"name"`
	Users []Raver
	Free  FreeList
}

//Repo is public
type Repo struct {
	client *firestore.Client
	ctx    context.Context
	events []Event
}

//StartUp starts server
func (r *Repo) StartUp() {
	r.ctx = context.Background()
	sa := option.WithCredentialsFile("creds.json")
	app, err := firebase.NewApp(r.ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	r.client, err = app.Firestore(r.ctx)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Startup")
}

//LoadEvents is public
func (r *Repo) LoadEvents() []FreeList {
	iter := r.client.Collection("Events").Documents(r.ctx)
	docs, _ := iter.GetAll()

	for _, doc := range docs {
		ev := struct {
			N    int    `firebase:"n"`
			Name string `firebase:"name"`
		}{}
		doc.DataTo(&ev)
		event := Event{N: ev.N, Name: ev.Name}
		iter := doc.Ref.Collection("Free").Documents(r.ctx)
		docs, _ := iter.GetAll()
		for _, doc := range docs {
			raver := Raver{}
			doc.DataTo(&raver)
			event.Free.Ravers[raver] = []Raver{}
			iter := doc.Ref.Collection("Friends").Documents(r.ctx)
			docs, _ := iter.GetAll()
			for _, doc := range docs {
				friend := Raver{}
				doc.DataTo(&friend)
				event.Free.Ravers[raver] = append(event.Free.Ravers[raver], friend)
			}
		}
		r.events = append(r.events, event)
	}
	return []FreeList{}
}

//SaveEvents is public
func (r *Repo) SaveEvents() {
	for _, event := range r.events {
		ev := struct {
			N    int    `firebase:"n"`
			Name string `firebase:"name"`
		}{
			N:    event.N,
			Name: event.Name,
		}
		r.client.Doc("Events/"+ev.Name).Set(r.ctx, ev)
		for raver, friends := range event.Free.Ravers {
			r.client.Doc("Events/"+ev.Name+"/Free/"+strconv.FormatInt(raver.ID, 10)).Set(r.ctx, raver)
			for _, friend := range friends {
				r.client.Doc("Events/"+ev.Name+"/Free/"+strconv.FormatInt(raver.ID, 10)+"/Friends/"+strconv.FormatInt(friend.ID, 10)).Set(r.ctx, friend)
			}
		}
	}
}

//ShowStatus is public
func (r *Repo) ShowStatus() {
	for _, event := range r.events {
		fmt.Print(event.Name)
		fmt.Print()
	}
}

func main() {
	s := Repo{}
	s.StartUp()
}
