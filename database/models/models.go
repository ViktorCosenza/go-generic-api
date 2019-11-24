package models

import (
	"github.com/jinzhu/gorm"
)

// User model
type User struct {
	gorm.Model
	Name     string `gorm:"unique, not null"`
	Password string
}

// Text model
type Text struct {
	gorm.Model
	Title       string
	Body        string
	Annotations []Annotation
}

// Label a single label for a text annotation
type Label struct {
	gorm.Model
	concept      string
	target       string
	AnnotationID uint
}

// Annotation the whole annotation
type Annotation struct {
	gorm.Model
	Labels []Label
	TextID uint
}
