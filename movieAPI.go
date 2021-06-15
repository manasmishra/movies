package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// DB stores the database session information needs to be initialized once
type DB struct {
	session *mgo.Session
	collection *mgo.Collection
}

type Movie struct {
	ID	bson.ObjectId `json:"id" bson:"_id", omitempty`
	Name string `json:"name" bson:"name"`
	Year string `json:"year" bson:"year"`
	Directors []string `json:"directors" bson:"directors"`
	Writers []string `json:"writers" bson:"writers"`
	BoxOffice BoxOffice `json:"boxOffice" bson:"boxOffice"`
}

type BoxOffice struct {
	Budget uint64
	Gross uint64
}

// Get Movie fetches a movie with a given Id
func (db *DB) GetMovie(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	var movie Movie
	err := db.collection.Find(bson.M{"id": bson.ObjectIdHex(vars["id"])}).One(&movie)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(movie)
		w.Write(response)
	}
}

// POST Movie adds a new movie to our MongoDB collection
func (db *DB) PostMovie(w http.ResponseWriter, r *http.Request)  {
	var movie Movie
	postBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(postBody, &movie)

	// Create a hash ID to insert
	movie.ID = bson.NewObjectId()
	err := db.collection.Insert(movie)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(movie)
		w.Write(response)
	}
}

// update movie
func (db *DB) UpdateMovie(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	var movie Movie
	putBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(putBody, &movie)
	err := db.collection.Update(bson.M{"_id": bson.ObjectIdHex(vars["id"])}, bson.M{"$set": &movie})
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text")
		w.Write([]byte("Uploaded Successfully!"))
	}
}

func (db *DB) DeleteMovie(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	err := db.collection.Remove(bson.M{"_id": bson.ObjectIdHex(vars["id"])})
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "Text")
		w.Write([]byte("Deleted Successfully!!!"))
	}
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	c := session.DB("appdb").C("movies")
	db := &DB{session: session, collection: c}
	// create a router
	r := mux.NewRouter()
	// Attach a path with handler to the created router
	r.HandleFunc("/v1/movies/{id:[a-z0-9A-Z]*}", db.GetMovie).Methods("GET")
	r.HandleFunc("/v1/movies", db.PostMovie).Methods("POST")
	r.HandleFunc("/v1/movies/{id:[a-z0-9A-Z]*}", db.UpdateMovie).Methods("PUT")
	r.HandleFunc("/v1/movies/{id:[a-z0-9A-Z]*}", db.DeleteMovie).Methods("DELETE")
	srv := &http.Server{
		Handler: r,
		Addr: "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout: 15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
