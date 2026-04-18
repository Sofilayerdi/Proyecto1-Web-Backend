package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"

	_ "github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func router(conn net.Conn, db *sql.DB) {

r:
	chi.NewRouter()

	r.Get("/series", listarSeries)
	r.Get("/series/{id}", obtenerSerie)
	r.Post("/series", crearSerie)
	r.Put("/series/{id}", editarSerie)
	r.Delete("/series/{id}", eliminarSerie)

	http.ListenAndServe(conn, r)
}

func main() {
	listener, _ := net.Listen("tcp", ":8080")
	defer listener.Close()

	db, err := sql.Open("pgx", "postgres://postgres:TU_PASSWORD@localhost:5432/series_tracker?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Print("listening to port 8000")

	for {
		conn, _ := listener.Accept()
		go router(conn, db)
	}
}
