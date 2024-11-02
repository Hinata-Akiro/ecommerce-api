package main

import (
	"context"
	"ecommerce-api/config"
	"ecommerce-api/database"
	"ecommerce-api/models"
	"ecommerce-api/routes"
	"ecommerce-api/docs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router *gin.Engine
}

func NewServer() (*Server, error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("orderStatus", func(fl validator.FieldLevel) bool {
			status, ok := fl.Field().Interface().(models.OrderStatus)
			if !ok {
				return false
			}
			return status.IsValid() == nil
		})

		v.RegisterValidation("password", models.PasswordValidation)
	}

	err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		return nil, err
	}

	migrations()

	server := &Server{}
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.New()

	appConfig := config.Config()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	docs.SwaggerInfo.BasePath = "/api/v1"
    docs.SwaggerInfo.Host = appConfig.SWAGGER_SERVER_URL

	// Swagger setup
	url := ginSwagger.URL( "http://" + appConfig.SWAGGER_SERVER_URL + "/swagger/doc.json")
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// API routes
	apiGroup := router.Group("/api/v1/")
	apiGroup.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Connected!"})
	})

	routes.AuthSetUpRoute(apiGroup, appConfig.JWT_SECRET, database.Database)

	routes.ProductSetUpRoute(apiGroup, database.Database)

	routes.OrderSetUpRoute(apiGroup, database.Database)

	server.router = router
}

func (server *Server) Start(address string) error {
	srv := &http.Server{
		Addr:    address,
		Handler: server.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
	return nil
}

func migrations() {
	db := database.Database
	err := db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderProduct{})
	if err != nil {
		panic("failed to auto migrate database: " + err.Error())
	}
}

// @ECOMMERCE-API
// @version 1.0
// @description Your API description.
// @host localhost:4000
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	appConfig := config.Config()

	srv, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	if err := srv.Start(appConfig.PORT); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
