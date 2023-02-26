package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Title   string             `json:"title,omitempty"`
	Content string             `json:"content,omitempty"`
	Author  string             `json:"author,omitempty"`
}

func (a *Article) validate() error {
	if len(strings.TrimSpace(a.Title)) == 0 {
		return errors.New("title is required")
	}

	if len(strings.TrimSpace(a.Content)) == 0 {
		return errors.New("content is required")
	}

	if len(strings.TrimSpace(a.Author)) == 0 {
		return errors.New("author is required")
	}

	return nil
}

func CreateArticle() {
	var article Article = Article{
		Title:   "spider man",
		Content: "spider man is exposed this morning",
		Author:  "peter parker",
	}

	//validation
	if err := article.validate(); err != nil {
		fmt.Println("validation failed. ", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := client.Database("mydb").Collection("articles").InsertOne(ctx, article)
	if err != nil {
		fmt.Println("unable to save in db", err)
		return
	}

	response := map[string]interface{}{
		"status":  201,
		"message": "Success",
		"data": map[string]string{
			"id": result.InsertedID.(primitive.ObjectID).Hex(),
		},
	}

	fmt.Println(response, "after success")
}

func main() {
	//first init() function executes present in init.go
	fmt.Println("came here after mongo connection")
	CreateArticle()
}
