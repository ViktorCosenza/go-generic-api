package server

import(
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"fama-api/routes"
	"fama-api/database"
)

//Server struct
type Server struct {
	db *gorm.DB
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
	r := routes.Start(db)
	return &Server{
		db: db,
		router: r,
	}, nil
}

