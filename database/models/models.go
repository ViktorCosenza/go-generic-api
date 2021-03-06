package models

import (
	"github.com/jinzhu/gorm"
)

// User model
type User struct {
	gorm.Model
	Username   string `gorm:"unique, not null"`
	Password   string `gorm:"not null" json:"-"`
	Assigments []Assigment
}

//Admin model
type Admin struct {
	gorm.Model
	UserID uint `gorm:"unique, not null"`
	User   User
}

// JSONOntology JSON tree representing the ontology (Gambiarra)
type JSONOntology struct {
	Value string `gorm:"type:json"`
}

// Label a single label for a text annotation
type Label struct {
	gorm.Model
	First       string `gorm:"not null" json:"first" binding:"required"`
	Second      string `json:"second" binding:"required"`
	Third       string `json:"third" binding:"required"`
	Fourth      string `json:"fourth" binding:"required"`
	Explicit    bool   `json:"explicit" binding:"required"`
	Start       uint   `gorm:"not null" json:"start" binding:"required"`
	End         uint   `gorm:"not null" json:"end" binding:"required"`
	AssigmentID uint   `gorm:"not null"`
}

// Assigment represents a text that can be annotated by a user
type Assigment struct {
	gorm.Model
	TextID uint `gorm:"not null;"`
	UserID uint `gorm:"not null;"`
	Labels []Label
}

// Text model
type Text struct {
	gorm.Model
	Name       string `gorm:"unique, not null"`
	Body       string `gorm:"not null"`
	AdminID    uint   `gorm:"not null"`
	Admin      Admin
	Assigments []Assigment
}
