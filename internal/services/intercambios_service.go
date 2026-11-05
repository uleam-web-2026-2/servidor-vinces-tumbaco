package services

import (
	"errors"
	"time"

	"github.com/uleam-web-2026-2/KEMA/internal/intercambios"
	"github.com/uleam-web-2026-2/KEMA/internal/repository"
	"gorm.io/gorm"
)

type IntercambiosService struct {
	Repositorio *repository.IntercambiosRepository
}

func NuevoIntercambiosService(repo *repository.IntercambiosRepository) *IntercambiosService {
	return &IntercambiosService{
		Repositorio: repo,
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
		return nil, errors.New("no se puede intercambiar un ejemplar por sí mismo")
	}

	if ofrecido.Disponibilidad != "disponible" ||
		deseado.Disponibilidad != "disponible" {
		return nil, errors.New("uno de los ejemplares no está disponible")
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
		EjemplarEntregadoID: ofrecido.ID,
		EjemplarRecibidoID:  deseado.ID,
		Cantidad:            1,
	}

	if err := s.Repositorio.Crear(intercambio, detalle); err != nil {
		return nil, err
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
		return nil, errors.New("los usuarios indicados no son válidos")
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
		return nil, errors.New("el estado indicado no es válido")
	}

	intercambio, err := s.Repositorio.Obtener(id)
	if err != nil {
		return nil, err
	}

	if !intercambios.TransicionPermitida(intercambio.Estado, nuevoEstado) {
		return nil, errors.New("la transición de estado no está permitida")
	}

	ahora := time.Now()
	intercambio.Estado = nuevoEstado

	if nuevoEstado == intercambios.EstadoAceptado {
		intercambio.FechaConfirmacion = &ahora
	}

	if nuevoEstado == intercambios.EstadoEntregado {
		intercambio.FechaFinalizacion = &ahora
	}

	if err := s.Repositorio.CambiarEstado(intercambio); err != nil {
		return nil, err
	}

	return intercambio, nil
}

func (s *IntercambiosService) Eliminar(id uint) error {
	return s.Repositorio.Eliminar(id)
}
