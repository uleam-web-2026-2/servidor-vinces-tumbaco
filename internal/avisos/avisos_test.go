package avisos

import "testing"

func TestEnviarSoloAlUsuarioSuscrito(t *testing.T) {
	a := Nuevo()

	canalUsuario := a.Suscribir("kenyis_admin")
	canalOtro := a.Suscribir("majo_user")

	aviso := Aviso{
		ID:     1,
		Estado: "aceptado",
	}

	a.Enviar("kenyis_admin", aviso)

	select {
	case recibido := <-canalUsuario:
		if recibido.ID != aviso.ID || recibido.Estado != aviso.Estado {
			t.Fatal("el aviso recibido no coincide")
		}
	default:
		t.Fatal("el usuario no recibió el aviso")
	}

	select {
	case <-canalOtro:
		t.Fatal("otro usuario recibió un aviso que no le corresponde")
	default:
	}

	a.Quitar("kenyis_admin", canalUsuario)
	a.Quitar("majo_user", canalOtro)
}

func TestQuitarConexion(t *testing.T) {
	a := Nuevo()

	canal := a.Suscribir("kenyis_admin")
	a.Quitar("kenyis_admin", canal)

	a.Enviar("kenyis_admin", Aviso{
		ID:     2,
		Estado: "cancelado",
	})

	if len(a.canales["kenyis_admin"]) != 0 {
		t.Fatal("la conexión no fue eliminada")
	}
}
