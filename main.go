package main

import (
	"database/sql"
	"log"
	"net/http"

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
