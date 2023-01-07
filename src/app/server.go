package app

import (
	"fmt"

	"github.com/Ayobami-00/booketo-mvc-go-postgres-gin/src/utils"
	token "github.com/Ayobami-00/booketo-mvc-go-postgres-gin/src/utils/auth"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our service.
type Server struct {
	config utils.Config
	router *gin.Engine
}

var (
	TokenMaker token.Maker
	Config     utils.Config
)

// NewServer creates a new HTTP server and set up routing.
func NewServer(config utils.Config) (*Server, error) {
	var err error

	TokenMaker, err = token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	Config = config

	server := &Server{
		config: config,
	}

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("currency", validCurrency)
	// }

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	mapUrls(router)

	// authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	// authRoutes.POST("/accounts", server.createAccount)
	// authRoutes.GET("/accounts/:id", server.getAccount)
	// authRoutes.GET("/accounts", server.listAccounts)

	// authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
