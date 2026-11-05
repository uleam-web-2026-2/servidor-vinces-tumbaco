package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/uleam-web-2026-2/KEMA/internal/config"
	"github.com/uleam-web-2026-2/KEMA/internal/intercambios"
	"github.com/uleam-web-2026-2/KEMA/internal/repository"
	"github.com/uleam-web-2026-2/KEMA/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// validarBaseDatos comprueba que la aplicación utilice
// exclusivamente la base de datos de recuperación.
func validarBaseDatos(databaseURL string) error {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return fmt.Errorf("la dirección de PostgreSQL no es válida")
	}

	if u.Hostname() != "localhost" && u.Hostname() != "127.0.0.1" {
		return fmt.Errorf("la conexión debe utilizar localhost o 127.0.0.1")
	}

	if u.Port() != "5434" {
		return fmt.Errorf("por seguridad, se requiere el puerto 5434")
	}

	if u.Path != "/intercambio_prueba" {
		return fmt.Errorf("por seguridad, se requiere la base de datos intercambio_prueba")
	}

	return nil
}

func main() {
	// Cargar configuración.
	cfg, err := config.Cargar()
	if err != nil {
		log.Fatal("Error de configuración: ", err)
	}

	// Verificar la base de datos antes de conectarse.
	if err := validarBaseDatos(cfg.DatabaseURL); err != nil {
		log.Fatal("Conexión bloqueada: ", err)
	}

	// Conectar con PostgreSQL.
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar con PostgreSQL: ", err)
	}

	fmt.Println("Conexión con PostgreSQL correcta")

	// Crear o actualizar las tablas de las 12 entidades.
	err = db.AutoMigrate(
		&intercambios.Usuario{},
		&intercambios.Libro{},
		&intercambios.Autor{},
		&intercambios.Categoria{},
		&intercambios.LibroAutor{},
		&intercambios.LibroCategoria{},
		&intercambios.EjemplarLibro{},
		&intercambios.SolicitudIntercambio{},
		&intercambios.Intercambio{},
		&intercambios.DetalleIntercambio{},
		&intercambios.Entrega{},
		&intercambios.Resena{},
	)
	if err != nil {
		log.Fatal("Error al preparar las tablas: ", err)
	}

	fmt.Println("Tablas preparadas correctamente")

	// Insertar datos iniciales.
	if err := intercambios.Sembrar(db); err != nil {
		log.Fatal("Error al insertar datos iniciales: ", err)
	}

	// CAPA 1: Repository.
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error al obtener sql.DB: %v", err)
	}
	repositorio := repository.NuevoIntercambiosRepository(db, sqlDB)

	// CAPA 2: Service.
	servicio := services.NuevoIntercambiosService(repositorio)

	// CAPA 3: Handler.
	handler := &intercambios.Handler{
		Servicio: servicio,
	}

	// Ruta principal.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"mensaje":"Servidor de intercambio de libros funcionando"}`)
	})

	// Ruta de salud.
	http.HandleFunc("/salud", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"estado":"ok"}`)
	})

	// Inicio de sesión.
	http.HandleFunc("/login", intercambios.Login)

	// Rutas de intercambios.
	rutasIntercambios := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			intercambios.Proteger(
				handler.Rutas,
				"administrador",
				"usuario",
			)(w, r)

		case http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete:

			intercambios.Proteger(
				handler.Rutas,
				"administrador",
			)(w, r)

		default:
			handler.Rutas(w, r)
		}
	}

	http.HandleFunc("/intercambios", rutasIntercambios)
	http.HandleFunc("/intercambios/", rutasIntercambios)

	// Configurar servidor HTTP.
	servidor := &http.Server{
		Addr:         ":" + cfg.Puerto,
		ReadTimeout:  cfg.TiempoEspera,
		WriteTimeout: cfg.TiempoEspera,
		Handler:      http.DefaultServeMux,
	}

	fmt.Println("Servidor iniciado en el puerto:", cfg.Puerto)

	// Iniciar servidor.
	if err := servidor.ListenAndServe(); err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
	}
}
