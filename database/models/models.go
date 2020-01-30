package models

import (
	"github.com/jinzhu/gorm"
)

// User model
type User struct {
	gorm.Model
	Username     string `gorm:"unique, not null"`
	Password     string `gorm:"not null"`
	Assigments   []Assigment
}

//Admin model
type Admin struct {
	gorm.Model
	UserID uint `gorm:"unique, not null"`
	User User
}

// Text model
type Text struct {
	gorm.Model
	Name        string `gorm:"unique, not null"`
	Body        string `gorm:"not null"`
	AdminID 	uint `gorm:"not null"`
	Admin       Admin
	Annotations []Annotation
}

// JSONOntology JSON tree representing the ontology (Gambiarra)
type JSONOntology struct {
	Value string `gorm:"type:json"`
}

/*
// Class ontology Model
type Class struct {
	Class      string `gorm:"unique, not null"`
	ParentOf []SubClassOf
	SonOf []SubClassOf
}


// SubClassOf ontology Relation
type SubClassOf struct {
	SonID uint
	Son Class
	ParentID uint
	Parent Class
} 
*/

// Label a single label for a text annotation
type Label struct {
	gorm.Model
	first        string `gorm:"not null"`
	second       string
	third		 string
	fourth	     string
	start 	     uint   `gorm:"not null"`
	end          uint   `gorm:"not null"`
	AnnotationID uint   `gorm:"not null"`
}

// Assigment represents a text that can be annotated by a user
type Assigment struct {
	gorm.Model
	TextID  uint `gorm:"not null;"`
	UserID  uint `gorm:"not null;"`
}

// Annotation the whole annotation
type Annotation struct {
	gorm.Model
	Labels []Label
	AssigmentID uint `gorm:"not null;"`
}

