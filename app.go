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

// Get a single article by ID
func getArticleById(w http.ResponseWriter, r *http.Request) {

	fmt.Println("GET article by ID.")
	response := map[string]interface{}{
		"status":  "",
		"message": "",
		"data":    nil,
	}

	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["article_id"])
	if err != nil {
		fmt.Println("invalid article ID")

		response["status"] = http.StatusBadRequest
		response["message"] = "invalid article ID"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	collection := client.Database("mydb").Collection("articles")
	// Find the article by ID in mongo
	var article models.Article
	err = collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&article)
	if err != nil {
		fmt.Printf("Article not found for id: %s\n", vars["article_id"])

		response["status"] = http.StatusNotFound
		response["message"] = "Article Not Found"
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response["status"] = 200
	response["message"] = "Success"
	response["data"] = []models.Article{article}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}

// Get all articles
func getAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET all articles")
	response := map[string]interface{}{
		"status":  "",
		"message": "",
		"data":    nil,
	}

	collection := client.Database("mydb").Collection("articles")
	// get all articles in the db
	var articles []models.Article
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		fmt.Println("Failed to fetch articles")

		response["status"] = http.StatusInternalServerError
		response["message"] = "Failed to fetch articles"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var article models.Article
		if err := cursor.Decode(&article); err != nil {
			fmt.Println("Failed to Decode articles")

			response["status"] = http.StatusInternalServerError
			response["message"] = "Failed to Decode articles"
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		articles = append(articles, article)
	}
	if err := cursor.Err(); err != nil {
		fmt.Println("Failed to Decode articles")

		response["status"] = http.StatusInternalServerError
		response["message"] = "Failed to Decode articles"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response["status"] = 200
	response["message"] = "Success"
	response["data"] = articles
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	var article models.Article
	// = models.Article{
	// 	Title:   "spider man",
	// 	Content: "spider man is exposed this morning",
	// 	Author:  "peter parker",
	// }
	fmt.Println("POST article request.")
	response := map[string]interface{}{
		"status":  "",
		"message": "",
		"data":    nil,
	}

	err := json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		fmt.Println("invalid request body")
		response["status"] = http.StatusBadRequest
		response["message"] = "Invalid request payload"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	//validation
	if err := article.Validate(); err != nil {
		fmt.Println("validation failed. ", err)
		response["status"] = 422
		response["message"] = "Validation failed!." + err.Error()
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(response)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := client.Database("mydb").Collection("articles").InsertOne(ctx, article)
	if err != nil {
		fmt.Println("Error creating article", err)
		response["status"] = http.StatusInternalServerError
		response["message"] = "Error creating article" + err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	fmt.Printf("Article with id %s inserted in MondoDB\n", result.InsertedID.(primitive.ObjectID).Hex())

	response["status"] = 201
	response["message"] = "Success"
	response["data"] = map[string]string{
		"id": result.InsertedID.(primitive.ObjectID).Hex(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func app() {
	// Set up router
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get home request")
		response := map[string]interface{}{
			"status":  201,
			"message": "it's a home page.",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)

	}).Methods("GET")

	router.HandleFunc("/articles", CreateArticle).Methods("POST")
	router.HandleFunc("/articles/{article_id}", getArticleById).Methods("GET")
	router.HandleFunc("/articles", getAllArticles).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
