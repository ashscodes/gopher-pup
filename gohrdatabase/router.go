package main

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/person", GetPeople).Methods("GET")
	router.HandleFunc("/person/{id}", GetPerson).Methods("GET")
	// router.HandleFunc("/person", CreatePerson).Methods("POST")
	// router.HandleFunc("/person/{id}", PatchPerson).Methods("PATCH")
	// router.HandleFunc("/person/{id}", UpdatePerson).Methods("PUT")
	// router.HandleFunc("/person/{id}", DeletePerson).Methods("DELETE")
	return router
}
