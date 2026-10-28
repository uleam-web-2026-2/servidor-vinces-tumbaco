package intercambios

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSinTokenLasRutasProtegidasResponden401(t *testing.T) {
	handler := &Handler{}

	rutas := RutasProtegidas(handler)

	pruebas := []struct {
		nombre string
		metodo string
		ruta   string
		cuerpo string
	}{
		{"Listar intercambios", http.MethodGet, "/intercambios", ""},
		{"Obtener intercambio", http.MethodGet, "/intercambios/1", ""},
		{"Crear intercambio", http.MethodPost, "/intercambios", `{}`},
		{"Actualizar intercambio", http.MethodPut, "/intercambios/1", `{}`},
		{"Cambiar estado", http.MethodPatch, "/intercambios/1", `{"estado":"aceptado"}`},
		{"Eliminar intercambio", http.MethodDelete, "/intercambios/1", ""},
	}

	for _, prueba := range pruebas {
		t.Run(prueba.nombre, func(t *testing.T) {
			peticion := httptest.NewRequest(
				prueba.metodo,
				prueba.ruta,
				strings.NewReader(prueba.cuerpo),
			)

			respuesta := httptest.NewRecorder()

			rutas(respuesta, peticion)

			if respuesta.Code != http.StatusUnauthorized {
				t.Errorf(
					"%s %s: se esperaba 401, se obtuvo %d",
					prueba.metodo,
					prueba.ruta,
					respuesta.Code,
				)
			}
		})
	}
}
