package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// package database sert à la gestion de la base de données
// github.com/lib/pq sert à la connexion avec PostgreSQL

type Database struct {
	DB *sql.DB
}
// ce typage sert à indiquer que la variable DB est une pointeur de sql.DB, ce qui est nécessaire pour la fonction NewDatabase

func NewDatabase(connectionString string) (*Database, error) {
	db, err := sql.Open("postgres", connectionString) // sql.Open sert à ouvrir une connexion avec la base de données
	if err != nil { // si la connexion echoue, nil sert à indiquer que la variable est vide
		return nil, err
	}
	// db sert à stocker la connexion avec la base de données

	if err = db.Ping(); err != nil { // db.Ping sert à tester la connexion avec la base de données
		return nil, err
	}

	if err = createTables(db); err != nil { // createTables sert à créer les tables
		return nil, err
	}

	return &Database{DB: db}, nil //  &Database{DB: db} sert à indiquer que la variable DB est un pointeur de sql.DB
	// on le retourne afin de pouvoir l'utiliser dans le main
}

func createTables(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS items (
			id VARCHAR(50) PRIMARY KEY,
			name VARCHAR(100) NOT NULL
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Erreur lors de la création des tables: %v", err)
		return err
	}
	return nil
}
