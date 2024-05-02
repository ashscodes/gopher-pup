package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var err error
var isInMemory bool
var peopleCollection *mongo.Collection
var people People

func ConnectDatabase() {
	dbPath := "mongodb://localhost:27017"
	connectToMongoDB(dbPath)

	if !isInMemory {
		isDatabaseConnected()
		if isEmpty, err := isCollectionEmpty(); err != nil {
			log.Fatal(err)
		} else if isEmpty {
			fmt.Println("Found no documents in collection, seeding the database...")
			seedDatabase()
		}
	} else {
		seedDatabase()
	}
}

func GetAllPeople(query bson.M) ([]*Person, error) {
	var result []*Person
	if (isInMemory) {
		return people.ConvertToSlice(), nil
	} else {
		// Returning multiple docs returns a 'cursor' object, decode one by one.
		cursor, err := peopleCollection.Find(context.TODO(), query, nil)
		if err != nil {
			return nil, err
		}

		defer cursor.Close(context.TODO())

		for cursor.Next(context.TODO()) {
			var person Person
			err = cursor.Decode(&person)
			if err != nil {
				return nil, err
			}

			result = append(result, &person)
		}

		if err = cursor.Err(); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func GetPersonByObjectId(id string) (*Person, error) {
	if isInMemory {
		return nil, fmt.Errorf("data is in memory and this endpoint is not available")
	} else {
		var person *Person
		docId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}

		filter := bson.D{{Key: "_id", Value: docId}}
		err = peopleCollection.FindOne(context.TODO(), filter).Decode(&person)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return nil, fmt.Errorf("person not found")
			}

			return nil, err
		}

		return person, nil
	}
}

func connectToMongoDB(path string) {
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(path))
	if err != nil {
		log.Println("No mongo database available. Reverting to in-memory storage")
		isInMemory = true
		return
	}

	defer client.Disconnect(context.TODO())
	fmt.Printf("Connected to %v!\n", path)
}

func generateUsers() People {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	firstnames := [25]string{"John", "Emma", "Michael", "Sophia", "William", "Olivia", "James", "Ava", "Alexander", "Isabella",
		"Ethan", "Mia", "Daniel", "Charlotte", "Matthew", "Amelia", "Benjamin", "Harper", "Joseph", "Evelyn",
		"Andrew", "Abigail", "David", "Emily", "Christopher"}

	lastnames := [25]string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez",
		"Martinez", "Wilson", "Anderson", "Taylor", "Thomas", "Jackson", "White", "Harris", "Clark", "Lewis",
		"Robinson", "Walker", "Hall", "Young", "Allen"}

	cities := [25]string{"London", "Paris", "New York", "Tokyo", "Berlin", "Sydney", "Los Angeles", "Toronto",
	"Madrid", "Rome", "Moscow", "Beijing", "Dubai", "Singapore", "Hong Kong", "Mumbai", "Rio de Janeiro",
	"Cape Town", "Bangkok", "Seoul", "Amsterdam", "Stockholm", "Oslo", "Helsinki", "Vienna"}

	countries := [25]string{"UK", "France", "USA", "Japan", "Germany", "Australia", "Canada", "Spain", "Italy", "Russia",
		"China", "UAE", "Singapore", "Hong Kong", "India", "Brazil", "South Africa", "Thailand", "South Korea", "Netherlands",
		"Sweden", "Norway", "Finland", "Austria"}

		var users []Person
	for i := 0; i < 100; i++ {
		firstname := firstnames[r.Intn(len(firstnames))]
		lastname := lastnames[r.Intn(len(lastnames))]
		location := &Location{
			City:    cities[r.Intn(len(cities))],
			Country: countries[r.Intn(len(countries))],
		}
		users = append(users, Person{
			Firstname: firstname,
			Lastname:  lastname,
			Location:  location,
		})
	}
	return users
}

func isCollectionEmpty() (bool, error) {
	var count int64
	count, err = peopleCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return false, fmt.Errorf("error counting documents in collection: %v", err)
	}

	return count == 0, nil
}

func isDatabaseConnected() {
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func seedDatabase() {
	if isInMemory {
		people = generateUsers()
	} else {
		_, err := peopleCollection.InsertMany(context.TODO(), generateUsers().ConvertToInterface())
		if err != nil {
			log.Fatal(err)
		}
	}
}
