package intercambios

import "net/http"

// RutasProtegidas configura los permisos de las rutas de intercambios.
// GET: administrador y usuario.
// POST, PUT, PATCH y DELETE: solo administrador.
func RutasProtegidas(handler *Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			Proteger(handler.Rutas, "admin", "usuario")(w, r)

		case http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete:
			Proteger(handler.Rutas, "admin")(w, r)

		default:
			handler.Rutas(w, r)
		}
	}
}


