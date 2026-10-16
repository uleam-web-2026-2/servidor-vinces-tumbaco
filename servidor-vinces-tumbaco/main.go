package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/uleam-web-2026-2/KEMA/internal/intercambios"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// ==========================================
	// CONEXION A POSTGRESQL
	// ==========================================

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		log.Fatal("Debes configurar DATABASE_URL antes de iniciar el servidor")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error al conectar con PostgreSQL: ", err)
	}

	fmt.Println("Conexion con PostgreSQL correcta")

	// ==========================================
	// MIGRACIONES DE LAS 12 ENTIDADES
	// ==========================================

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
		log.Fatal("Error al crear las tablas: ", err)
	}

	fmt.Println("Tablas preparadas correctamente")

	// ==========================================
	// DATOS INICIALES
	// ==========================================

	if err := intercambios.Sembrar(db); err != nil {
		log.Fatal("Error al insertar datos iniciales: ", err)
	}

	fmt.Println("Datos iniciales preparados correctamente")

	// ==========================================
	// VALIDAR CONFIGURACION JWT
	// ==========================================

	if len(os.Getenv("JWT_SECRET")) < 32 {
		log.Fatal("JWT_SECRET debe tener al menos 32 caracteres")
	}

	if os.Getenv("DEMO_ADMIN_PASSWORD") == "" {
		log.Fatal("Debes configurar DEMO_ADMIN_PASSWORD")
	}

	if os.Getenv("DEMO_USER_PASSWORD") == "" {
		log.Fatal("Debes configurar DEMO_USER_PASSWORD")
	}

	fmt.Println("Configuracion JWT correcta")

	// ==========================================
	// HANDLER DE INTERCAMBIOS
	// ==========================================

	handler := &intercambios.Handler{
		DB: db,
	}

	// ==========================================
	// RUTA PRINCIPAL
	// ==========================================

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(
			w,
			`{"mensaje":"Servidor de intercambio de libros funcionando"}`,
		)
	})

	// ==========================================
	// RUTA DE SALUD
	// ==========================================

	http.HandleFunc("/salud", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(
			w,
			`{"estado":"ok"}`,
		)
	})

	// ==========================================
	// INICIO DE SESION JWT
	// ==========================================

	http.HandleFunc("/login", intercambios.Login)

	// ==========================================
	// CRUD DE INTERCAMBIOS CON PERMISOS
	// ==========================================

	rutasIntercambios := func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			// Administrador y usuario pueden consultar.
			intercambios.Proteger(
				handler.Rutas,
				"admin",
				"usuario",
			)(w, r)

		case http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete:

			// Solo el administrador puede modificar.
			intercambios.Proteger(
				handler.Rutas,
				"admin",
			)(w, r)

		default:
			handler.Rutas(w, r)
		}
	}

	http.HandleFunc(
		"/intercambios",
		rutasIntercambios,
	)

	http.HandleFunc(
		"/intercambios/",
		rutasIntercambios,
	)

	// ==========================================
	// INICIAR SERVIDOR
	// ==========================================

	fmt.Println("")
	fmt.Println("==========================================")
	fmt.Println(" PLATAFORMA DE INTERCAMBIO DE LIBROS")
	fmt.Println("==========================================")
	fmt.Println("")
	fmt.Println("Servidor: http://localhost:8080")
	fmt.Println("")
	fmt.Println("Entidades del modelo:")
	fmt.Println("1. Usuario")
	fmt.Println("2. Libro")
	fmt.Println("3. Autor")
	fmt.Println("4. Categoria")
	fmt.Println("5. LibroAutor")
	fmt.Println("6. LibroCategoria")
	fmt.Println("7. EjemplarLibro")
	fmt.Println("8. SolicitudIntercambio")
	fmt.Println("9. Intercambio")
	fmt.Println("10. DetalleIntercambio")
	fmt.Println("11. Entrega")
	fmt.Println("12. Resena")
	fmt.Println("")
	fmt.Println("Endpoints disponibles:")
	fmt.Println("POST   /login")
	fmt.Println("GET    /intercambios")
	fmt.Println("GET    /intercambios/{id}")
	fmt.Println("POST   /intercambios")
	fmt.Println("PUT    /intercambios/{id}")
	fmt.Println("PATCH  /intercambios/{id}")
	fmt.Println("DELETE /intercambios/{id}")
	fmt.Println("GET    /salud")
	fmt.Println("")
	fmt.Println("Permisos:")
	fmt.Println("GET: administrador y usuario")
	fmt.Println("POST, PUT, PATCH, DELETE: solo administrador")
	fmt.Println("")
	fmt.Println("Servidor iniciado correctamente")

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
	}
}
