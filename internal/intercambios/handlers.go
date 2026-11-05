package intercambios

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// IntercambiosService define las operaciones que necesita el Handler.
// El Handler no accede directamente a PostgreSQL.
type IntercambiosService interface {
	Listar() ([]Intercambio, error)
	Obtener(id uint) (*Intercambio, error)
	Crear(input CrearIntercambioInput) (*Intercambio, error)
	Actualizar(id uint, input ActualizarIntercambioInput) (*Intercambio, error)
	CambiarEstado(id uint, nuevoEstado string) (*Intercambio, error)
	Eliminar(id uint) error
}

// Handler recibe las peticiones HTTP.
type Handler struct {
	Servicio IntercambiosService
}

// Datos necesarios para crear un intercambio.
type CrearIntercambioInput struct {
	UsuarioSolicitanteID uint   `json:"id_usuario_solicitante"`
	EjemplarOfrecidoID   uint   `json:"id_ejemplar_ofrecido"`
	EjemplarDeseadoID    uint   `json:"id_ejemplar_deseado"`
	Observaciones        string `json:"observaciones"`
}

// Datos necesarios para actualizar un intercambio.
type ActualizarIntercambioInput struct {
	UsuarioOfreceID uint   `json:"id_usuario_ofrece"`
	UsuarioRecibeID uint   `json:"id_usuario_recibe"`
	Observaciones   string `json:"observaciones"`
}

// Datos necesarios para cambiar el estado.
type solicitudEstado struct {
	Estado string `json:"estado"`
}

// Responder solicitudes con formato JSON.
func responderJSON(w http.ResponseWriter, codigo int, datos any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigo)
	_ = json.NewEncoder(w).Encode(datos)
}

// Responder errores.
func responderError(w http.ResponseWriter, codigo int, mensaje string) {
	responderJSON(w, codigo, map[string]string{
		"error": mensaje,
	})
}

// Listar todos los intercambios.
func (h *Handler) Listar(w http.ResponseWriter, r *http.Request) {
	lista, err := h.Servicio.Listar()

	if err != nil {
		responderError(w, http.StatusInternalServerError,
			"No se pudieron consultar los intercambios")
		return
	}

	responderJSON(w, http.StatusOK, lista)
}

// Obtener un intercambio por su identificador.
func (h *Handler) Obtener(w http.ResponseWriter, r *http.Request, id uint) {
	intercambio, err := h.Servicio.Obtener(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		responderError(w, http.StatusNotFound,
			"Intercambio no encontrado")
		return
	}

	if err != nil {
		responderError(w, http.StatusInternalServerError,
			"No se pudo consultar el intercambio")
		return
	}

	responderJSON(w, http.StatusOK, intercambio)
}

// Crear un intercambio.
func (h *Handler) Crear(w http.ResponseWriter, r *http.Request) {
	var solicitud CrearIntercambioInput

	if err := json.NewDecoder(r.Body).Decode(&solicitud); err != nil {
		responderError(w, http.StatusBadRequest,
			"El cuerpo de la solicitud no es válido")
		return
	}

	intercambio, err := h.Servicio.Crear(solicitud)

	if err != nil {
		mensaje := err.Error()
		codigo := http.StatusBadRequest

		if strings.Contains(mensaje, "no está disponible") {
			codigo = http.StatusConflict
		}

		responderError(w, codigo, mensaje)
		return
	}

	responderJSON(w, http.StatusCreated, intercambio)
}

// Actualizar los datos de un intercambio pendiente.
func (h *Handler) Actualizar(w http.ResponseWriter, r *http.Request, id uint) {
	var solicitud ActualizarIntercambioInput

	if err := json.NewDecoder(r.Body).Decode(&solicitud); err != nil {
		responderError(w, http.StatusBadRequest,
			"El cuerpo de la solicitud no es válido")
		return
	}

	intercambio, err := h.Servicio.Actualizar(id, solicitud)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		responderError(w, http.StatusNotFound,
			"Intercambio no encontrado")
		return
	}

	if err != nil {
		codigo := http.StatusBadRequest

		if strings.Contains(err.Error(), "pendientes") {
			codigo = http.StatusConflict
		}

		responderError(w, codigo, err.Error())
		return
	}

	responderJSON(w, http.StatusOK, intercambio)
}

// Verificar si un cambio de estado está permitido.
func TransicionPermitida(actual, nuevo string) bool {
	switch actual {
	case EstadoPendiente:
		return nuevo == EstadoAceptado || nuevo == EstadoCancelado

	case EstadoAceptado:
		return nuevo == EstadoEntregado || nuevo == EstadoCancelado

	default:
		return false
	}
}

// Cambiar el estado de un intercambio.
func (h *Handler) CambiarEstado(w http.ResponseWriter, r *http.Request, id uint) {
	var solicitud solicitudEstado

	if err := json.NewDecoder(r.Body).Decode(&solicitud); err != nil {
		responderError(w, http.StatusBadRequest,
			"El cuerpo de la solicitud no es válido")
		return
	}

	if !EstadoValido(solicitud.Estado) {
		responderError(w, http.StatusBadRequest,
			"El estado indicado no es válido")
		return
	}

	intercambio, err := h.Servicio.CambiarEstado(id, solicitud.Estado)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		responderError(w, http.StatusNotFound,
			"Intercambio no encontrado")
		return
	}

	if err != nil {
		codigo := http.StatusInternalServerError

		if strings.Contains(err.Error(), "transición") {
			codigo = http.StatusConflict
		}

		responderError(w, codigo, err.Error())
		return
	}

	responderJSON(w, http.StatusOK, intercambio)
}

// Eliminar un intercambio y su detalle.
func (h *Handler) Eliminar(w http.ResponseWriter, r *http.Request, id uint) {
	err := h.Servicio.Eliminar(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		responderError(w, http.StatusNotFound,
			"Intercambio no encontrado")
		return
	}

	if err != nil {
		responderError(w, http.StatusInternalServerError,
			"No se pudo eliminar el intercambio")
		return
	}

	responderJSON(w, http.StatusOK, map[string]string{
		"mensaje": "Intercambio eliminado correctamente",
	})
}

// Configurar las rutas de intercambios.
func (h *Handler) Rutas(w http.ResponseWriter, r *http.Request) {
	if h.Servicio == nil {
		responderError(w, http.StatusInternalServerError,
			"El servicio de intercambios no está configurado")
		return
	}

	if r.URL.Path == "/intercambios" ||
		r.URL.Path == "/intercambios/" {

		switch r.Method {
		case http.MethodGet:
			h.Listar(w, r)

		case http.MethodPost:
			h.Crear(w, r)

		default:
			responderError(w, http.StatusMethodNotAllowed,
				"Método no permitido")
		}

		return
	}

	partes := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(partes) != 2 || partes[0] != "intercambios" {
		http.NotFound(w, r)
		return
	}

	id64, err := strconv.ParseUint(partes[1], 10, 32)

	if err != nil || id64 == 0 {
		responderError(w, http.StatusBadRequest,
			"El identificador no es válido")
		return
	}

	id := uint(id64)

	switch r.Method {
	case http.MethodGet:
		h.Obtener(w, r, id)

	case http.MethodPut:
		h.Actualizar(w, r, id)

	case http.MethodPatch:
		h.CambiarEstado(w, r, id)

	case http.MethodDelete:
		h.Eliminar(w, r, id)

	default:
		responderError(w, http.StatusMethodNotAllowed,
			"Método no permitido")
	}
}
