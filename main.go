package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
)


// Person is a struct that represents a person in this application
type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

// Song is a struct that represents a single song
type Song struct {
	Title 	 string `json:"title"`
	Duration string `json:"duration"`
	Singer 	 Person `json:"singer"`
}

// songs is an array of Song
var songs []Song = []Song{}


var (
	tpl *template.Template
	tplName   = "index.gohtml" //template listed at the end.
	finalHTML *bytes.Buffer  //an io.Writer for template.Execute()
	gotForm   = true
	noForm    = false
)

func init() {
	tpl = template.Must(template.ParseFiles(tplName))
	finalHTML = bytes.NewBuffer([]byte(""))
}


func main(){
	router := mux.NewRouter()

	router.HandleFunc("/songs", addItem).Methods("POST")
	router.HandleFunc("/songs", getAllSongs).Methods("GET")
	router.HandleFunc("/songs/{id}", getSong).Methods("GET")
	router.HandleFunc("/songs/{id}", updateSong).Methods("PUT")
	router.HandleFunc("/songs/{id}", patchSong).Methods("PATCH")
	router.HandleFunc("/songs/{id}", deleteSong).Methods("DELETE")

	http.ListenAndServe(":8080", router)
}


func deleteSong(w http.ResponseWriter, req *http.Request){
	//get ID of the song from the root parameter
	var idParam string = mux.Vars(req)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID can't be converted to an int"))
		return
	}

	//err checking
	if id >= len(songs) {
		w.WriteHeader(404)
		w.Write([]byte("ID is out of range"))
		return
	}

	//Delete the song from the slice
	songs = append(songs[:id], songs[id+1:]...)

	w.WriteHeader(200)
}

func patchSong(w http.ResponseWriter, req *http.Request){
	//get ID of the song from the root parameter
	var idParam string = mux.Vars(req)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID can't be converted to an int"))
		return
	}

	//err checking
	if id >= len(songs) {
		w.WriteHeader(404)
		w.Write([]byte("ID is out of range"))
		return
	}

	//get the current value
	song := &songs[id]
	json.NewDecoder(req.Body).Decode(song)

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(*song)
}


func updateSong(w http.ResponseWriter, req *http.Request){
	//get ID of the song from the root parameter
	var idParam string = mux.Vars(req)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID can't be converted to an int"))
		return
	}

	//err checking
	if id >= len(songs) {
		w.WriteHeader(404)
		w.Write([]byte("ID is out of range"))
		return
	}

	//get the value from JSON body
	var updatedSong Song
	json.NewDecoder(req.Body).Decode(&updatedSong)

	songs[id] = updatedSong

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(updatedSong)
}


func getSong(w http.ResponseWriter, req *http.Request){
	//get the ID of the song from the route parameter
	var idParam string = mux.Vars(req)["id"]
	id, err :=strconv.Atoi(idParam)
	if err != nil{
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converter to integer"))
		return
	}

	//err checking
	//id == index into a specific database
	if id >= len(songs) {
		w.WriteHeader(404)
		w.Write([]byte("No song found with specified ID"))
		return
	}
	song := songs[id]

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(song)
}

func getAllSongs(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(songs)
}


func addItem(w http.ResponseWriter, req *http.Request){
	//routeVariable := mux.Vars(req)["item"]

	// get item value from JSON body
	var newSong Song
	json.NewDecoder(req.Body).Decode(&newSong)

	songs = append(songs, newSong)

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(songs)
}
