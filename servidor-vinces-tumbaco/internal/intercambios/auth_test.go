package intercambios

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestLogin(t *testing.T) {
	t.Setenv("JWT_SECRET", strings.Repeat("a", 48))
	t.Setenv("DEMO_ADMIN_PASSWORD", "clave-admin-prueba")
	t.Setenv("DEMO_USER_PASSWORD", "clave-usuario-prueba")

	pruebas := []struct {
		nombre   string
		cuerpo   string
		esperado int
		rol      string
	}{
		{
			"Login administrador",
			`{"usuario":"admin","clave":"clave-admin-prueba"}`,
			200,
			"admin",
		},
		{
			"Login usuario",
			`{"usuario":"usuario","clave":"clave-usuario-prueba"}`,
			200,
			"usuario",
		},
		{
			"Contrasena incorrecta",
			`{"usuario":"admin","clave":"incorrecta"}`,
			401,
			"",
		},
		{
			"Usuario desconocido",
			`{"usuario":"desconocido","clave":"123"}`,
			401,
			"",
		},
		{
			"JSON incorrecto",
			`{incorrecto}`,
			400,
			"",
		},
	}

	for _, prueba := range pruebas {
		t.Run(prueba.nombre, func(t *testing.T) {
			peticion := httptest.NewRequest(
				http.MethodPost,
				"/login",
				strings.NewReader(prueba.cuerpo),
			)

			respuesta := httptest.NewRecorder()
			Login(respuesta, peticion)

			if respuesta.Code != prueba.esperado {
				t.Fatalf(
					"Codigo obtenido: %d; esperado: %d",
					respuesta.Code,
					prueba.esperado,
				)
			}

			if prueba.rol != "" {
				var datos struct {
					Token string `json:"token"`
				}

				if err := json.Unmarshal(respuesta.Body.Bytes(), &datos); err != nil {
					t.Fatal(err)
				}

				if datos.Token == "" {
					t.Fatal("No se genero el token JWT")
				}

				claims := jwt.MapClaims{}

				_, err := jwt.ParseWithClaims(
					datos.Token,
					claims,
					func(token *jwt.Token) (interface{}, error) {
						return []byte(strings.Repeat("a", 48)), nil
					},
					jwt.WithValidMethods([]string{"HS256"}),
				)

				if err != nil {
					t.Fatalf("Token invalido: %v", err)
				}

				if claims["rol"] != prueba.rol {
					t.Errorf(
						"Rol obtenido: %v; esperado: %s",
						claims["rol"],
						prueba.rol,
					)
				}
			}
		})
	}
}

func TestProteger(t *testing.T) {
	clave := strings.Repeat("b", 48)
	t.Setenv("JWT_SECRET", clave)

	generarToken := func(rol string) string {
		token := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			jwt.MapClaims{"rol": rol},
		)

		firmado, err := token.SignedString([]byte(clave))
		if err != nil {
			t.Fatal(err)
		}

		return firmado
	}

	rutaProtegida := Proteger(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
		"admin",
	)

	pruebas := []struct {
		nombre         string
		autorizacion   string
		codigoEsperado int
	}{
		{"Sin token", "", 401},
		{"Token incorrecto", "Bearer token-invalido", 401},
		{"Usuario sin permisos", "Bearer " + generarToken("usuario"), 403},
		{"Administrador autorizado", "Bearer " + generarToken("admin"), 200},
	}

	for _, prueba := range pruebas {
		t.Run(prueba.nombre, func(t *testing.T) {
			peticion := httptest.NewRequest(
				http.MethodGet,
				"/intercambios",
				nil,
			)

			if prueba.autorizacion != "" {
				peticion.Header.Set(
					"Authorization",
					prueba.autorizacion,
				)
			}

			respuesta := httptest.NewRecorder()
			rutaProtegida(respuesta, peticion)

			if respuesta.Code != prueba.codigoEsperado {
				t.Errorf(
					"Codigo obtenido: %d; esperado: %d",
					respuesta.Code,
					prueba.codigoEsperado,
				)
			}
		})
	}
}
