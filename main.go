package main

import (
    "fmt"
	"log"
	"encoding/json"
    "math/rand"
    "net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type Movie struct {
	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
    Firstname string `json:"firstname"`
    Lastname string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if(item.ID == params["id"]){
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie (w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(10000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request){
	// set content type
	w.Header().Set("Content-Type", "application/json")
	// take params
	params := mux.Vars(r)
	// loop and take the movie by id
	for index, item := range movies {
		if item.ID == params["id"] {
			// delete old movie
			movies = append(movies[:index], movies[index+1:]...)
			
			// create new movie w what w send
			var movie Movie
			// decode body to movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
	// assign new values
}

func main(){
	r := mux.NewRouter()

	movies = append(
		movies, 
		Movie{
			ID: "1", 
			Isbn: "438227", 
			Title: "Interstellar", 
			// & untuk memberi address / * untuk mengakses address / pointer 
			Director: &Director {
				Firstname: "John",
				Lastname: "Doe",
			},
		},
	)

	movies = append(
		movies, 
		Movie{
			ID: "2", 
			Isbn: "45455", 
			Title: "The Matrix", 
			// & untuk memberi address / * untuk mengakses address / pointer 
			Director: &Director {
				Firstname: "Steve",
				Lastname: "Smith",
			},
		},
	)

	port := "8081"

	r.HandleFunc("/api/movies", getMovies).Methods("GET")
	r.HandleFunc("/api/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/api/movies", createMovie).Methods("POST")
	r.HandleFunc("/api/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/api/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port: %s\n", port)
	log.Fatal(http.ListenAndServe("localhost:"+port, r))

}
