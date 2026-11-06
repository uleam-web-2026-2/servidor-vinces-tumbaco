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

// Datos internos para crear un intercambio.
type CrearIntercambioInput struct {
	UsuarioSolicitanteID uint
	EjemplarOfrecidoID   uint
	EjemplarDeseadoID    uint
	Observaciones        string
}

// Datos internos para actualizar un intercambio.
type ActualizarIntercambioInput struct {
	UsuarioOfreceID uint
	UsuarioRecibeID uint
	Observaciones   string
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

// Convertir una lista de modelos a DTO de respuesta.
func listaARespuestaDTO(lista []Intercambio) []IntercambioRespuestaDTO {
	resultado := make([]IntercambioRespuestaDTO, 0, len(lista))

	for i := range lista {
		resultado = append(resultado, IntercambioARespuestaDTO(&lista[i]))
	}

	return resultado
}

// Listar todos los intercambios.
func (h *Handler) Listar(w http.ResponseWriter, r *http.Request) {
	lista, err := h.Servicio.Listar()

	if err != nil {
		responderError(w, http.StatusInternalServerError,
			"No se pudieron consultar los intercambios")
		return
	}

	responderJSON(w, http.StatusOK, listaARespuestaDTO(lista))
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

	responderJSON(w, http.StatusOK, IntercambioARespuestaDTO(intercambio))
}

// Crear un intercambio.
func (h *Handler) Crear(w http.ResponseWriter, r *http.Request) {
	var solicitud CrearIntercambioDTO

	if err := json.NewDecoder(r.Body).Decode(&solicitud); err != nil {
		responderError(w, http.StatusBadRequest,
			"El cuerpo de la solicitud no es válido")
		return
	}

	intercambio, err := h.Servicio.Crear(solicitud.AInput())

	if err != nil {
		mensaje := err.Error()
		codigo := http.StatusBadRequest

		if strings.Contains(mensaje, "no está disponible") {
			codigo = http.StatusConflict
		}

		responderError(w, codigo, mensaje)
		return
	}

	responderJSON(w, http.StatusCreated, IntercambioARespuestaDTO(intercambio))
}

// Actualizar los datos de un intercambio pendiente.
func (h *Handler) Actualizar(w http.ResponseWriter, r *http.Request, id uint) {
	var solicitud ActualizarIntercambioDTO

	if err := json.NewDecoder(r.Body).Decode(&solicitud); err != nil {
		responderError(w, http.StatusBadRequest,
			"El cuerpo de la solicitud no es válido")
		return
	}

	intercambio, err := h.Servicio.Actualizar(id, solicitud.AInput())

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

	responderJSON(w, http.StatusOK, IntercambioARespuestaDTO(intercambio))
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
	var solicitud CambiarEstadoDTO

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

	responderJSON(w, http.StatusOK, IntercambioARespuestaDTO(intercambio))
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
