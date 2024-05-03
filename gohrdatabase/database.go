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
var personMap map[string]Person

// ConnectDatabase establishes a connection to the database and initializes it.
// If the application is running in memory mode, it seeds the database with mock data.
// If the application is using MongoDB, it checks if the collection is empty and seeds the database if necessary.
func ConnectDatabase() {
	dbPath := "mongodb://localhost:27017"
	connectToMongoDB(dbPath)

	if isInMemory {
		err = seedDatabase()
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	isDatabaseConnected()
	if isEmpty, err := isCollectionEmpty(); err != nil {
		log.Fatal(err)
	} else if isEmpty {
		fmt.Println("Found no documents in collection, seeding the database...")
		err = seedDatabase()
		if err != nil {
			log.Fatal(err)
		}
	}
}

// CreatePersonRecord creates a new person record in the database.
// If the application is in memory mode, it adds the person to the in-memory map.
// If the application is using MongoDB, it inserts the person into the database.
func CreatePersonRecord(person Person) (*mongo.InsertOneResult, error) {
	if isInMemory {
		objectId := primitive.NewObjectID()
		personMap[objectId.String()] = person
		return &mongo.InsertOneResult{InsertedID: objectId}, nil
	}

	result, err := peopleCollection.InsertOne(context.TODO(), person)
	if err != nil {
		return &mongo.InsertOneResult{InsertedID: primitive.NilObjectID}, err
	}

	return result, nil
}

// DeletePersonRecord deletes a person record from the database by its ObjectID.
// If the application is in memory mode, it deletes the person from the in-memory map.
// If the application is using MongoDB, it deletes the person record from the database.
func DeletePersonRecord(id string) (*mongo.DeleteResult, error) {
	if isInMemory {
		if _, ok := personMap[id]; ok {
			delete(personMap, id)
			return &mongo.DeleteResult{DeletedCount: 1}, nil
		} else {
			return nil, fmt.Errorf("no person with the given id was found")
		}
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectId}
	result, err := peopleCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetAllPeople retrieves all person records from the database that match the provided query.
// If the application is in memory mode, it retrieves all person records from the in-memory map.
// If the application is using MongoDB, it queries the database for matching records.
func GetAllPeople(query bson.M) ([]*Person, error) {

	if isInMemory {
		//	TODO: Handle query on personMap.
		return ConvertToSlice(personMap), nil
	}

	var result []*Person

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

	return result, nil
}

// GetPersonByObjectId retrieves a person record from the database by its ObjectID.
// If the application is in memory mode, it retrieves the person from the in-memory map.
// If the application is using MongoDB, it queries the database for the person record.
func GetPersonByObjectId(id string) (*Person, error) {
	if isInMemory {
		person, ok := personMap[id]
		if ok {
			return person.Clone(), nil
		} else {
			return nil, fmt.Errorf("person not found")
		}
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

// UpdatePersonRecord updates an existing person record in the database.
// If the application is in memory mode, it updates the person in the in-memory map.
// If the application is using MongoDB, it updates the person in the database.
func UpdatePersonRecord(person Person, id string) (*mongo.UpdateResult, error) {
	if isInMemory {
		_, ok := personMap[id]
		if !ok {
			return &mongo.UpdateResult{UpsertedID: primitive.NilObjectID, UpsertedCount: 0}, fmt.Errorf("person not found")
		}

		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return &mongo.UpdateResult{UpsertedID: primitive.NilObjectID, UpsertedCount: 0}, err
		}

		personMap[id] = person
		return &mongo.UpdateResult{UpsertedID: objectId, UpsertedCount: 1}, nil
	} else {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return &mongo.UpdateResult{UpsertedID: primitive.NilObjectID, UpsertedCount: 0}, err
		}

		filter := bson.M{"_id": objectId}
		update := bson.M{"$set": person}

		result, err := peopleCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return &mongo.UpdateResult{UpsertedID: primitive.NilObjectID, UpsertedCount: 0}, err
		}

		return result, nil
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

func generatePeople() People {
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

func seedDatabase() error {
	if isInMemory {
		for _, person := range generatePeople() {
			objectId := primitive.NewObjectID()
			personMap[objectId.String()] = person
		}

		return nil
	} else {
		_, err := peopleCollection.InsertMany(context.TODO(), generatePeople().ConvertToInterface())
		if err != nil {
			return err
		}

		return nil
	}
}
