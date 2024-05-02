package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func GetPeople(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryFilters := getPeopleQueryFilter()
	filter := parseQuery(req.URL.Query(), queryFilters)
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

func getPeopleQueryFilter() []QueryFilter {
	return []QueryFilter{
		{Name: "firstname"},
		{Name: "lastname"},
		{Name: "city", ParentPath: "location"},
		{Name: "country", ParentPath: "location"},
	}
}

func parseQuery(queryValues url.Values, queryFilters []QueryFilter) bson.M {
	filter := bson.M{}
	if len(queryValues) != 0 {
		for _, queryFilter := range queryFilters {
			item, found := queryValues[queryFilter.Name]
			if found {
				var names []string
				if len(item) == 1 {
					names = strings.Split(item[0], ",")
				} else {
					names = item
				}

				jsonPath := queryFilter.Name
				if queryFilter.ParentPath != "" {
					jsonPath = queryFilter.ParentPath + "." + jsonPath
				}

				primitive := bson.M{"$in": names}
				filter[jsonPath] = primitive
			}
		}
	}

	return filter
}
