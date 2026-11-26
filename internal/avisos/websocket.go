package avisos

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/golang-jwt/jwt/v5"
)

func extraerClaveUsuario(claims jwt.MapClaims) string {
	if rawID, existe := claims["id"]; existe && rawID != nil {
		switch v := rawID.(type) {
		case float64:
			if v > 0 {
				return strconv.Itoa(int(v))
			}
		case int:
			if v > 0 {
				return strconv.Itoa(v)
			}
		case string:
			if v != "" {
				return v
			}
		}
	}

	if sub, ok := claims["sub"].(string); ok && sub != "" {
		return sub
	}

	if subNum, ok := claims["sub"].(float64); ok && subNum > 0 {
		return strconv.Itoa(int(subNum))
	}

	return ""
}

// Conectar valida el token y establece una conexion WebSocket.
func (a *Avisos) Conectar(w http.ResponseWriter, r *http.Request) {
	tokenTexto := r.URL.Query().Get("token")
	if tokenTexto == "" {
		http.Error(w, "token requerido", http.StatusUnauthorized)
		return
	}

	clave := os.Getenv("JWT_SECRET")
	if len(clave) < 32 {
		http.Error(w, "configuracion JWT incorrecta", http.StatusInternalServerError)
		return
	}

	token, err := jwt.Parse(
		tokenTexto,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(clave), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil || !token.Valid {
		http.Error(w, "token invalido o vencido", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "token invalido", http.StatusUnauthorized)
		return
	}

	usuario := extraerClaveUsuario(claims)
	if usuario == "" {
		http.Error(w, "usuario invalido", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"http://localhost:5173"},
	})
	if err != nil {
		return
	}
	defer conn.CloseNow()

	canal := a.Suscribir(usuario)
	defer a.Quitar(usuario, canal)

	ctx := conn.CloseRead(context.Background())

	for {
		select {
		case <-ctx.Done():
			return

		case aviso, abierto := <-canal:
			if !abierto {
				return
			}

			if err := wsjson.Write(ctx, conn, aviso); err != nil {
				return
			}
		}
	}
}
