package router

import (
	"hello_world/internal/handler"
	"hello_world/internal/repository"
	"hello_world/internal/service"
	
				"net/http"
	"github.com/gin-gonic/gin"
		
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	gormDB := s.db.GetDB()
	api := r.Group("/api/v1")

	
		userRepo := repository.NewUserRepo(gormDB)
		userService := service.NewUserService(userRepo)
		userHandler := handler.NewUserHandler(userService)

		{
			user := api.Group("/user")
			{
				user.GET("", userHandler.GetUsers)
			}
		}
	

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context)  {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	 c.JSON(http.StatusOK,  resp)
}

func (s *Server) healthHandler(c *gin.Context)  {
	 c.JSON(http.StatusOK,  s.db.Health())
}
