package http

import (
	_ "embed"
	"net/http"
	"trainer/internal/interfaces/http/handler"

	"github.com/gorilla/mux"
)

//go:embed docs/swagger.json
var swaggerJSON []byte

//go:embed docs/swagger.html
var swaggerHTML []byte

func NewRouter(
	authMiddleware mux.MiddlewareFunc,
	adminMiddleware mux.MiddlewareFunc,
	mentorMiddleware mux.MiddlewareFunc,
	userHandler *handler.UserHandler,
	loginHandler *handler.AuthHandler,
) http.Handler {
	r := mux.NewRouter()

	r.Use(authMiddleware)
	//r.Use(corsMiddleware)

	r.HandleFunc("swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(swaggerJSON)
	}).Methods("GET")

	r.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(swaggerHTML))
	}).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger", http.StatusMovedPermanently)
	}).Methods("GET")

	r.HandleFunc("/auth/access_token", loginHandler.AccessToken).Methods("POST")
	r.HandleFunc("/auth/refresh_token", loginHandler.RefreshToken).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(authMiddleware)

	adminRoutes := api.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(adminMiddleware)
	adminRoutes.HandleFunc("/users", userHandler.ListUser).Methods("GET")
	adminRoutes.HandleFunc("/users", userHandler.CreateUser).Methods("PUT")
	adminRoutes.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("POST")
	adminRoutes.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	adminRoutes.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	return r
}
