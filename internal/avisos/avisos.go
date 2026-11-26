package avisos

import "sync"

// Aviso representa una notificacion de dominio de intercambio.
type Aviso struct {
	ID     uint   `json:"id"`
	Estado string `json:"estado"`
	Evento string `json:"evento,omitempty"`
}

type suscripcion struct {
	canal chan Aviso
	rol   string
}

// Avisos administra los canales de notificacion por usuario y rol.
type Avisos struct {
	mu      sync.Mutex
	canales map[string][]suscripcion
}

// Nuevo crea un administrador de avisos.
func Nuevo() *Avisos {
	return &Avisos{
		canales: make(map[string][]suscripcion),
	}
}

// Suscribir registra una conexion para un usuario sin rol especifico.
func (a *Avisos) Suscribir(usuario string) chan Aviso {
	return a.SuscribirConRol(usuario, "")
}

// SuscribirConRol registra una conexion para un usuario asociada a un rol.
func (a *Avisos) SuscribirConRol(usuario string, rol string) chan Aviso {
	canal := make(chan Aviso, 10)

	a.mu.Lock()
	a.canales[usuario] = append(a.canales[usuario], suscripcion{
		canal: canal,
		rol:   rol,
	})
	a.mu.Unlock()

	return canal
}

// Quitar elimina una conexion de un usuario.
func (a *Avisos) Quitar(usuario string, canal chan Aviso) {
	a.mu.Lock()
	defer a.mu.Unlock()

	conexiones := a.canales[usuario]

	for i, existente := range conexiones {
		if existente.canal == canal {
			a.canales[usuario] = append(
				conexiones[:i],
				conexiones[i+1:]...,
			)
			close(canal)
			break
		}
	}

	if len(a.canales[usuario]) == 0 {
		delete(a.canales, usuario)
	}
}

// Enviar envia un aviso a todas las conexiones del usuario.
func (a *Avisos) Enviar(usuario string, aviso Aviso) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, sub := range a.canales[usuario] {
		select {
		case sub.canal <- aviso:
		default:
			// No bloquear si una conexion esta lenta.
		}
	}
}

// EnviarPorRol envia un aviso solo a las conexiones activas con el rol indicado.
func (a *Avisos) EnviarPorRol(rol string, aviso Aviso) {
	if rol == "" {
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	for _, subs := range a.canales {
		for _, sub := range subs {
			if sub.rol != rol {
				continue
			}

			select {
			case sub.canal <- aviso:
			default:
				// No bloquear si una conexion esta lenta.
			}
		}
	}
}
