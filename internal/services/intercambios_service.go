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

// NotificadorCreacion abstrae el aviso al crear un registro sin que el servicio conozca roles ni transporte.
type NotificadorCreacion interface {
	NotificarCreacion(intercambio *intercambios.Intercambio, creadorID uint)
}

type IntercambiosService struct {
	Repositorio         *repository.IntercambiosRepository
	Avisos              EmisorAvisos
	NotificadorCreacion NotificadorCreacion
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

	if s.NotificadorCreacion != nil {
		s.NotificadorCreacion.NotificarCreacion(intercambio, input.UsuarioSolicitanteID)
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

	if input.Observaciones != "" {
		intercambio.Observaciones = input.Observaciones
	}

	if err := s.Repositorio.Actualizar(intercambio); err != nil {
		return nil, err
	}

	if s.Avisos != nil && intercambio.UsuarioRecibeID != nil {
		s.Avisos.Enviar(
			strconv.FormatUint(uint64(*intercambio.UsuarioRecibeID), 10),
			avisos.Aviso{
				ID:     intercambio.ID,
				Estado: intercambio.Estado,
			},
		)
	}

	return intercambio, nil
}

func (s *IntercambiosService) CambiarEstado(
	id uint,
	nuevoEstado string,
) (*intercambios.Intercambio, error) {
	intercambio, err := s.Repositorio.Obtener(id)
	if err != nil {
		return nil, err
	}

	if nuevoEstado != intercambios.EstadoPendiente &&
		nuevoEstado != intercambios.EstadoAceptado &&
		nuevoEstado != intercambios.EstadoEntregado &&
		nuevoEstado != intercambios.EstadoCancelado {
		return nil, errors.New("estado no permitido")
	}

	transicionValida := false
	switch intercambio.Estado {
	case intercambios.EstadoPendiente:
		transicionValida = nuevoEstado == intercambios.EstadoAceptado ||
			nuevoEstado == intercambios.EstadoCancelado
	case intercambios.EstadoAceptado:
		transicionValida = nuevoEstado == intercambios.EstadoEntregado ||
			nuevoEstado == intercambios.EstadoCancelado
	}

	if !transicionValida {
		return nil, errors.New("transicion de estado no permitida")
	}

	ahora := time.Now()
	intercambio.Estado = nuevoEstado

	switch nuevoEstado {
	case intercambios.EstadoAceptado:
		intercambio.FechaConfirmacion = &ahora
	case intercambios.EstadoEntregado:
		intercambio.FechaFinalizacion = &ahora
	}

	if err := s.Repositorio.CambiarEstado(intercambio); err != nil {
		return nil, err
	}

	if s.Avisos != nil {
		notificados := make(map[string]struct{})
		for _, idPtr := range []*uint{intercambio.UsuarioOfreceID, intercambio.UsuarioRecibeID} {
			if idPtr == nil || *idPtr == 0 {
				continue
			}
			clave := strconv.FormatUint(uint64(*idPtr), 10)
			if _, ya := notificados[clave]; ya {
				continue
			}
			notificados[clave] = struct{}{}
			s.Avisos.Enviar(clave, avisos.Aviso{
				ID:     intercambio.ID,
				Estado: intercambio.Estado,
				Evento: "estado_cambiado",
			})
		}
	}

	return intercambio, nil
}

func (s *IntercambiosService) Eliminar(id uint) error {
	intercambio, err := s.Repositorio.Obtener(id)
	if err != nil {
		return err
	}

	if intercambio.Estado != intercambios.EstadoPendiente &&
		intercambio.Estado != intercambios.EstadoCancelado {
		return errors.New("solo se pueden eliminar intercambios pendientes o cancelados")
	}

	return s.Repositorio.Eliminar(id)
}
