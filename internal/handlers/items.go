package handlers

import (
	"database/sql"
	"encoding/json"
	"example/hello/internal/database"
	"example/hello/internal/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ItemHandler struct {
	db *database.Database
}

func NewItemHandler(db *database.Database) *ItemHandler {
	return &ItemHandler{db: db}
}

func (h *ItemHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println("Récupération de tous les items...")
	rows, err := h.db.DB.Query("SELECT id, name FROM items")
	if err != nil {
		log.Printf("Erreur de requête: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			log.Printf("Erreur de scan: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	log.Printf("Nombre d'items trouvés: %d\n", len(items))
	json.NewEncoder(w).Encode(items)
}

func (h *ItemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var item models.Item
	err := h.db.DB.QueryRow("SELECT id, name FROM items WHERE id = $1", params["id"]).Scan(&item.ID, &item.Name)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		log.Printf("Erreur de décodage JSON: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Création d'un nouvel item: ID=%s, Name=%s\n", item.ID, item.Name)
	_, err := h.db.DB.Exec("INSERT INTO items (id, name) VALUES ($1, $2)", item.ID, item.Name)
	if err != nil {
		log.Printf("Erreur d'insertion en base: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Item créé avec succès: %+v\n", item)
	json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.db.DB.Exec("UPDATE items SET name = $1 WHERE id = $2", item.Name, params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	item.ID = params["id"]
	json.NewEncoder(w).Encode(item)
}

func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := h.db.DB.Exec("DELETE FROM items WHERE id = $1", params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
