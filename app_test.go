package main

import (
	"Golang-CRUD-Api/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCreateArticle(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		inputJSON      string
		expectedStatus int
		expectedResp   string
	}{
		{
			name: "Valid Request",
			inputJSON: `{
                "title": "Test Article",
                "content": "Test Content",
                "author": "Test Author"
            }`,
			expectedStatus: http.StatusCreated,
			expectedResp: `{
				"data": {
					"id": "%s"
				},
				"message": "Success",
				"status": 201
			}`,
		},
		{
			name: "Empty Title",
			inputJSON: `{
		        "title": "",
		        "content": "Test Content",
		        "author": "Test Author"
		    }`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedResp: `{
		        "status": 422,
		        "message": "Validation failed!. title must not be empty",
		        "data": null
		    }`,
		},
		{
			name: "Empty Author",
			inputJSON: `{
		        "title": "Test Article",
		        "content": "Test Content",
		        "author": ""
		    }`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedResp: `{
		        "status": 422,
		        "message": "Validation failed!. author must not be empty",
		        "data": null
		    }`,
		},
		{
			name: "Empty Content",
			inputJSON: `{
		        "title": "Test Article",
		        "content": "",
		        "author": "Test Author"
		    }`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedResp: `{
		        "status": 422,
		        "message": "Validation failed!. content must not be empty",
		        "data": null
		    }`,
		},
		{
			name: "Invalid Body",
			inputJSON: `[{
		        "content": "",
		        "author": "Test Author"
		    }]`,
			expectedStatus: http.StatusBadRequest,
			expectedResp: `{
		        "status": 400,
		        "message": "Invalid request payload",
		        "data": null
		    }`,
		},
	}

	// Set up router
	router := mux.NewRouter()
	router.HandleFunc("/articles", createArticle).Methods("POST")

	// Iterate over test cases
	for tNum, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("%d. %s\n", tNum+1, tc.name)
			req, err := http.NewRequest("POST", "/articles", bytes.NewBuffer([]byte(tc.inputJSON)))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder to record the response
			rr := httptest.NewRecorder()

			// Call the CreateArticle API handler function with the request and response recorder
			router.ServeHTTP(rr, req)

			// Check the status code returned by the API handler function
			assert.Equal(t, tc.expectedStatus, rr.Code)
			actualResponse := rr.Body.String()

			if tc.name == "Valid Request" {
				var result map[string]interface{}
				err = json.NewDecoder(rr.Body).Decode(&result)
				if err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}
				if _, ok := result["data"].(map[string]interface{})["id"]; !ok {
					t.Errorf("Expected 'id' field in response, but got %v", result)
				}
				tc.expectedResp = fmt.Sprintf(tc.expectedResp, result["data"].(map[string]interface{})["id"].(string))
				fmt.Printf("expectedResponse: %s\n", tc.expectedResp)
				fmt.Printf("actualResponse: %s\n", actualResponse)
				assert.JSONEq(t, tc.expectedResp, actualResponse)
			} else {
				fmt.Printf("expectedResponse: %s\n", tc.expectedResp)
				fmt.Printf("actualResponse: %s\n", actualResponse)
				assert.JSONEq(t, tc.expectedResp, actualResponse)
			}
		})
	}
}

