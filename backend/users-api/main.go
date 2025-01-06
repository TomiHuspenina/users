package main

import (
	"net/http"
	controller "users-api/controller"
	"users-api/middleware"
	repository "users-api/repositories"
	service "users-api/service"

	"github.com/gin-gonic/gin"
)

/*
type Controller interface {
	GetUserByID(ctx *gin.Context)
	Login(ctx *gin.Context)
	InsertUser(ctx *gin.Context)
	GetUserByName(ctx *gin.Context)
}*/

func main() {

	/*
		sqlconfig := repo.SQLConfig{
			Name: "users",
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASSWORD"),
			Host: os.Getenv("DB_HOST"),
		}*/
	/*
		sqlconfig := repo.SQLConfig{
			Name: "users",             // El nombre de la base de datos
			User: "root",              // El usuario de la base de datos
			Pass: "Tomihuspenina2003", // La contraseña de la base de datos
			Host: "localhost",         // El host donde se encuentra la base de datos
		}*/

	/*
		cacheConfig := repository.CacheConfig{
			MaxSize:      100000,
			ItemsToPrune: 100,
		}*/

	//username y password de tomi: root / root
	mongoConfig := repository.MongoConfig{
		Host:       "mongo",
		Port:       "27017",
		Username:   "root", //fran:
		Password:   "root", //fran:
		Database:   "hotels",
		Collection: "hotels",
	}

	//mainRepo := repo.NewSql(sqlconfig)
	//cacheRepository := repository.NewCache(cacheConfig)

	mainRepo := repository.NewMongo(mongoConfig)
	Service := service.NewService(mainRepo)
	Controller := controller.NewController(Service)
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type, X-Auth-Token")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	router.POST("/users", Controller.InsertUser)
	router.POST("/users/login", Controller.Login)
	router.GET("/users/token", Controller.Extrac)
	router.GET("/users/:id", Controller.GetUserById)

	authorized := router.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("/users", Controller.GetUserByName)

	}
	router.Run(":8080")
}
