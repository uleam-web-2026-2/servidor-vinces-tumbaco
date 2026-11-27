package avisos

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/uleam-web-2026-2/KEMA/internal/intercambios"
)

type notificadorAdministracion struct {
	hub *Avisos
	rol string
}

func (n *notificadorAdministracion) NotificarCreacion(intercambio *intercambios.Intercambio, creadorID uint) {
	if n == nil || n.hub == nil || intercambio == nil {
		return
	}

	n.hub.EnviarPorRol(
		n.rol,
		Aviso{
			ID:     intercambio.ID,
			Estado: intercambio.Estado,
			Evento: "intercambio_creado",
		},
		strconv.FormatUint(uint64(creadorID), 10),
	)
}

func TestRetoAvisaALosQueAdministran(t *testing.T) {
	hub := Nuevo()

	admin1 := hub.SuscribirConRol("10", "administrador")
	admin2 := hub.SuscribirConRol("20", "administrador")
	creadorAdmin := hub.SuscribirConRol("3", "administrador")
	otroRol := hub.SuscribirConRol("99", "usuario")

	defer hub.Quitar("10", admin1)
	defer hub.Quitar("20", admin2)
	defer hub.Quitar("3", creadorAdmin)
	defer hub.Quitar("99", otroRol)

	notif := &notificadorAdministracion{hub: hub, rol: "administrador"}
	notif.NotificarCreacion(&intercambios.Intercambio{
		ID:     42,
		Estado: intercambios.EstadoPendiente,
	}, 3)

	for nombre, ch := range map[string]chan Aviso{
		"admin1": admin1,
		"admin2": admin2,
	} {
		select {
		case msg := <-ch:
			raw, err := json.Marshal(msg)
			if err != nil {
				t.Fatalf("error serializando aviso para %s: %v", nombre, err)
			}
			if msg.ID != 42 || msg.Evento != "intercambio_creado" {
				t.Fatalf("aviso inesperado para %s: %s", nombre, string(raw))
			}
		case <-time.After(200 * time.Millisecond):
			t.Fatalf("no llego aviso a %s", nombre)
		}
	}

	for nombre, ch := range map[string]chan Aviso{
		"creadorAdmin": creadorAdmin,
		"otroRol":      otroRol,
	} {
		select {
		case msg := <-ch:
			t.Fatalf("no debia llegar aviso a %s, llego: %+v", nombre, msg)
		case <-time.After(100 * time.Millisecond):
		}
	}
}
