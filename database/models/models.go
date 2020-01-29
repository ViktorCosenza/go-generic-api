package models

import (
	"github.com/jinzhu/gorm"
)

// User model
type User struct {
	gorm.Model
	Username     string `gorm:"unique, not null"`
	Password     string `gorm:"not null"`
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
	Body        string `gorm:"not null"`
	AdminID     uint `gorm:"not null"`
	Admin       Admin
	Annotations []Annotation
}

// JSONOntology JSON tree representing the ontology (Gambiarra)
type JSONOntology struct {
	Value string `gorm:"type:json"`
}

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

// Label a single label for a text annotation
type Label struct {
	gorm.Model
	first        string
	second       string
	third		 string
	fourth	     string
	start 	     uint
	end          uint
	AnnotationID uint
	Annotation   Annotation
}

// Annotation the whole annotation
type Annotation struct {
	gorm.Model
	Labels []Label
	TextID uint
	Text   Text
}
