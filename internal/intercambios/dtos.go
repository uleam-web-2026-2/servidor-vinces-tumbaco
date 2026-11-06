package intercambios

import "time"

// DTO para recibir los datos al crear un intercambio.
type CrearIntercambioDTO struct {
	UsuarioSolicitanteID uint   `json:"id_usuario_solicitante"`
	EjemplarOfrecidoID   uint   `json:"id_ejemplar_ofrecido"`
	EjemplarDeseadoID    uint   `json:"id_ejemplar_deseado"`
	Observaciones        string `json:"observaciones"`
}

// DTO para recibir los datos al actualizar un intercambio.
type ActualizarIntercambioDTO struct {
	UsuarioOfreceID uint   `json:"id_usuario_ofrece"`
	UsuarioRecibeID uint   `json:"id_usuario_recibe"`
	Observaciones   string `json:"observaciones"`
}

// DTO para recibir cambios de estado.
type CambiarEstadoDTO struct {
	Estado string `json:"estado"`
}

// DTO para devolver los datos de un intercambio.
type IntercambioRespuestaDTO struct {
	ID                uint       `json:"id_intercambio"`
	LibroID           *uint      `json:"libro_id,omitempty"`
	Persona           string     `json:"persona,omitempty"`
	Estado            string     `json:"estado"`
	UsuarioOfreceID   *uint      `json:"usuario_ofrece_id,omitempty"`
	UsuarioRecibeID   *uint      `json:"usuario_recibe_id,omitempty"`
	FechaSolicitud    *time.Time `json:"fecha_solicitud"`
	FechaConfirmacion *time.Time `json:"fecha_confirmacion"`
	FechaFinalizacion *time.Time `json:"fecha_finalizacion"`
	Observaciones     string     `json:"observaciones"`
	SolicitanteID     *uint      `json:"solicitante_id,omitempty"`
	Creado            *time.Time `json:"creado,omitempty"`
}

// Convertir los datos recibidos al formato del servicio.
func (dto CrearIntercambioDTO) AInput() CrearIntercambioInput {
	return CrearIntercambioInput{
		UsuarioSolicitanteID: dto.UsuarioSolicitanteID,
		EjemplarOfrecidoID:   dto.EjemplarOfrecidoID,
		EjemplarDeseadoID:    dto.EjemplarDeseadoID,
		Observaciones:        dto.Observaciones,
	}
}

// Convertir los datos recibidos al formato del servicio.
func (dto ActualizarIntercambioDTO) AInput() ActualizarIntercambioInput {
	return ActualizarIntercambioInput{
		UsuarioOfreceID: dto.UsuarioOfreceID,
		UsuarioRecibeID: dto.UsuarioRecibeID,
		Observaciones:   dto.Observaciones,
	}
}

// Convertir el modelo de intercambio a DTO de respuesta.
func IntercambioARespuestaDTO(modelo *Intercambio) IntercambioRespuestaDTO {
	return IntercambioRespuestaDTO{
		ID:                modelo.ID,
		LibroID:           modelo.LibroID,
		Persona:           modelo.Persona,
		Estado:            modelo.Estado,
		UsuarioOfreceID:   modelo.UsuarioOfreceID,
		UsuarioRecibeID:   modelo.UsuarioRecibeID,
		FechaSolicitud:    modelo.FechaSolicitud,
		FechaConfirmacion: modelo.FechaConfirmacion,
		FechaFinalizacion: modelo.FechaFinalizacion,
		Observaciones:     modelo.Observaciones,
		SolicitanteID:     modelo.SolicitanteID,
		Creado:            modelo.Creado,
	}
}
