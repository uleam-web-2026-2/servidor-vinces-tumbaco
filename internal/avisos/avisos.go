package avisos

import "sync"

// Aviso representa un cambio de estado que debe notificarse.
type Aviso struct {
	ID     uint   `json:"id"`
	Estado string `json:"estado"`
}

// Avisos administra los canales de notificación por usuario.
type Avisos struct {
	mu      sync.Mutex
	canales map[string][]chan Aviso
}

// Nuevo crea un administrador de avisos.
func Nuevo() *Avisos {
	return &Avisos{
		canales: make(map[string][]chan Aviso),
	}
}

// Suscribir registra una conexión para un usuario.
func (a *Avisos) Suscribir(usuario string) chan Aviso {
	canal := make(chan Aviso, 10)

	a.mu.Lock()
	a.canales[usuario] = append(a.canales[usuario], canal)
	a.mu.Unlock()

	return canal
}

// Quitar elimina una conexión de un usuario.
func (a *Avisos) Quitar(usuario string, canal chan Aviso) {
	a.mu.Lock()
	defer a.mu.Unlock()

	conexiones := a.canales[usuario]

	for i, existente := range conexiones {
		if existente == canal {
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

// Enviar envía un aviso a todas las conexiones del usuario.
func (a *Avisos) Enviar(usuario string, aviso Aviso) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, canal := range a.canales[usuario] {
		select {
		case canal <- aviso:
		default:
			// No bloquear si una conexión está lenta.
		}
	}
}
