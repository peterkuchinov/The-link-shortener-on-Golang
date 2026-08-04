package http

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/logger"
	"github.com/peterkuchinov/The-link-shortener-on-Golang/internal/service"
	"go.uber.org/zap"
)

type Server struct {
	srv *http.Server
}

func (s *Server) Shutdown(ctx context.Context) error {
    return s.srv.Shutdown(ctx)
}

func NewServer(addr string, baseURL string, log *zap.Logger, svc *service.LinkService) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(RequestIDAndLoggerMiddleware(log))
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "server is running successfully!"})
	})
	
	r.POST("/shorten", func(c *gin.Context) {
		var req struct {
			URL        string `json:"url" binding:"required"`
			CustomCode string `json:"custom_code"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}
		
		code, err := svc.Shorten(c.Request.Context(), req.URL, req.CustomCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"short_url": fmt.Sprintf("%s/%s", baseURL, code), "code": code,})
	})
	
	r.GET("/:code", func(c *gin.Context) {
		code := c.Param("code")
		if code == "" || code == "favicon.ico" {
			c.Status(http.StatusNotFound)
			return
		}
		
		originalURL, err := svc.GetOriginalURL(c.Request.Context(), code)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "link not found"})
			return
		}
		
		c.Redirect(http.StatusTemporaryRedirect, originalURL)
	})

	return &Server{
		srv: &http.Server{
			Addr:         addr,
			Handler:      r,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func RequestIDAndLoggerMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t1 := time.Now()
		
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateUUID()
		}
		
		c.Header("X-Request-ID", requestID)
		
		requestLogger := log.With(zap.String("request_id", requestID))
		
		ctx := logger.ToContext(c.Request.Context(), requestLogger)
		c.Request = c.Request.WithContext(ctx)
		
		c.Next()
		
		latency := time.Since(t1)
		requestLogger.Info("http request handled",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
		)
	}
}

func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}