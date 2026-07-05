package http

import (
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

func NewServer(addr string, log *zap.Logger, svc *service.LinkService) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(RequestIDAndLoggerMiddleware(log))
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "server is running successfully!"})
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
