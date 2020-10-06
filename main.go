package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Song is a struct that represents a single song
type Song struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Duration string `json:"duration"`
	Singer   string `json:"singer"`
}

var (
	// lenSong will store the len(songs). songs is an array of Song (songs []Song)
	lenSong int
	db      *sql.DB
	err     error
)

func main() {

	// Open the DB: "testAPI" with the username "root"
	db, err = sql.Open("mysql", "root:@/testAPI")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// NewRouter returns a new router instance.
	router := mux.NewRouter()

	// HandleFunc registers a new route with a matcher for the URL path.
	// Methods adds a matcher for HTTP methods.
	// It accepts a sequence of one or more methods to be matched (eg: GET, POST, PUT, PATCH, DELETE ...)
	router.HandleFunc("/", general).Methods("GET")
	router.HandleFunc("/songs", addSong).Methods("POST")
	router.HandleFunc("/songs", getAllSongs).Methods("GET")
	router.HandleFunc("/songs/{id}", getSong).Methods("GET")
	router.HandleFunc("/songs/{id}", updateSong).Methods("PUT")
	router.HandleFunc("/songs/{id}", deleteSong).Methods("DELETE")
	router.HandleFunc("/songss/delAll", deleteAll).Methods("GET")

	// starting the localhost:8080
	http.ListenAndServe(":8080", router)
}

// This function is the main pattern: "http://localhost:8080/"
func general(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte (`<h1>You are in the main page, please add pattern:<br> - "/songs"         -> to see all songs <br>
																				  - "/songs/{id}"    -> to see a specific song <br>
																				  - "/songss/delAll" -> to delete all songs stored </h1>`))
}

// This function add a song with method POST from JSON body
func addSong(w http.ResponseWriter, req *http.Request) {
	// stmt = statement, is a prepared statement
	// Prepare -> Execute, are used for actions which DOESN'T returns any rows
	stmt, err := db.Prepare("INSERT INTO songs(title, duration, singer) VALUES (?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}

	// body transform the data stored in req.body JSON to []byte
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err.Error())
	}

	// Extract data from req.body in keyVal
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["title"]
	duration := keyVal["duration"]
	singer := keyVal["singer"]

	// Execute the prepared syntax with the values set into req.body
	_, err = stmt.Exec(title, duration, singer)
	if err != nil {
		panic(err.Error())
	}

	// This message will print into ResponeWriter
	fmt.Fprintf(w, "New song was added")
}

// This function get all songs stored into DB songs
func getAllSongs(w http.ResponseWriter, req *http.Request) {
	// songs is an array of Song
	var songs []Song
	// Query is uesed for actions wihch returns rows
	// Store all records in result from DB songs, if any
	result, err := db.Query("SELECT * from songs")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	// Checking each record stored in result and append it to songs
	for result.Next() {
		var song Song
		err := result.Scan(&song.ID, &song.Title, &song.Duration, &song.Singer)
		if err != nil {
			panic(err.Error())
		}
		songs = append(songs, song)
	}
	lenSong = len(songs)

	// setting the header “Content-Type” to “application/json”.
	w.Header().Set("Content-type", "application/json")
	// encode songs to JSON and send them to interface
	json.NewEncoder(w).Encode(songs)
}

// This function get a specified song with method POST from JSON body
func getSong(w http.ResponseWriter, req *http.Request) {
	// params is a map from route paramete
	params := mux.Vars(req)

	// result store the record that satisfies the condition from Query
	result, err := db.Query("SELECT * FROM songs WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	// convert the song ID into an int
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		// 400 indicates server was unable to process the request sent by the client due to invalid syntax
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to int"))
		return
	}

	// checking if id exists in my DB
	if (id > lenSong) || (id < 1) {
		// 404 indicates page not found
		w.WriteHeader(404)
		w.Write([]byte("No song found with specified ID"))
		return
	}

	var song Song

	// Scan copies the columns in the current row into the values pointed
	// at by dest. The number of values in dest must be the same as the
	// number of columns in Rows.
	for result.Next() {
		err := result.Scan(&song.ID, &song.Title, &song.Duration, &song.Singer)
		if err != nil {
			panic(err.Error())
		}
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(song)
}

func updateSong(w http.ResponseWriter, req *http.Request) {
	// params is a map from route parameter
	params := mux.Vars(req)

	// convert the song ID into an int
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to int"))
		return
	}

	// checking if id exists in my DB
	if (id > lenSong) || (id < 1) {
		w.WriteHeader(404)
		w.Write([]byte("ID is out of range"))
		return
	}

	stmt, err := db.Prepare("UPDATE songs SET title = ?, duration=?, singer =? WHERE id= ?")
	if err != nil {
		panic(err.Error())
	}

	// body transform the data stored in req.body JSON to []byte
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err.Error())
	}

	// Extract data from req.body in keyVal
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newTitle := keyVal["title"]
	newDuration := keyVal["duration"]
	newSinger := keyVal["singer"]

	// Execute the prepared syntax with the values set into req.body
	_, err = stmt.Exec(newTitle, newDuration, newSinger, params["id"])
	if err != nil {
		panic(err.Error())
	}

	// This message will print into ResponeWriter
	fmt.Fprintf(w, "Song with ID = %s was updated", params["id"])
}

func deleteSong(w http.ResponseWriter, req *http.Request) {
	// params is a map from route parameter
	params := mux.Vars(req)

	// convert the song ID into an int
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to int"))
		return
	}

	// checking if id exists in my DB
	if (id > lenSong) || (id < 1) {
		w.WriteHeader(404)
		w.Write([]byte("ID is out of range"))
		return
	}

	stmt, err := db.Prepare("DELETE FROM songs WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}

	// 200 - OK success status response code indicates that the request has succeeded
	w.WriteHeader(200)
	fmt.Fprintf(w, "Post with ID = %s was deleted", params["id"])
}

func deleteAll(w http.ResponseWriter, req *http.Request) {
	stmt, err := db.Query("TRUNCATE TABLE songs")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	// 200 - OK success status response code indicates that the request has succeeded
	w.WriteHeader(200)

	// This message will print into ResponeWriter
	w.Write([]byte("All records from table songs have been deleted"))
}
