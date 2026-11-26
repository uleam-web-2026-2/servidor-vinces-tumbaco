package intercambios

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Credenciales struct {
	Usuario string `json:"usuario"`
	Clave   string `json:"clave"`
}

func extraerUsuarioIDDesdeClaims(claims jwt.MapClaims) uint {
	if rawID, existe := claims["id"]; existe && rawID != nil {
		switch v := rawID.(type) {
		case float64:
			if v > 0 {
				return uint(v)
			}
		case int:
			if v > 0 {
				return uint(v)
			}
		case string:
			if n, err := strconv.Atoi(v); err == nil && n > 0 {
				return uint(n)
			}
		}
	}

	if subNum, ok := claims["sub"].(float64); ok && subNum > 0 {
		return uint(subNum)
	}

	if subStr, ok := claims["sub"].(string); ok && subStr != "" {
		if n, err := strconv.Atoi(subStr); err == nil && n > 0 {
			return uint(n)
		}
	}

	return 0
}

// ExtraerSolicitanteDesdeRequest obtiene el ID del usuario autenticado desde el JWT del request.
func ExtraerSolicitanteDesdeRequest(r *http.Request) uint {
	autorizacion := r.Header.Get("Authorization")
	partes := strings.Fields(autorizacion)
	if len(partes) != 2 || !strings.EqualFold(partes[0], "Bearer") {
		return 0
	}

	clave := os.Getenv("JWT_SECRET")
	if len(clave) < 32 {
		return 0
	}

	token, err := jwt.Parse(
		partes[1],
		func(t *jwt.Token) (interface{}, error) {
			return []byte(clave), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil || !token.Valid {
		return 0
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0
	}

	return extraerUsuarioIDDesdeClaims(claims)
}

func claveValida(ingresada string, nombresEnv ...string) bool {
	if ingresada == "" {
		return false
	}

	for _, envName := range nombresEnv {
		esperada := os.Getenv(envName)
		if esperada != "" && subtle.ConstantTimeCompare([]byte(ingresada), []byte(esperada)) == 1 {
			return true
		}
	}

	return false
}

// Login valida los usuarios de demostracion y entrega un JWT.
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responderError(w, http.StatusMethodNotAllowed,
			"metodo no permitido")
		return
	}

	var datos Credenciales

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		responderError(w, http.StatusBadRequest, "JSON invalido")
		return
	}

	var rol string
	var usuarioID uint

	switch {
	case datos.Usuario == "kenyis_admin" && claveValida(datos.Clave, "DEMO_ADMIN_PASSWORD"):
		rol = "administrador"
		usuarioID = 3

	case datos.Usuario == "majo_user" && claveValida(datos.Clave, "DEMO_USER_PASSWORD"):
		rol = "usuario"
		usuarioID = 1

	case datos.Usuario == "carlos_user" && claveValida(datos.Clave, "DEMO_SECOND_USER_PASSWORD", "DEMO_USER_PASSWORD"):
		rol = "usuario"
		usuarioID = 2

	default:
		responderError(w, http.StatusUnauthorized,
			"credenciales incorrectas")
		return
	}

	clave := os.Getenv("JWT_SECRET")
	if len(clave) < 32 {
		responderError(w, http.StatusInternalServerError,
			"JWT_SECRET no esta configurado correctamente")
		return
	}

	claims := jwt.MapClaims{
		"id":  usuarioID,
		"sub": datos.Usuario,
		"rol": rol,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(30 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	firmado, err := token.SignedString([]byte(clave))
	if err != nil {
		responderError(w, http.StatusInternalServerError,
			"no se pudo generar el token")
		return
	}

	responderJSON(w, http.StatusOK, map[string]string{
		"token": firmado,
	})
}

// Proteger exige un token valido y un rol permitido.
func Proteger(siguiente http.HandlerFunc, roles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		autorizacion := r.Header.Get("Authorization")
		partes := strings.Fields(autorizacion)

		if len(partes) != 2 ||
			!strings.EqualFold(partes[0], "Bearer") {
			responderError(w, http.StatusUnauthorized,
				"token requerido")
			return
		}

		clave := os.Getenv("JWT_SECRET")
		if len(clave) < 32 {
			responderError(w, http.StatusInternalServerError,
				"JWT_SECRET no esta configurado correctamente")
			return
		}

		token, err := jwt.Parse(
			partes[1],
			func(t *jwt.Token) (interface{}, error) {
				return []byte(clave), nil
			},
			jwt.WithValidMethods([]string{"HS256"}),
		)

		if err != nil || !token.Valid {
			responderError(w, http.StatusUnauthorized,
				"token invalido o vencido")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			responderError(w, http.StatusUnauthorized,
				"token invalido")
			return
		}

		rol, ok := claims["rol"].(string)
		if !ok {
			responderError(w, http.StatusUnauthorized,
				"rol invalido")
			return
		}

		for _, permitido := range roles {
			if rol == permitido {
				siguiente(w, r)
				return
			}
		}

		responderError(w, http.StatusForbidden,
			"no tienes permiso para esta operacion")
	}
}
