package libros

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Manejador struct {
	DB *gorm.DB
}

func (m Manejador) ListarLibros(w http.ResponseWriter, r *http.Request) {
	var libros []Libro

	query := m.DB.Debug().Preload("Prestamos")

	estado := r.URL.Query().Get("estado")

	if estado != "" {
		if !EstadosValidos[estado] {
			http.Error(w, "Estado invalido", http.StatusUnprocessableEntity)
			return
		}

		query = query.Where("estado = ?", estado)
	}

	if err := query.Find(&libros).Error; err != nil {
		http.Error(w, "Error al consultar los libros", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(libros)
}

func (m Manejador) CrearLibro(w http.ResponseWriter, r *http.Request) {
	var libro Libro

	if err := json.NewDecoder(r.Body).Decode(&libro); err != nil {
		http.Error(w, "JSON invalido", http.StatusBadRequest)
		return
	}

	if libro.Titulo == "" || libro.Autor == "" {
		http.Error(w, "Titulo y autor son obligatorios", http.StatusUnprocessableEntity)
		return
	}

	if libro.Titulo == libro.Autor {
		http.Error(w, "Titulo y autor no pueden ser iguales", http.StatusUnprocessableEntity)
		return
	}

	if !EstadosValidos[libro.Estado] {
		http.Error(w, "Estado invalido", http.StatusUnprocessableEntity)
		return
	}

	if err := m.DB.Create(&libro).Error; err != nil {
		http.Error(w, "Error al crear el libro", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(libro)
}
func (m Manejador) ObtenerLibro(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, err := strconv.ParseUint(id, 10, 64); err != nil {
		http.Error(w, "ID invalido", http.StatusBadRequest)
		return
	}

	var libro Libro

	if err := m.DB.First(&libro, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Libro no encontrado", http.StatusNotFound)
			return
		}

		http.Error(w, "Error al consultar el libro", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(libro)
}

func (m Manejador) ActualizarLibro(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var libro Libro

	if err := m.DB.First(&libro, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Libro no encontrado", http.StatusNotFound)
			return
		}

		http.Error(w, "Error al consultar el libro", http.StatusInternalServerError)
		return
	}

	var datos Libro

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		http.Error(w, "JSON invalido", http.StatusBadRequest)
		return
	}

	if datos.Titulo == "" || datos.Autor == "" {
		http.Error(w, "Titulo y autor son obligatorios", http.StatusUnprocessableEntity)
		return
	}

	if !EstadosValidos[datos.Estado] {
		http.Error(w, "Estado invalido", http.StatusUnprocessableEntity)
		return
	}

	libro.Titulo = datos.Titulo
	libro.Autor = datos.Autor
	libro.Estado = datos.Estado

	if err := m.DB.Save(&libro).Error; err != nil {
		http.Error(w, "Error al actualizar el libro", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(libro)
}

func (m Manejador) EliminarLibro(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var libro Libro

	if err := m.DB.First(&libro, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Libro no encontrado", http.StatusNotFound)
			return
		}

		http.Error(w, "Error al consultar el libro", http.StatusInternalServerError)
		return
	}

	var cantidad int64

	if err := m.DB.Model(&Prestamo{}).Where("libro_id = ?", id).Count(&cantidad).Error; err != nil {
		http.Error(w, "Error al consultar los prestamos", http.StatusInternalServerError)
		return
	}

	if cantidad > 0 {
		http.Error(w, "No se puede eliminar un libro con prestamos", http.StatusUnprocessableEntity)
		return
	}

	if err := m.DB.Delete(&libro).Error; err != nil {
		http.Error(w, "Error al eliminar el libro", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (m Manejador) Rutas(r chi.Router) {
	r.Get("/libros", m.ListarLibros)
	r.Post("/libros", m.CrearLibro)
	r.Get("/libros/{id}", m.ObtenerLibro)
	r.Put("/libros/{id}", m.ActualizarLibro)
	r.Delete("/libros/{id}", m.EliminarLibro)
}
