package main

import (
	"flag"
	"fmt"
	"golang-crud-api-minimal/database"
	models "golang-crud-api-minimal/internal"

	"log"
	"net/http"
	"os"
)

type config struct {
	port int
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	models   models.Models
}

const (
	DSN = "data.sqlite"
)

func main() {

	var cfg config
	cfg.port = 9090

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := database.ConnectSQLite(DSN)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	defer db.SQL.Close()

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		models:   models.New(db.SQL),
	}

	migrate := flag.Bool("migrate", false, "Create the tables in the database")
	flag.Parse()

	if *migrate {
		app.infoLog.Println("Creating tables ...")

		if err := createTables(); err != nil {
			log.Fatal(err)
		}
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) serve() error {
	app.infoLog.Println("API listening on port", app.config.port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}
	return srv.ListenAndServe()
}

func createTables() error {
	db, _ := database.ConnectSQLite(DSN)
	query := `CREATE TABLE IF NOT EXISTS categories (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name VARCHAR(64) NULL,
				description VARCHAR(200) NULL,
				created_at TIMESTAMP DEFAULT DATETIME,
				updated_at TIMESTAMP NOT NULL
			  );
			  CREATE TABLE IF NOT EXISTS products (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name VARCHAR(64) NULL,
				description VARCHAR(200) NULL,
				price INTEGER NOT NULL,
				category_id INTEGER NOT NULL,
				created_at TIMESTAMP DEFAULT DATETIME,
				updated_at TIMESTAMP NOT NULL
			  );
			  INSERT INTO categories (name, description, created_at, updated_at)
				VALUES ('Electrónica', 'Productos electrónicos', DATETIME('now'), DATETIME('now'));

			  INSERT INTO categories (name, description, created_at, updated_at)
				VALUES ('Ropa', 'Prendas de vestir', DATETIME('now'), DATETIME('now'));

			  INSERT INTO categories (name, description, created_at, updated_at)
				VALUES ('Hogar', 'Artículos para el hogar', DATETIME('now'), DATETIME('now'));

			  INSERT INTO categories (name, description, created_at, updated_at)
				VALUES ('Deportes', 'Equipos deportivos', DATETIME('now'), DATETIME('now'));

			  INSERT INTO categories (name, description, created_at, updated_at)
				VALUES ('Juguetes', 'Juguetes para niños', DATETIME('now'), DATETIME('now'));

			-- Inserts para la tabla 'products'
			  INSERT INTO products (name, description, price, category_id, created_at, updated_at)
				VALUES ('Teléfono móvil', 'Smartphone de última generación', 799, 1, DATETIME('now'), DATETIME('now'));

			  INSERT INTO products (name, description, price, category_id, created_at, updated_at)
				VALUES ('Camiseta', 'Camiseta de algodón', 20, 2, DATETIME('now'), DATETIME('now'));

			  INSERT INTO products (name, description, price, category_id, created_at, updated_at)
				VALUES ('Sartén antiadherente', 'Sartén para cocinar', 35, 3, DATETIME('now'), DATETIME('now'));

			  INSERT INTO products (name, description, price, category_id, created_at, updated_at)
				VALUES ('Balón de fútbol', 'Balón oficial de la FIFA', 50, 4, DATETIME('now'), DATETIME('now'));

			  INSERT INTO products (name, description, price, category_id, created_at, updated_at)
				VALUES ('Muñeca', 'Muñeca de peluche para niños', 15, 5, DATETIME('now'), DATETIME('now'));

		  `

	_, err := db.SQL.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
