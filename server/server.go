package server

import (
	"fama-api/database"
	"fama-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//Server struct
type Server struct {
	db     *gorm.DB
	router *gin.Engine
}

//Listen to port
func (s *Server) Listen(port string) {
	s.router.Run(port)
}

//Start returns a server with default options
func Start() (*Server, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	r, err := routes.Start(db)
	if err != nil {
		return nil, err
	}

	return &Server{
		db:     db,
		router: r,
	}, nil
}
