package main

import (
	"example/hello/internal/config"
	"example/hello/internal/database"
	"example/hello/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Charger la configuration
	cfg := config.NewConfig()

	// Initialiser la base de données
	db, err := database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Erreur de connexion à la base de données: %v", err)
	}

	// Créer le handler // un handler sert à faire une action sur une ressource
	itemHandler := handlers.NewItemHandler(db)

	// Configurer le router // le router sert à recevoir les requêtes et à les traiter
	router := mux.NewRouter()

	// Définir les routes // les routes sont les endpoints (les URLs)
	router.HandleFunc("/items", itemHandler.GetItems).Methods("GET")
	router.HandleFunc("/items/{id}", itemHandler.GetItem).Methods("GET")
	router.HandleFunc("/items", itemHandler.CreateItem).Methods("POST")
	router.HandleFunc("/items/{id}", itemHandler.UpdateItem).Methods("PUT")
	router.HandleFunc("/items/{id}", itemHandler.DeleteItem).Methods("DELETE")

	// Démarrer le serveur // le serveur sert à recevoir les requêtes et à les traiter
	log.Printf("Serveur démarré sur %s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(cfg.ServerPort, router))
}
