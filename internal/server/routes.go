package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"hello_world/internal/handler"
	"hello_world/internal/middleware"
	"hello_world/internal/repository"
	"hello_world/internal/service"
)

func (s *Server) registerRoutes() http.Handler {
	if s.cfg.IsDevelopment() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// ── Global middleware ──────────────────────────────────────────────────────
	r.Use(
		middleware.RequestID(),
		middleware.RecoveryWithJSON(),
		middleware.Logger(),
		// 100 sustained req/s per IP, burst up to 200.
		middleware.RateLimit(rate.Limit(100), 200),
		middleware.CORS(s.cfg.CORSAllowedOrigins),
	)

	// ── Health probes (no auth required) ─────────────────────────────────────
	r.GET("/healthz/live", s.livenessHandler)
	r.GET("/healthz/ready", s.readinessHandler)

	// Legacy root – kept for backwards compatibility.
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
	})

	// ── Wire dependencies ─────────────────────────────────────────────────────
	db := s.db.GetDB()
	userRepo := repository.NewUserRepo(db)
	userSvc := service.NewUserService(userRepo, s.cfg.JWTSecret, s.cfg.JWTExpiryHours)
	userHandler := handler.NewUserHandler(userSvc)
	authHandler := handler.NewAuthHandler(userSvc)

	// ── Public auth routes ────────────────────────────────────────────────────
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// ── Protected API routes ──────────────────────────────────────────────────
	api := r.Group("/api/v1")
	api.Use(middleware.AuthRequired(s.cfg.JWTSecret))
	{
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)
		}
	}

	return r
}

func (s *Server) livenessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) readinessHandler(c *gin.Context) {
	health := s.db.Health(c.Request.Context())
	if health["status"] != "up" {
		c.JSON(http.StatusServiceUnavailable, health)
		return
	}
	c.JSON(http.StatusOK, health)
}
