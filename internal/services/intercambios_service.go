package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/uleam-web-2026-2/KEMA/internal/avisos"
	"github.com/uleam-web-2026-2/KEMA/internal/intercambios"
	"github.com/uleam-web-2026-2/KEMA/internal/repository"
	"gorm.io/gorm"
)

// EmisorAvisos permite enviar notificaciones de dominio sin acoplar el servicio a transporte.
type EmisorAvisos interface {
	Enviar(usuario string, aviso avisos.Aviso)
}

// EmisorPorRol permite notificar eventos de negocio a conexiones filtradas por rol.
type EmisorPorRol interface {
	EnviarPorRol(rol string, aviso avisos.Aviso)
}

type IntercambiosService struct {
	Repositorio *repository.IntercambiosRepository
	Avisos      EmisorAvisos
}

func NuevoIntercambiosService(
	repo *repository.IntercambiosRepository,
	gestorAvisos ...EmisorAvisos,
) *IntercambiosService {
	var emisor EmisorAvisos
	if len(gestorAvisos) > 0 {
		emisor = gestorAvisos[0]
	}

	return &IntercambiosService{
		Repositorio: repo,
		Avisos:      emisor,
	}
}

func (s *IntercambiosService) Listar() ([]intercambios.Intercambio, error) {
	return s.Repositorio.Listar()
}

func (s *IntercambiosService) Obtener(id uint) (*intercambios.Intercambio, error) {
	return s.Repositorio.Obtener(id)
}

func (s *IntercambiosService) Crear(
	input intercambios.CrearIntercambioInput,
) (*intercambios.Intercambio, error) {
	if input.UsuarioSolicitanteID == 0 ||
		input.EjemplarOfrecidoID == 0 ||
		input.EjemplarDeseadoID == 0 {
		return nil, errors.New("debe indicar el usuario y los dos ejemplares")
	}

	ofrecido, err := s.Repositorio.BuscarEjemplar(input.EjemplarOfrecidoID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el ejemplar ofrecido no existe")
	}
	if err != nil {
		return nil, err
	}

	deseado, err := s.Repositorio.BuscarEjemplar(input.EjemplarDeseadoID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el ejemplar deseado no existe")
	}
	if err != nil {
		return nil, err
	}

	if ofrecido.UsuarioID != input.UsuarioSolicitanteID {
		return nil, errors.New("el ejemplar ofrecido no pertenece al usuario solicitante")
	}

	if ofrecido.ID == deseado.ID {
		return nil, errors.New("no se puede intercambiar un ejemplar por si mismo")
	}

	if ofrecido.Disponibilidad != "disponible" ||
		deseado.Disponibilidad != "disponible" {
		return nil, errors.New("uno de los ejemplares no esta disponible")
	}

	ahora := time.Now()
	usuarioOfreceID := ofrecido.UsuarioID
	usuarioRecibeID := deseado.UsuarioID

	intercambio := &intercambios.Intercambio{
		UsuarioOfreceID: &usuarioOfreceID,
		UsuarioRecibeID: &usuarioRecibeID,
		FechaSolicitud:  &ahora,
		Estado:          intercambios.EstadoPendiente,
		Observaciones:   input.Observaciones,
	}

	detalle := &intercambios.DetalleIntercambio{
		IntercambioID:       intercambio.ID,
		EjemplarEntregadoID: ofrecido.ID,
		EjemplarRecibidoID:  deseado.ID,
		Cantidad:            1,
	}

	if err := s.Repositorio.Crear(intercambio, detalle); err != nil {
		return nil, err
	}

	if emisorRol, ok := s.Avisos.(EmisorPorRol); ok && emisorRol != nil {
		emisorRol.EnviarPorRol("administrador", avisos.Aviso{
			ID:     intercambio.ID,
			Estado: intercambio.Estado,
			Evento: "creado",
		})
	}

	return intercambio, nil
}

func (s *IntercambiosService) Actualizar(
	id uint,
	input intercambios.ActualizarIntercambioInput,
) (*intercambios.Intercambio, error) {
	intercambio, err := s.Repositorio.Obtener(id)
	if err != nil {
		return nil, err
	}

	if intercambio.Estado != intercambios.EstadoPendiente {
		return nil, errors.New("solo se pueden actualizar intercambios pendientes")
	}

	if input.UsuarioOfreceID == 0 ||
		input.UsuarioRecibeID == 0 ||
		input.UsuarioOfreceID == input.UsuarioRecibeID {
		return nil, errors.New("los usuarios indicados no son validos")
	}

	intercambio.UsuarioOfreceID = &input.UsuarioOfreceID
	intercambio.UsuarioRecibeID = &input.UsuarioRecibeID
	intercambio.Observaciones = input.Observaciones

	if err := s.Repositorio.Actualizar(intercambio); err != nil {
		return nil, err
	}

	return intercambio, nil
}

func (s *IntercambiosService) CambiarEstado(
	id uint,
	nuevoEstado string,
) (*intercambios.Intercambio, error) {
	if !intercambios.EstadoValido(nuevoEstado) {
		return nil, errors.New("el estado indicado no es valido")
	}

	intercambio, err := s.Repositorio.Obtener(id)
	if err != nil {
		return nil, err
	}

	if !intercambios.TransicionPermitida(intercambio.Estado, nuevoEstado) {
		return nil, errors.New("la transicion de estado no esta permitida")
	}

	ahora := time.Now()
	intercambio.Estado = nuevoEstado

	if nuevoEstado == intercambios.EstadoAceptado {
		intercambio.FechaConfirmacion = &ahora
	}

	if nuevoEstado == intercambios.EstadoEntregado {
		intercambio.FechaFinalizacion = &ahora
	}

	// Guardar primero: nunca avisar si la persistencia falla.
	if err := s.Repositorio.CambiarEstado(intercambio); err != nil {
		return nil, err
	}

	if s.Avisos != nil {
		aviso := avisos.Aviso{
			ID:     intercambio.ID,
			Estado: intercambio.Estado,
			Evento: "estado_cambiado",
		}
		notificados := make(map[string]struct{})

		for _, participante := range []*uint{intercambio.UsuarioOfreceID, intercambio.UsuarioRecibeID} {
			if participante == nil || *participante == 0 {
				continue
			}
			clave := strconv.Itoa(int(*participante))
			if _, yaNotificado := notificados[clave]; yaNotificado {
				continue
			}
			notificados[clave] = struct{}{}
			s.Avisos.Enviar(clave, aviso)
		}
	}

	return intercambio, nil
}

func (s *IntercambiosService) Eliminar(id uint) error {
	return s.Repositorio.Eliminar(id)
}
