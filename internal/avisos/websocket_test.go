package avisos

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/golang-jwt/jwt/v5"
)

func firmarTokenPrueba(t *testing.T, secreto string, sub string) string {
	t.Helper()
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	})
	firmado, err := tok.SignedString([]byte(secreto))
	if err != nil {
		t.Fatalf("no se pudo firmar token: %v", err)
	}
	return firmado
}

func TestConectarRechazaSinTokenOInvalido(t *testing.T) {
	clave := strings.Repeat("k", 48)
	t.Setenv("JWT_SECRET", clave)

	gestor := Nuevo()
	srv := httptest.NewServer(http.HandlerFunc(gestor.Conectar))
	defer srv.Close()

	respSinToken, err := http.Get(srv.URL + "/ws")
	if err != nil {
		t.Fatalf("error en peticion sin token: %v", err)
	}
	_ = respSinToken.Body.Close()
	if respSinToken.StatusCode != http.StatusUnauthorized {
		t.Fatalf("sin token: obtenido %d, esperado %d", respSinToken.StatusCode, http.StatusUnauthorized)
	}

	respInvalido, err := http.Get(srv.URL + "/ws?token=token-invalido")
	if err != nil {
		t.Fatalf("error en peticion con token invalido: %v", err)
	}
	_ = respInvalido.Body.Close()
	if respInvalido.StatusCode != http.StatusUnauthorized {
		t.Fatalf("token invalido: obtenido %d, esperado %d", respInvalido.StatusCode, http.StatusUnauthorized)
	}
}

func TestConectarRecibeAvisoYQuitaAlCerrar(t *testing.T) {
	clave := strings.Repeat("m", 48)
	t.Setenv("JWT_SECRET", clave)

	gestor := Nuevo()
	srv := httptest.NewServer(http.HandlerFunc(gestor.Conectar))
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws?token=" + firmarTokenPrueba(t, clave, "42")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
		HTTPHeader: http.Header{
			"Origin": []string{"http://localhost:5173"},
		},
	})
	if err != nil {
		t.Fatalf("no se pudo conectar por websocket: %v", err)
	}

	gestor.Enviar("42", Aviso{ID: 7, Estado: "aceptado"})
	gestor.Enviar("99", Aviso{ID: 8, Estado: "cancelado"})

	var recibido Aviso
	if err := wsjson.Read(ctx, conn, &recibido); err != nil {
		t.Fatalf("no se recibio aviso por websocket: %v", err)
	}
	if recibido.ID != 7 || recibido.Estado != "aceptado" {
		t.Fatalf("aviso recibido inesperado: %+v", recibido)
	}

	_ = conn.Close(websocket.StatusNormalClosure, "fin de prueba")
	time.Sleep(50 * time.Millisecond)

	gestor.mu.Lock()
	quedan := len(gestor.canales["42"])
	gestor.mu.Unlock()
	if quedan != 0 {
		t.Fatalf("se esperaba 0 canales activos para el usuario tras cerrar, quedan %d", quedan)
	}
}
