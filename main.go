package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/akhiltn/rssagg-go/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type ApiConfig struct{
  DB *database.Queries
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		portString = "8080"
		log.Println("Using Default Port: ", portString)
	} else {
		log.Println("Using Port: ", portString)
	}


	dbURL := os.Getenv("DB_URL")
	if dbURL  == "" {
		log.Fatal("DB_URL is not found in the environment")
	} else {
		log.Println("Using DB_URL: ", dbURL)
	}

  conn, err := sql.Open("postgres", dbURL)
  if err!=nil {
    log.Fatal("Can't connect to database", err)
  }

  apiCfg := ApiConfig{
    DB: database.New(conn),
  }

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
  v1Router.Post("/users", apiCfg.handlerCreateUser)
  v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
  v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
