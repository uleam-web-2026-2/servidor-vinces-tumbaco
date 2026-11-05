package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uleam-web-2026-2/KEMA/internal/almacen/filas"
	"github.com/uleam-web-2026-2/KEMA/internal/intercambios"
	"gorm.io/gorm"
)

// IntercambiosRepository contiene las operaciones de acceso a datos combinando GORM y sqlc.
type IntercambiosRepository struct {
	DB      *gorm.DB
	queries *filas.Queries
}

// NuevoIntercambiosRepository crea un repositorio con soporte para GORM y sqlc.
func NuevoIntercambiosRepository(db *gorm.DB, sqlDB *sql.DB) *IntercambiosRepository {
	return &IntercambiosRepository{
		DB:      db,
		queries: filas.New(sqlDB),
	}
}

// Listar obtiene todos los intercambios usando GORM.
func (repo *IntercambiosRepository) Listar() ([]intercambios.Intercambio, error) {
	var lista []intercambios.Intercambio

	err := repo.DB.
		Preload("UsuarioOfrece").
		Preload("UsuarioRecibe").
		Find(&lista).Error

	return lista, err
}

// Obtener busca un intercambio por su identificador utilizando sqlc.
func (repo *IntercambiosRepository) Obtener(id uint) (*intercambios.Intercambio, error) {
	ctx := context.Background()
	
	// sqlc requiere el tipo adecuado (ej. int32)
	row, err := repo.queries.BuscarIntercambio(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	// Mapeamos el resultado de sqlc al modelo de GORM/Intercambio existente
	intercambio := &intercambios.Intercambio{
		ID:     uint(row.ID),
		Estado: row.Estado,
	}

	return intercambio, nil
}

// BuscarEjemplar obtiene un ejemplar por su identificador usando GORM.
func (repo *IntercambiosRepository) BuscarEjemplar(id uint) (*intercambios.EjemplarLibro, error) {
	var ejemplar intercambios.EjemplarLibro

	err := repo.DB.First(&ejemplar, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &ejemplar, nil
}

// Crear registra el intercambio y su detalle en una transacción usando GORM.
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

// Actualizar guarda los cambios de un intercambio usando GORM.
func (repo *IntercambiosRepository) Actualizar(intercambio *intercambios.Intercambio) error {
	return repo.DB.Save(intercambio).Error
}

// CambiarEstado actualiza el estado de un intercambio utilizando sqlc.
func (repo *IntercambiosRepository) CambiarEstado(intercambio *intercambios.Intercambio) error {
	ctx := context.Background()

	rowsAffected, err := repo.queries.GuardarIntercambio(ctx, filas.GuardarIntercambioParams{
		ID:     int32(intercambio.ID),
		Estado: intercambio.Estado,
	})
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Eliminar borra el detalle y el intercambio en una transacción usando GORM.
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