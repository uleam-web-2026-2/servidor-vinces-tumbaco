package avisos

import (
	"context"
	"net/http"
	"os"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/golang-jwt/jwt/v5"
)

// Conectar valida el token y establece una conexión WebSocket.
func (a *Avisos) Conectar(w http.ResponseWriter, r *http.Request) {
	// Obtener el token enviado en la URL.
	tokenTexto := r.URL.Query().Get("token")
	if tokenTexto == "" {
		http.Error(w, "token requerido", http.StatusUnauthorized)
		return
	}

	// Validar la firma y vigencia del JWT.
	clave := os.Getenv("JWT_SECRET")
	if len(clave) < 32 {
		http.Error(w, "configuración JWT incorrecta", http.StatusInternalServerError)
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
		http.Error(w, "token inválido o vencido", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "token inválido", http.StatusUnauthorized)
		return
	}

	usuario, ok := claims["sub"].(string)
	if !ok || usuario == "" {
		http.Error(w, "usuario inválido", http.StatusUnauthorized)
		return
	}

	// Aceptar únicamente conexiones del cliente autorizado.
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:5173"},
	})
	if err != nil {
		return
	}
	defer conn.CloseNow()

	// Registrar la conexión del usuario.
	canal := a.Suscribir(usuario)
	defer a.Quitar(usuario, canal)

	// Mantener la lectura de controles de la conexión.
	ctx := conn.CloseRead(context.Background())

	// Enviar los avisos recibidos para este usuario.
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
