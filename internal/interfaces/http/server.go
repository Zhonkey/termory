package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trainer/internal/app"
	"trainer/internal/interfaces/http/middleware"

	"trainer/internal/infrastructure/database"
	"trainer/internal/interfaces/http/handler"
)

type Server struct {
	db         *database.DB
	httpServer *http.Server
}

func NewServer(db *database.DB) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) Run(ctx context.Context) error {
	c, err := app.NewContainer(s.db)
	if err != nil {
		return err
	}

	tokenHandler := handler.NewAuthTokenHandler(c.AccessTokenUC, c.RefreshTokenUC)
	userHandler := handler.NewUserHandler(c.CreateUserUC, c.UpdateUserUC, c.DeleteUserUC, c.GetUserUC, c.ListUserUC)

	authMiddleware := middleware.AuthMiddleware(c.TokenService)
	adminMiddleware := middleware.RoleMiddleware("admin")
	mentorMiddleware := middleware.RoleMiddleware("admin", "mentor")

	router := NewRouter(authMiddleware, adminMiddleware, mentorMiddleware, userHandler, tokenHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Starting server on port %s", port)
		serverErrors <- s.httpServer.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Printf("Received signal %v, starting graceful shutdown", sig)

		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.httpServer.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}

		log.Println("Server stopped gracefully")
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}
