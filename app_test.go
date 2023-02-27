package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
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
	}

	// Set up router
	router := mux.NewRouter()
	router.HandleFunc("/articles", CreateArticle).Methods("POST")

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
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

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