func TestGetAllArticles(t *testing.T) {
	// Clear the articles collection
	err := clearCollection(client, "articles")
	if err != nil {
		t.Fatalf("Error clearing collection: %v", err)
	}

	// Insert some test articles into the collection
	articles := []models.Article{
		{
			Title:   "Article 1",
			Author:  "Author 1",
			Content: "This is the content of article 1.",
		},
		{
			Title:   "Article 2",
			Author:  "Author 2",
			Content: "This is the content of article 2.",
		},
		{
			Title:   "Article 3",
			Author:  "Author 3",
			Content: "This is the content of article 3.",
		},
	}
	for _, article := range articles {
		_, err := client.Database("mydb").Collection("articles").InsertOne(context.Background(), article)
		if err != nil {
			t.Fatalf("Error inserting article: %v", err)
		}
	}

	// Test case 1: Get all articles successfully
	fmt.Println("valid case of GET all articles")
	req1, err := http.NewRequest("GET", "/articles", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	rr1 := httptest.NewRecorder()
	handler1 := http.HandlerFunc(getAllArticles)
	handler1.ServeHTTP(rr1, req1)

	expectedRespLen := 3
	assert.Equal(t, http.StatusOK, rr1.Code)
	assert.Equal(t, expectedRespLen, len(articles))
}

func TestGetArticleById(t *testing.T) {
	// Define test cases
	// testCases := []struct {
	// 	name           string
	// 	inputParam     string
	// 	expectedStatus int
	// 	expectedResp   string
	// }{
	// 	{
	// 		name:           "Valid Request",
	// 		inputParam:     "%s",
	// 		expectedStatus: http.StatusOK,
	// 		expectedResp: `{
	// 			"data": {
	// 				"id": "%s"
	// 			},
	// 			"message": "Success",
	// 			"status": 200
	// 		}`,
	// 	},
	// 	{
	// 		name:           "Invalid Param",
	// 		inputParam:     "abcde",
	// 		expectedStatus: http.StatusBadRequest,
	// 		expectedResp: `{
	// 	        "status": 400,
	// 	        "message": "Invalid article ID",
	// 	        "data": null
	// 	    }`,
	// 	},
	// 	{
	// 		name:           "Not Found",
	// 		inputParam:     "63fbbb319358d3501cc3b356",
	// 		expectedStatus: http.StatusNotFound,
	// 		expectedResp: `{
	// 	        "status": 404,
	// 	        "message": "Article Not Found",
	// 	        "data": null
	// 	    }`,
	// 	},
	// }

	// Set up router
	// Test case 1: Get all articles successfully
	fmt.Println("1. valid case of GET an articles")
	router := mux.NewRouter()
	router.HandleFunc("/articles/{article_id}", getArticleById).Methods("GET")

	var article = models.Article{
		Title:   "Test Title",
		Content: "Test Content",
		Author:  "Test Author",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	dbRes, err := client.Database("mydb").Collection("articles").InsertOne(ctx, article)
	if err != nil {
		t.Errorf("Error creating article: %v", err)
	}
	article_id := dbRes.InsertedID.(primitive.ObjectID).Hex()

	// Test valid case
	rr1 := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/articles/"+article_id, nil)
	router.ServeHTTP(rr1, req)

	expected := `{
			"data": [{
			"id": "` + article_id + `",
			"title": "Test Title",
			"content": "Test Content",
			"author": "Test Author"
		}],
		"status": 200,
		"message": "Success"
	}`

	fmt.Printf("expectedResponse: %s\n", expected)
	fmt.Printf("actualResponse: %s\n", rr1.Body.String())
	assert.Equal(t, http.StatusOK, rr1.Code)
	assert.JSONEq(t, expected, rr1.Body.String())

	//2. invalid article id
	fmt.Println("2. Invalid Article Id case of GET article by ID")
	article_id = "abcde"
	req2, err := http.NewRequest("GET", "/articles/"+article_id, nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	expectedResp2 := `{
		"data": null,
		"status": 400,
		"message": "Invalid article ID"
	}`
	fmt.Printf("expectedResponse: %s\n", expectedResp2)
	fmt.Printf("actualResponse: %s\n", rr2.Body.String())
	assert.Equal(t, http.StatusBadRequest, rr2.Code)
	assert.JSONEq(t, expectedResp2, rr2.Body.String())

	//3. article not found
	fmt.Println("3. Article not found case of GET article by ID")
	article_id = "63fbbb319358d3501cc3b356"
	req3, err := http.NewRequest("GET", "/articles/"+article_id, nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	rr3 := httptest.NewRecorder()
	router.ServeHTTP(rr3, req3)

	expectedResp3 := `{
		"data": null,
		"status": 404,
		"message": "Article Not Found"
	}`
	fmt.Printf("expectedResponse: %s\n", expectedResp3)
	fmt.Printf("actualResponse: %s\n", rr3.Body.String())
	assert.Equal(t, http.StatusNotFound, rr3.Code)
	assert.JSONEq(t, expectedResp3, rr3.Body.String())
}

func clearCollection(client *mongo.Client, col string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := client.Database("mydb").Collection(col).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
