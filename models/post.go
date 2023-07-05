package models

import "github.com/google/uuid"

type Post struct {
	ID			uuid.UUID	`gorm:"type:uuid;column:id;primaryKey;default:gen_random_uuid()" json:"id"`

	Title	   	string 		`json:"title"`
	Content 	string 		`json:"content"`
	Slug    	string 		`json:"slug" gorm:"uniqueIndex"`
}