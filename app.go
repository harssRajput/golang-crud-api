package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"Golang-CRUD-Api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gorilla/mux"
)

func app() {
	// Set up router
	router := mux.NewRouter()

	router.HandleFunc("/articles", createArticle).Methods("POST")
	router.HandleFunc("/articles/{article_id}", getArticleById).Methods("GET")
	router.HandleFunc("/articles", getAllArticles).Methods("GET")

	fmt.Println("Server is ready!")
	log.Fatal(http.ListenAndServe(":"+SERVER_PORT, router))
}

// Get a single article by ID
func getArticleById(w http.ResponseWriter, r *http.Request) {

	fmt.Println("GET article by ID.")

	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["article_id"])
	if err != nil {
		fmt.Println("invalid article ID")
		sendResponse(400, "Invalid article ID", nil, w)
		return
	}

	collection := client.Database("mydb").Collection("articles")
	// Find the article by ID in mongo
	var article models.Article
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&article)
	if err != nil {
		fmt.Printf("Article not found for id: %s\n", vars["article_id"])
		sendResponse(404, "Article Not Found", nil, w)
		return
	}

	fmt.Println("got article ", article)
	sendResponse(200, "Success", []models.Article{article}, w)
}

// Get all articles
func getAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET all articles")

	collection := client.Database("mydb").Collection("articles")
	// get all articles in the db
	var articles []models.Article
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		fmt.Println("Failed to fetch articles")
		sendResponse(500, "Failed to fetch articles", nil, w)
		return
	}

	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var article models.Article
		if err := cursor.Decode(&article); err != nil {
			fmt.Println("Failed to Decode articles")
			sendResponse(500, "Failed to Decode articles", nil, w)
			return
		}
		articles = append(articles, article)
	}
	if err := cursor.Err(); err != nil {
		fmt.Println("Failed to Decode articles")
		sendResponse(500, "Failed to Decode articles", nil, w)
		return
	}

	sendResponse(200, "Success", articles, w)
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	var article models.Article
	fmt.Println("POST article request.")

	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		fmt.Println("invalid request body")
		sendResponse(400, "Invalid request payload", nil, w)
		return
	}

	//validation
	if err := article.Validate(); err != nil {
		fmt.Println("Validation failed!.", err)
		sendResponse(422, "Validation failed!. "+err.Error(), nil, w)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := client.Database("mydb").Collection("articles").InsertOne(ctx, article)
	if err != nil {
		fmt.Println("Error creating article", err)
		sendResponse(500, "Error creating article "+err.Error(), nil, w)
		return
	}

	fmt.Printf("Article with id %s inserted in MondoDB\n", result.InsertedID.(primitive.ObjectID).Hex())
	sendResponse(201, "Success", map[string]string{
		"id": result.InsertedID.(primitive.ObjectID).Hex(),
	}, w)
}

func sendResponse(status int, message string, data interface{}, w http.ResponseWriter) {
	response := map[string]interface{}{
		"status":  status,
		"message": message,
		"data":    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
