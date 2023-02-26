package models

import (
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Title   string             `json:"title,omitempty"`
	Content string             `json:"content,omitempty"`
	Author  string             `json:"author,omitempty"`
}

func (a *Article) Validate() error {
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
