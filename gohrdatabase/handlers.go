package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// CreatePerson handles the HTTP POST request to create a new person record.
// It decodes the JSON request body into a Person struct, creates the record
// in the database, and returns the result as a JSON response.
func CreatePerson(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var person Person
	_ = json.NewDecoder(req.Body).Decode(&person)

	result, err := CreatePersonRecord(person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// DeletePerson handles the HTTP DELETE request to delete a person record by ID.
// It retrieves the person ID from the request parameters, deletes the record
// from the database, and returns the result as a JSON response.
func DeletePerson(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(req)
	id := params["id"]

	result, err := DeletePersonRecord(id)
	if err != nil {
		if err.Error() == "person not found" {
			http.Error(w, "Person not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// GetPeople handles the HTTP GET request to retrieve multiple person records
// based on query parameters.
// It parses the query parameters, constructs MongoDB filter criteria, retrieves
// matching records from the database, and returns them as a JSON response.
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

// GetPerson handles the HTTP GET request to retrieve a single person record by ID.
// It retrieves the person ID from the request parameters, fetches the record
// from the database, and returns it as a JSON response.
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

// UpdatePerson handles the HTTP PUT request to update an existing person record.
// It decodes the JSON request body into a Person struct, retrieves the person ID
// from the request parameters, updates the record in the database, and returns
// the result as a JSON response.
func UpdatePerson(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var person Person
	_ = json.NewDecoder(req.Body).Decode(&person)

	params := mux.Vars(req)
	id := params["id"]

	result, err := UpdatePersonRecord(person, id)
	if err != nil {
		if err.Error() == "person not found" {
			http.Error(w, "Person not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func getPeopleQueryFilter() []QueryFilter {
	return []QueryFilter{
		{Name: "firstname"},
		{Name: "lastname"},
		{Name: "city", ParentPath: "location"},
		{Name: "country", ParentPath: "location"},
	}
}

// parseQuery parses the query parameters from an HTTP request into MongoDB filter criteria.
// It takes the URL query values and a list of QueryFilter structs, constructs filter criteria
// based on the query parameters, and returns a BSON filter document suitable for MongoDB queries.
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
