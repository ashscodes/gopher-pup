package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func GetPeople(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	filter := bson.M{}
	query := req.URL.Query()
	if len(query) != 0 {
		firstname, found := query["firstname"]
		if found {
			var names []string
			if len(firstname) == 1 {
				names = strings.Split(firstname[0], ",")
			} else {
				names = firstname
			}
			fnQuery := bson.M{"$in": names}
			filter["firstname"] = fnQuery
		}

		lastname, found := query["lastname"]
		if found {
			var names []string
			if len(lastname) == 1 {
				names = strings.Split(lastname[0], ",")
			} else {
				names = lastname
			}
			lnQuery := bson.M{"$in": names}
			filter["lastname"] = lnQuery
		}

		city, found := query["city"]
		if found {
			cityQuery := bson.M{"$in": city}
			filter["location.city"] = cityQuery
		}

		country, found := query["country"]
		if found {
			countryQuery := bson.M{"$in": country}
			filter["location.country"] = countryQuery
		}
	}

	people, err := GetAllPeople(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(people) == 0 {
		http.Error(w, "No people found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(people)
}

func GetPerson(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	id := params["id"]
	person, err := GetPersonByObjectId(id)
	if err != nil {
		if err.Error() == "person not found" {
			http.Error(w, "Person not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(person)
}