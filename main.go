package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Franso/go-graphql-api/gql"
	"github.com/Franso/go-graphql-api/postgres"
	"github.com/Franso/go-graphql-api/server"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/graphql-go/graphql"
)

func main() {
	// Initialize our api and return a pointer to our router for http.ListenAndServer
	// and a pointer to our db to defer its closing when main() is finished
	router, db := initializeAPI()
	defer db.Close()

	// listen on port 4000 and if there is an error, log it and exit
	log.Fatal(http.ListenAndServe(":4000", router))
}

func initializeAPI() (*chi.Mux, *postgres.Db) {
	// Create a new  router
	router := chi.NewRouter()

	// Create a new connection to our postgres database
	db, err := postgres.New(
		postgres.ConnString("localhost", 5432, "muasya", "go_graphql_db"),
	)
	if err != nil {
		fmt.Println("Error creating schema: ", err)
	}

	// create our root query for qraphql
	rootQuery := gql.NewRoot(db)

	// Create a new graphql schema, passing in the root query
	sc, err := graphql.NewSchema(
		graphql.SchemaConfig{Query: rootQuery.Query},
	)
	if err != nil {
		fmt.Println("Error creating schema: ", err)
	}

	// create a server struct that holds a pointer to our database
	// as the addews to our graphql schema
	s := server.Server{
		GqlSchema: &sc,
	}

	// add some middleware to our router
	router.Use(
		// Set the content-type as application/json
		render.SetContentType(render.ContentTypeJSON),
		// Log api requests
		middleware.Logger,
		// Compress results, mostly gzipping assets and json
		middleware.Compress(5),
		// match routes with a trailing slash, strip it, and continue routing through the mux
		middleware.StripSlashes,
		// recover from panics without crashing the server
		middleware.Recoverer,
	)

	// Create the graphql route with a Server method to handle it
	router.Post("/graphql", s.GraphQL())
	return router, db

}
