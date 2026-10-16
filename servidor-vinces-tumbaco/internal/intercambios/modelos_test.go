package intercambios

import "testing"

func TestEstadoValido(t *testing.T) {
	pruebas := []struct {
		nombre   string
		estado   string
		esperado bool
	}{
		{"Pendiente valido", "pendiente", true},
		{"Aceptado valido", "aceptado", true},
		{"Entregado valido", "entregado", true},
		{"Cancelado valido", "cancelado", true},
		{"Estado rechazado", "rechazado", false},
		{"Estado vacio", "", false},
		{"Estado desconocido", "desconocido", false},
	}

	for _, prueba := range pruebas {
		t.Run(prueba.nombre, func(t *testing.T) {
			resultado := EstadoValido(prueba.estado)

			if resultado != prueba.esperado {
				t.Errorf(
					"EstadoValido(%q) = %v; esperado %v",
					prueba.estado,
					resultado,
					prueba.esperado,
				)
			}
		})
	}
}
func TestTransicionPermitida(t *testing.T) {
	pruebas := []struct {
		nombre   string
		actual   string
		nuevo    string
		esperado bool
	}{
		{"Pendiente a aceptado", "pendiente", "aceptado", true},
		{"Pendiente a cancelado", "pendiente", "cancelado", true},
		{"Aceptado a entregado", "aceptado", "entregado", true},
		{"Aceptado a cancelado", "aceptado", "cancelado", true},

		{"Pendiente a entregado", "pendiente", "entregado", false},
		{"Aceptado a pendiente", "aceptado", "pendiente", false},
		{"Entregado a pendiente", "entregado", "pendiente", false},
		{"Entregado a aceptado", "entregado", "aceptado", false},
		{"Cancelado a aceptado", "cancelado", "aceptado", false},
		{"Cancelado a entregado", "cancelado", "entregado", false},
		{"Pendiente a pendiente", "pendiente", "pendiente", false},
		{"Estado desconocido", "desconocido", "aceptado", false},
	}

	for _, prueba := range pruebas {
		t.Run(prueba.nombre, func(t *testing.T) {
			resultado := TransicionPermitida(
				prueba.actual,
				prueba.nuevo,
			)

			if resultado != prueba.esperado {
				t.Errorf(
					"TransicionPermitida(%q, %q) = %v; esperado %v",
					prueba.actual,
					prueba.nuevo,
					resultado,
					prueba.esperado,
				)
			}
		})
	}
}
