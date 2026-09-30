package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/uleam-web-2026-2/KEMA/internal/config"
	"github.com/uleam-web-2026-2/KEMA/internal/libros"
)

func main() {
	reset := flag.Bool("reset", false, "reiniciar las tablas")
	flag.Parse()

	cfg, err := config.Cargar()
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Error conectando a PostgreSQL:", err)
	}

	if *reset {
		if err := db.Migrator().DropTable(&libros.Prestamo{}, &libros.Libro{}); err != nil {
			log.Fatal("Error eliminando tablas:", err)
		}
	}

	if err := db.AutoMigrate(&libros.Libro{}, &libros.Prestamo{}); err != nil {
		log.Fatal("Error en AutoMigrate:", err)
	}

	if *reset {
		if err := libros.Sembrar(db); err != nil {
			log.Fatal("Error sembrando datos:", err)
		}
	}

	manejador := libros.Manejador{DB: db}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(cfg.TiempoEspera))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.Get("/libros", manejador.ListarLibros)
	r.Post("/libros", manejador.CrearLibro)
	r.Get("/libros/{id}", manejador.ObtenerLibro)
	r.Put("/libros/{id}", manejador.ActualizarLibro)
	r.Delete("/libros/{id}", manejador.EliminarLibro)

	puerto := ":" + cfg.Puerto

	servidor := &http.Server{
		Addr:         puerto,
		Handler:      r,
		ReadTimeout:  cfg.TiempoEspera,
		WriteTimeout: cfg.TiempoEspera,
	}

	log.Println("Servidor escuchando en", puerto)

	if err := servidor.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
