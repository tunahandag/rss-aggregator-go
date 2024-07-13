package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/tunahandag/rss-aggregator/internal/database"
)

type apiConfig struct {
	db *database.Queries
}

func main() {
	
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	
	if portString == "" {
		log.Fatal("PORT must be set")
	}
	fmt.Println("PORT is set to: ", portString)

	dbUrl := os.Getenv("DB_URL")
	
	if dbUrl == "" {
		log.Fatal("DB_URL must be set")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("sql.Open: ", err)
	}
	
	apiCfg := apiConfig {
		db: database.New(conn),
	}


	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: true,
		MaxAge: 300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handleCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollows))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollow))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}
	log.Printf("Starting server on port %v", portString)
	
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}