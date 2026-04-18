package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

type Serie struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	CurrentEpisode int    `json:"current_ep"`
	TotalEpisodes  int    `json:"total_ep"`
	Img            string `json:"img"`
}

func listarSeries(w http.ResponseWriter, r *http.Request) {

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 2
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = 1
	}

	offset := limit * (page - 1)

	//obtener las series de la base de datos
	rows, err := db.Query("SELECT id, name, current_ep, total_ep, img FROM series LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		http.Error(w, "Error al obtener las series", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// guardar las series en un slice
	var series []Serie
	for rows.Next() {
		var s Serie
		rows.Scan(&s.ID, &s.Name, &s.CurrentEpisode, &s.TotalEpisodes, &s.Img)
		series = append(series, s)
	}

	//respuesta JSON con CORS
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	json.NewEncoder(w).Encode(series)
}

func verSerie(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	row := db.QueryRow("SELECT id, name, current_ep, total_ep, img FROM series WHERE id = $1", id)

	var s Serie
	err := row.Scan(&s.ID, &s.Name, &s.CurrentEpisode, &s.TotalEpisodes, &s.Img)
	if err != nil {
		http.Error(w, "Serie no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	json.NewEncoder(w).Encode(s)
}

func crearSerie(w http.ResponseWriter, r *http.Request) {
	var s Serie

	//obteniendo informacion
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "Contenido inválido", http.StatusBadRequest)
		return
	}

	if s.Name == "" || s.CurrentEpisode < 0 || s.TotalEpisodes < 0 || s.Img == "" {
		http.Error(w, "Datos incompletos", http.StatusBadRequest)
		return
	}

	if s.CurrentEpisode > s.TotalEpisodes {
		http.Error(w, "Episodio actual debe ser menor o igual al total", http.StatusBadRequest)
		return
	}

	//insertar a la base de datos
	db.Exec(
		"INSERT INTO series (name, current_ep, total_ep, img) VALUES ($1, $2, $3, $4)",
		s.Name, s.CurrentEpisode, s.TotalEpisodes, s.Img,
	)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.WriteHeader(http.StatusCreated)
}

func editarSerie(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var s Serie

	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "Contenido inválido", http.StatusBadRequest)
		return
	}

	if s.Name == "" || s.CurrentEpisode < 0 || s.TotalEpisodes < 0 || s.Img == "" {
		http.Error(w, "Datos incompletos", http.StatusBadRequest)
		return
	}

	if s.CurrentEpisode > s.TotalEpisodes {
		http.Error(w, "Episodio actual debe ser menor o igual al total", http.StatusBadRequest)
		return
	}

	row, err := db.Exec(
		"UPDATE series SET name = $1, current_ep = $2, total_ep = $3, img = $4 WHERE id = $5",
		s.Name, s.CurrentEpisode, s.TotalEpisodes, s.Img, id,
	)

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.WriteHeader(http.StatusOK)
}

func eliminarSerie(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	row, err := db.Exec("DELETE FROM series WHERE id = $1", id)
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	//Inicializar database
	var err error
	db, err = sql.Open("pgx", "postgres://postgres:Duquesa0@localhost:5432/series?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := chi.NewRouter()

	r.Get("/series", listarSeries)
	r.Get("/series/{id}", verSerie)
	r.Post("/series", crearSerie)
	r.Put("/series/{id}", editarSerie)
	r.Delete("/series/{id}", eliminarSerie)

	log.Print("listening to port 8000")

	http.ListenAndServe(":8000", r)
}
