package main

import (
	"context"
	"elk-example/ent"
	elk "elk-example/ent/http"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func main() {
	// Create the ent client.
	c, err := ent.Open("sqlite3", "./ent.db?_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer c.Close()
	// Run the auto migration tool.
	if err := c.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	// Router, Logger and Validator.
	r, l, v := chi.NewRouter(), zap.NewExample(), validator.New()
	// Create the pet handler.
	r.Route("/pets", func(r chi.Router) {
		elk.NewPetHandler(c, l, v).Mount(r, elk.PetRoutes)
	})
	// Create the user handler.
	r.Route("/users", func(r chi.Router) {
		elk.NewUserHandler(c, l, v).Mount(r, elk.UserRoutes)
	})
	// Start listen to incoming requests.
	fmt.Println("Server running")
	defer fmt.Println("Server stopped")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
