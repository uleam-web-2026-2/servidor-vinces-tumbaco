package intercambios

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Credenciales struct {
	Usuario string `json:"usuario"`
	Clave   string `json:"clave"`
}

// Login valida usuarios de demostración y entrega un JWT.
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responderError(w, http.StatusMethodNotAllowed,
			"método no permitido")
		return
	}

	var datos Credenciales

	if err := json.NewDecoder(r.Body).Decode(&datos); err != nil {
		responderError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	var rol string

	switch {
	case datos.Usuario == "admin" &&
		subtle.ConstantTimeCompare(
			[]byte(datos.Clave),
			[]byte(os.Getenv("DEMO_ADMIN_PASSWORD")),
		) == 1 &&
		os.Getenv("DEMO_ADMIN_PASSWORD") != "":
		rol = "admin"

	case datos.Usuario == "usuario" &&
		subtle.ConstantTimeCompare(
			[]byte(datos.Clave),
			[]byte(os.Getenv("DEMO_USER_PASSWORD")),
		) == 1 &&
		os.Getenv("DEMO_USER_PASSWORD") != "":
		rol = "usuario"

	default:
		responderError(w, http.StatusUnauthorized,
			"credenciales incorrectas")
		return
	}

	clave := os.Getenv("JWT_SECRET")
	if len(clave) < 32 {
		responderError(w, http.StatusInternalServerError,
			"JWT_SECRET no está configurado correctamente")
		return
	}

	claims := jwt.MapClaims{
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

// Proteger exige un token válido y un rol permitido.
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
				"JWT_SECRET no está configurado correctamente")
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
				"token inválido o vencido")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			responderError(w, http.StatusUnauthorized,
				"token inválido")
			return
		}

		rol, ok := claims["rol"].(string)
		if !ok {
			responderError(w, http.StatusUnauthorized,
				"rol inválido")
			return
		}

		for _, permitido := range roles {
			if rol == permitido {
				siguiente(w, r)
				return
			}
		}

		responderError(w, http.StatusForbidden,
			"no tienes permiso para esta operación")
	}
}
