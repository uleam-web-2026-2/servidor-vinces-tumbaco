package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/uleam-web-2026-2/KEMA/internal/avisos"
	"github.com/uleam-web-2026-2/KEMA/internal/config"
	"github.com/uleam-web-2026-2/KEMA/internal/intercambios"
	"github.com/uleam-web-2026-2/KEMA/internal/repository"
	"github.com/uleam-web-2026-2/KEMA/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// validarBaseDatos comprueba que la aplicaciÃƒÆ’Ã‚Â³n utilice
// exclusivamente la base de datos de recuperaciÃƒÆ’Ã‚Â³n.
func validarBaseDatos(databaseURL string) error {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return fmt.Errorf("la direcciÃƒÆ’Ã‚Â³n de PostgreSQL no es vÃƒÆ’Ã‚Â¡lida")
	}

	host := u.Hostname()
	puerto := u.Port()

	conexionLocal := (host == "localhost" || host == "127.0.0.1") && puerto == "5434"
	conexionDocker := host == "db" && puerto == "5432"

	if !conexionLocal && !conexionDocker {
		return fmt.Errorf("conexiÃƒÆ’Ã‚Â³n no permitida: use localhost:5434 o db:5432")
	}

	if u.Path != "/intercambio_prueba" {
		return fmt.Errorf("por seguridad, se requiere la base de datos intercambio_prueba")
	}

	return nil
}

// porIP identifica al cliente mediante su direcciÃƒÆ’Ã‚Â³n IP.
func porIP(r *http.Request) (string, error) {
	ip := middleware.GetClientIP(r.Context())
	return httprate.CanonicalizeIP(ip), nil
}

func main() {
	// Cargar configuraciÃƒÆ’Ã‚Â³n.
	cfg, err := config.Cargar()
	if err != nil {
		log.Fatal("Error de configuraciÃƒÆ’Ã‚Â³n: ", err)
	}

	// Verificar la base de datos antes de conectarse.
	if err := validarBaseDatos(cfg.DatabaseURL); err != nil {
		log.Fatal("ConexiÃƒÆ’Ã‚Â³n bloqueada: ", err)
	}

	// Conectar con PostgreSQL.
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar con PostgreSQL: ", err)
	}

	fmt.Println("ConexiÃƒÆ’Ã‚Â³n con PostgreSQL correcta")

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
		if strings.Contains(err.Error(), "must be owner of table") || strings.Contains(err.Error(), "42501") {
			log.Println("Aviso: omitiendo cambios DDL de AutoMigrate por permisos del usuario de base de datos:", err)
		} else {
			log.Fatal("Error al preparar las tablas: ", err)
		}
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

	// Administrador unico de conexiones y avisos WebSocket.
	administradorAvisos := avisos.Nuevo()

	// CAPA 2: Service.
	servicio := services.NuevoIntercambiosService(repositorio, administradorAvisos)

	// CAPA 3: Handler.
	handler := &intercambios.Handler{
		Servicio: servicio,
	}

	// Ruta para establecer conexiones WebSocket.
	http.HandleFunc("/ws", administradorAvisos.Conectar)

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

	// Limitar el inicio de sesiÃƒÆ’Ã‚Â³n a 5 intentos por minuto y por IP.
	loginLimitado := httprate.LimitBy(
		5,
		time.Minute,
		porIP,
	)(http.HandlerFunc(intercambios.Login))

	http.Handle("/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			intercambios.Login(w, r)
			return
		}

		loginLimitado.ServeHTTP(w, r)
	}))

	rutasIntercambios := intercambios.RutasProtegidas(handler)

	http.HandleFunc("/intercambios", rutasIntercambios)
	http.HandleFunc("/intercambios/", rutasIntercambios)

	// ConfiguraciÃƒÆ’Ã‚Â³n CORS.
	manejadorCORS := cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:5173",
		},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
		},
		AllowedHeaders: []string{
			"Authorization",
			"Content-Type",
		},
		MaxAge: 300,
	})

	// Configurar servidor HTTP.
	servidor := &http.Server{
		Addr:         ":" + cfg.Puerto,
		ReadTimeout:  cfg.TiempoEspera,
		WriteTimeout: cfg.TiempoEspera,
		Handler: manejadorCORS(
			middleware.ClientIPFromRemoteAddr(
				middleware.Logger(http.DefaultServeMux),
			),
		),
	}

	fmt.Println("Servidor iniciado en el puerto:", cfg.Puerto)

	// Iniciar servidor.
	if err := servidor.ListenAndServe(); err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
	}
}


type notificadorCreacionRol struct {
hub *avisos.Avisos
rol string
}

func (n *notificadorCreacionRol) NotificarCreacion(intercambio *intercambios.Intercambio, creadorID uint) {
if n == nil || n.hub == nil || intercambio == nil {
return
}
n.hub.EnviarPorRol(
n.rol,
avisos.Aviso{
ID:     intercambio.ID,
Estado: intercambio.Estado,
Evento: "intercambio_creado",
},
strconv.FormatUint(uint64(creadorID), 10),
)
}
