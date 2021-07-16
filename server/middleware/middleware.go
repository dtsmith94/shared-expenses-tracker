package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dtsmith94/shared-expenses-tracker/server/models"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database Name
const dbName = "myFirstDatabase"

// Collection name
const collName = "sharedExpenses"

// collection object/instance
var collection *mongo.Collection

// create connection with mongo db
func init() {

	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal("Failed to open config file: ", err)
	}
	defer f.Close()

	var config models.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("Failed to decode config file: ", err)
	}

	clientOptions := options.Client().ApplyURI(config.Database.ConnectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	collection = client.Database(dbName).Collection(collName)
	fmt.Println("Collection instance created")

}

func setAccessControlHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Location")

}

// GetAllExpenses get all the expenses route
func GetAllExpenses(w http.ResponseWriter, r *http.Request) {
	setAccessControlHeaders(w)
	expenses, err := getAllExpenses()

	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(expenses)
}

// CreateExpense create expense route
func CreateExpense(w http.ResponseWriter, r *http.Request) {
	setAccessControlHeaders(w)
	var expense models.Expense
	_ = json.NewDecoder(r.Body).Decode(&expense)

	if len(expense.Name) == 0 || expense.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := insertExpense(expense)

	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if id != "" {
		w.Header().Set("Location", r.RequestURI+"/"+id)
		w.WriteHeader(http.StatusCreated)
		fmt.Println("Location", r.RequestURI+"/"+id)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// EditExpense mutate the expense record
func EditExpense(w http.ResponseWriter, r *http.Request) {

	setAccessControlHeaders(w)

	params := mux.Vars(r)

	var expense models.Expense
	err := json.NewDecoder(r.Body).Decode(&expense)

	if len(expense.Name) == 0 || expense.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	editExpense(params["id"], expense)

	json.NewEncoder(w).Encode(params["id"])
}

// DeleteTask delete one task route
func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	setAccessControlHeaders(w)
	params := mux.Vars(r)
	err := deleteExpense(params["id"])

	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(params["id"])
}

// DeleteAllExpenses delete all expenses route
func DeleteAllExpenses(w http.ResponseWriter, r *http.Request) {
	setAccessControlHeaders(w)
	count, err := deleteAllExpenses()

	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(count)

}

// get all expenses from the DB and return them
func getAllExpenses() ([]models.Expense, error) {
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	var results []models.Expense
	for cursor.Next(context.Background()) {
		var result models.Expense
		err = cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(context.Background())
	return results, err
}

// Insert one task in the DB
func insertExpense(expense models.Expense) (string, error) {
	expense.Created = time.Now()
	expense.Modified = time.Now()

	insertResult, err := collection.InsertOne(context.Background(), expense)

	if err != nil {
		return "", err
	}

	return insertResult.InsertedID.(primitive.ObjectID).Hex(), err
}

// update expense settled status
func editExpense(id string, expense models.Expense) error {
	fmt.Println(expense)
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"settled": expense.Settled, "name": expense.Name, "amount": expense.Amount, "modified": time.Now()}}
	result, err := collection.UpdateOne(context.Background(), filter, update)

	fmt.Println("modified count: ", result.ModifiedCount)
	return err
}

// delete one task from the DB, delete by ID
func deleteExpense(id string) error {
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}
	d, err := collection.DeleteOne(context.Background(), filter)

	fmt.Println("Deleted Document", d.DeletedCount)
	return err
}

// delete all the expenses from the DB
func deleteAllExpenses() (int64, error) {
	d, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)

	fmt.Println("Deleted Document", d.DeletedCount)
	return d.DeletedCount, err
}
