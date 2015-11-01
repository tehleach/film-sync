package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	path := "/media/kyle/01CF57AD7D315EE0/FRANKENSTEIN/Movies"
	films := getFilmNames(path)
	for i := 0; i < len(films); i++ {
		film := buildFilmFromName(films[i])
		addFilmToDb(session, film)
	}
}

//Film represents a film
type Film struct {
	Name       string `json:"Title"`
	Director   string
	Runtime    string
	Year       int `json:",string"`
	Rated      string
	Metascore  int     `json:",string"`
	IMDBRating float64 `json:"imdbRating,string"`
}

func getFilmNames(dir string) []string {
	files, _ := ioutil.ReadDir(dir)
	var films []string
	for _, f := range files {
		name := strings.TrimSpace(f.Name())
		if len(name) > 0 {
			films = append(films, name)
		}
	}
	return films
}

func buildFilmFromName(name string) *Film {
	name = strings.Split(name, ".")[0]
	name = strings.Replace(name, " ", "+", -1)
	url := fmt.Sprintf("http://www.omdbapi.com/?t=%v", name)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var film Film
	err = json.NewDecoder(resp.Body).Decode(&film)

	return &film
}

func addFilmToDb(session *mgo.Session, film *Film) {
	c := session.DB("film-api").C("films")
	fmt.Println("Writing: ", film)
	_, err := c.Upsert(bson.M{"name": film.Name}, film)
	if err != nil {
		panic(err)
	}
}
