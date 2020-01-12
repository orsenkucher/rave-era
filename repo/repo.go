package repo

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
	Events []Event
	Users  []Raver
}

//NewRepo is public
func NewRepo() *Repo {
	r := &Repo{}
	r.startUp()
	return r
}

func (r *Repo) startUp() {
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

	r.load()
	r.loadRavers()
	fmt.Println("Startup")
}

func (r *Repo) loadRavers() {
	iter := r.client.Collection("Users").Documents(r.ctx)
	docs, _ := iter.GetAll()

	for _, doc := range docs {
		var user Raver
		doc.DataTo(&user)
		fmt.Println(user)
		r.Users = append(r.Users, user)
	}
}

func (r *Repo) load() {
	iter := r.client.Collection("Events").Documents(r.ctx)
	docs, _ := iter.GetAll()

	for _, doc := range docs {
		ev := struct {
			N    int    `firebase:"n"`
			Name string `firebase:"name"`
		}{}
		doc.DataTo(&ev)
		fmt.Println(ev)
		event := Event{N: ev.N, Name: ev.Name}
		event.Free.Ravers = make(map[Raver][]Raver)
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
		r.Events = append(r.Events, event)
	}
}

//UserRegistrated UserRegistrated
func (r *Repo) UserRegistrated(id int64) bool {
	for _, user := range r.Users {
		if id == user.ID {
			return true
		}
	}
	return false
}

func (r *Repo) Find(id int64) (Raver, bool) {
	for _, user := range r.Users {
		if id == user.ID {
			return user, true
		}
	}
	return Raver{}, false
}

//AddUser is public
func (r *Repo) AddUser(user Raver) {
	r.Users = append(r.Users, user)
	r.client.Doc("Users/"+strconv.FormatInt(user.ID, 10)).Set(r.ctx, user)
}

//Subscribe Subscribe
func (r *Repo) Subscribe(user Raver, friends []Raver, eventName string) {
	for _, ev := range r.Events {
		if ev.Name == eventName {
			ev.Free.Ravers[user] = friends
			r.client.Doc("Events/"+ev.Name+"/Free/"+strconv.FormatInt(user.ID, 10)).Set(r.ctx, user)
			fmt.Print("Events/" + ev.Name + "/Free/" + strconv.FormatInt(user.ID, 10))
			for _, friend := range friends {
				r.client.Doc("Events/"+ev.Name+"/Free/"+strconv.FormatInt(user.ID, 10)+"/Friends/"+strconv.FormatInt(friend.ID, 10)).Set(r.ctx, friend)
			}
		}
	}
}

//Save is public
func (r *Repo) Save() {
	for _, event := range r.Events {
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
	for _, event := range r.Events {
		fmt.Print(event.Name)
		fmt.Print()
	}
}
