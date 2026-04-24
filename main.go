package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

var db *sql.DB

type Serie struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	CurrentEpisode int    `json:"current_ep"`
	TotalEpisodes  int    `json:"total_ep"`
	Img            string `json:"img"`
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func listarSeries(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 5
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	offset := limit * (page - 1)

	rows, err := db.Query("SELECT id, name, current_ep, total_ep, image_url FROM series LIMIT ? OFFSET ?;", limit, offset)
	if err != nil {
		http.Error(w, "Error al obtener las series", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var series []Serie
	for rows.Next() {
		var s Serie
		err := rows.Scan(&s.ID, &s.Name, &s.CurrentEpisode, &s.TotalEpisodes, &s.Img)
		if err != nil {
			http.Error(w, "Error leyendo datos", http.StatusInternalServerError)
			return
		}
		series = append(series, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(series)
}

func verSerie(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	id := chi.URLParam(r, "id")

	var s Serie
	row := db.QueryRow("SELECT id, name, current_ep, total_ep, image_url FROM series WHERE id = ?;", id)
	err := row.Scan(&s.ID, &s.Name, &s.CurrentEpisode, &s.TotalEpisodes, &s.Img)
	if err == sql.ErrNoRows {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error al obtener la serie", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func crearSerie(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	var s Serie

	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "Contenido inválido", http.StatusBadRequest)
		return
	}

	// Validación
	if s.Name == "" || s.CurrentEpisode <= 0 || s.TotalEpisodes <= 0 || s.Img == "" {
		http.Error(w, "Datos incompletos", http.StatusBadRequest)
		return
	}

	if s.CurrentEpisode > s.TotalEpisodes {
		http.Error(w, "Episodio actual debe ser menor o igual al total", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(
		"INSERT INTO series (name, current_ep, total_ep, image_url) VALUES (?,?,?,?);",
		s.Name, s.CurrentEpisode, s.TotalEpisodes, s.Img,
	)
	if err != nil {
		http.Error(w, "Error al crear la serie", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	s.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func editarSerie(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	id := chi.URLParam(r, "id")

	var s Serie

	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "Contenido inválido", http.StatusBadRequest)
		return
	}

	if s.Name == "" || s.CurrentEpisode <= 0 || s.TotalEpisodes <= 0 || s.Img == "" {
		http.Error(w, "Datos incompletos", http.StatusBadRequest)
		return
	}

	if s.CurrentEpisode > s.TotalEpisodes {
		http.Error(w, "Episodio actual debe ser menor o igual al total", http.StatusBadRequest)
		return
	}

	row, err := db.Exec(
		"UPDATE series SET name = ?, current_ep = ?, total_ep = ?, image_url = ? WHERE id = ?;",
		s.Name, s.CurrentEpisode, s.TotalEpisodes, s.Img, id,
	)
	if err != nil {
		http.Error(w, "Error al actualizar la serie", http.StatusInternalServerError)
		return
	}

	affected, _ := row.RowsAffected()
	if affected == 0 {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func eliminarSerie(w http.ResponseWriter, r *http.Request) {
	enableCors(w)
	id := chi.URLParam(r, "id")

	row, err := db.Exec("DELETE FROM series WHERE id = ?;", id)
	if err != nil {
		http.Error(w, "Error al eliminar la serie", http.StatusInternalServerError)
		return
	}

	affected, _ := row.RowsAffected()
	if affected == 0 {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	var err error
	db, err = sql.Open("sqlite", "series.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := chi.NewRouter()

	r.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		enableCors(w)
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/series", listarSeries)
	r.Get("/series/{id}", verSerie)
	r.Post("/series", crearSerie)
	r.Put("/series/{id}", editarSerie)
	r.Delete("/series/{id}", eliminarSerie)

	log.Println("Servidor corriendo en http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
