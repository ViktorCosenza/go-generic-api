package routes

import (
	"fama-api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Start returns a router containing all the app routes
func Start(db *gorm.DB) (*gin.Engine, error) {
	db.LogMode(true)
	authMiddleware, err := middleware.GetJwtMiddleware(db)
	if err != nil {
		return nil, err
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	auth := r.Group("/auth")
	{
		auth.POST("/login", authMiddleware.LoginHandler)
	}

	admin := r.Group("/admin")
	admin.Use(authMiddleware.MiddlewareFunc())
	{
		auth.POST("/signup", createUser(db))           // Create new User
		admin.GET("/users", getUsers(db))              // Get all users, just for dev purpose TODO: DELETE THIS WHEN DONE//
		admin.POST("/text", createTexts(db))           // Add Texts via Zip file upload
		admin.POST("/ontology", createOntology(db))    // Insert json ontology
		admin.GET("/assigment", getAssigments(db))     // Get assigments count for each user
		admin.POST("/assigment", createAssigments(db)) // Assign texts to given users
	}

	annotation := r.Group("/annotation")
	{
		annotation.Use(authMiddleware.MiddlewareFunc())
		annotation.GET("/ontology", getOntology(db)) // Get Json ontology
		annotation.POST("/", createAnnotation(db))   // Submit annotaiton for a given text
		annotation.GET("/assigment", getAssigment(db))
	}

	return r, nil
}
