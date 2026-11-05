package repository

import (
	"errors"

	"github.com/uleam-web-2026-2/KEMA/internal/intercambios"
	"gorm.io/gorm"
)

// IntercambiosRepository contiene las operaciones de acceso
// a la base de datos para los intercambios.
type IntercambiosRepository struct {
	DB *gorm.DB
}

// NuevoIntercambiosRepository crea un repositorio.
func NuevoIntercambiosRepository(db *gorm.DB) *IntercambiosRepository {
	return &IntercambiosRepository{
		DB: db,
	}
}

// Listar obtiene todos los intercambios.
func (repo *IntercambiosRepository) Listar() ([]intercambios.Intercambio, error) {
	var lista []intercambios.Intercambio

	err := repo.DB.
		Preload("UsuarioOfrece").
		Preload("UsuarioRecibe").
		Find(&lista).Error

	return lista, err
}

// Obtener busca un intercambio por su identificador.
func (repo *IntercambiosRepository) Obtener(id uint) (*intercambios.Intercambio, error) {
	var intercambio intercambios.Intercambio

	err := repo.DB.
		Preload("UsuarioOfrece").
		Preload("UsuarioRecibe").
		First(&intercambio, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &intercambio, nil
}

// BuscarEjemplar obtiene un ejemplar por su identificador.
func (repo *IntercambiosRepository) BuscarEjemplar(id uint) (*intercambios.EjemplarLibro, error) {
	var ejemplar intercambios.EjemplarLibro

	err := repo.DB.First(&ejemplar, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &ejemplar, nil
}

// Crear registra el intercambio y su detalle en una transacción.
func (repo *IntercambiosRepository) Crear(
	intercambio *intercambios.Intercambio,
	detalle *intercambios.DetalleIntercambio,
) error {
	return repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(intercambio).Error; err != nil {
			return err
		}

		detalle.IntercambioID = intercambio.ID

		return tx.Create(detalle).Error
	})
}

// Actualizar guarda los cambios de un intercambio.
func (repo *IntercambiosRepository) Actualizar(intercambio *intercambios.Intercambio) error {
	return repo.DB.Save(intercambio).Error
}

// CambiarEstado actualiza el estado de un intercambio.
func (repo *IntercambiosRepository) CambiarEstado(intercambio *intercambios.Intercambio) error {
	return repo.DB.Save(intercambio).Error
}

// Eliminar borra el detalle y el intercambio en una transacción.
func (repo *IntercambiosRepository) Eliminar(id uint) error {
	return repo.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("intercambio_id = ?", id).
			Delete(&intercambios.DetalleIntercambio{}).Error; err != nil {
			return err
		}

		resultado := tx.Delete(&intercambios.Intercambio{}, "id = ?", id)

		if resultado.Error != nil {
			return resultado.Error
		}

		if resultado.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
}
