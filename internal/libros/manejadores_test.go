package libros

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func enrutador() http.Handler {
	r := chi.NewRouter()
	(&Manejador{DB: nil}).Rutas(r)
	return r
}

func TestCrearLibro(t *testing.T) {
	casos := []struct {
		nombre string
		cuerpo string
		codigo int
	}{
		{
			nombre: "JSON roto responde 400",
			cuerpo: `{"titulo": "Libro"`,
			codigo: http.StatusBadRequest,
		},
		{
			nombre: "titulo vacio responde 422",
			cuerpo: `{"titulo": "", "autor": "Autor", "estado": "disponible"}`,
			codigo: http.StatusUnprocessableEntity,
		},
		{
			nombre: "estado invalido responde 422",
			cuerpo: `{"titulo": "Libro", "autor": "Autor", "estado": "urgente"}`,
			codigo: http.StatusUnprocessableEntity,
		},
		{
			nombre: "autor vacio responde 422",
			cuerpo: `{"titulo": "Libro", "autor": "", "estado": "disponible"}`,
			codigo: http.StatusUnprocessableEntity,
		},
		{
			nombre: "titulo y autor iguales responden 422",
			cuerpo: `{"titulo": "Gabriel Garcia Marquez", "autor": "Gabriel Garcia Marquez", "estado": "disponible"}`,
			codigo: http.StatusUnprocessableEntity,
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			peticion := httptest.NewRequest(
				http.MethodPost,
				"/libros",
				strings.NewReader(caso.cuerpo),
			)

			grabadora := httptest.NewRecorder()

			enrutador().ServeHTTP(grabadora, peticion)

			if grabadora.Code != caso.codigo {
				t.Fatalf("se esperaba %d y llego %d", caso.codigo, grabadora.Code)
			}
		})
	}
}
func TestObtenerLibroIDInvalido(t *testing.T) {
	peticion := httptest.NewRequest(
		http.MethodGet,
		"/libros/abc",
		nil,
	)

	grabadora := httptest.NewRecorder()

	enrutador().ServeHTTP(grabadora, peticion)

	if grabadora.Code != http.StatusBadRequest {
		t.Fatalf("se esperaba %d y llego %d", http.StatusBadRequest, grabadora.Code)
	}
}
