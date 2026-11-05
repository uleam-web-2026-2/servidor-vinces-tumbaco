package intercambios

import "time"

// USUARIO
type Usuario struct {
	ID            uint      `gorm:"column:id;primaryKey" json:"id_usuario"`
	Nombre        string    `gorm:"size:50;not null" json:"nombre"`
	Apellido      string    `gorm:"size:50" json:"apellido"`
	Correo        string    `gorm:"size:100;unique;not null" json:"correo"`
	Contrasena    string    `gorm:"column:contrasena;size:255" json:"-"`
	Telefono      string    `gorm:"size:20" json:"telefono"`
	Direccion     string    `gorm:"size:150" json:"direccion"`
	Ciudad        string    `gorm:"size:50" json:"ciudad"`
	FechaRegistro time.Time `json:"fecha_registro"`
	Estado        string    `gorm:"size:20" json:"estado"`
	Rol           string    `gorm:"size:20;not null" json:"rol"`
}

func (Usuario) TableName() string {
	return "usuarios"
}

// LIBRO
type Libro struct {
	ID              uint   `gorm:"column:id;primaryKey" json:"id_libro"`
	Titulo          string `gorm:"size:150;not null" json:"titulo"`
	Autor           string `gorm:"type:text" json:"autor,omitempty"`
	Genero          string `gorm:"type:text" json:"genero,omitempty"`
	Estado          string `gorm:"type:text" json:"estado,omitempty"`
	ISBN            string `gorm:"size:20;uniqueIndex" json:"isbn"`
	Descripcion     string `gorm:"type:text" json:"descripcion"`
	AnioPublicacion int    `gorm:"column:anio_publicacion" json:"anio_publicacion"`
	Editorial       string `gorm:"size:100" json:"editorial"`
	Idioma          string `gorm:"size:50" json:"idioma"`
	Paginas         int    `json:"paginas"`
	Disponible      bool   `gorm:"not null;default:true" json:"disponible"`
	UsuarioID       *uint  `gorm:"column:usuario_id" json:"usuario_id,omitempty"`
}

func (Libro) TableName() string {
	return "libros"
}

// AUTOR
type Autor struct {
	ID              uint       `gorm:"column:id;primaryKey" json:"id_autor"`
	Nombre          string     `gorm:"size:50" json:"nombre"`
	Apellido        string     `gorm:"size:50" json:"apellido"`
	Nacionalidad    string     `gorm:"size:50" json:"nacionalidad"`
	FechaNacimiento *time.Time `json:"fecha_nacimiento"`
}

func (Autor) TableName() string {
	return "autors"
}

// CATEGORIA
type Categoria struct {
	ID          uint   `gorm:"column:id;primaryKey" json:"id_categoria"`
	Nombre      string `gorm:"size:50" json:"nombre"`
	Descripcion string `gorm:"size:200" json:"descripcion"`
}

func (Categoria) TableName() string {
	return "categoria"
}

// LIBRO_AUTOR
type LibroAutor struct {
	LibroID uint   `gorm:"column:libro_id;primaryKey" json:"libro_id"`
	AutorID uint   `gorm:"column:autor_id;primaryKey" json:"autor_id"`
	Rol     string `gorm:"size:30" json:"rol"`

	Libro Libro `gorm:"foreignKey:LibroID;references:ID" json:"-"`
	Autor Autor `gorm:"foreignKey:AutorID;references:ID" json:"-"`
}

func (LibroAutor) TableName() string {
	return "libro_autors"
}

// LIBRO_CATEGORIA
type LibroCategoria struct {
	LibroID     uint `gorm:"column:libro_id;primaryKey" json:"libro_id"`
	CategoriaID uint `gorm:"column:categoria_id;primaryKey" json:"categoria_id"`

	Libro     Libro     `gorm:"foreignKey:LibroID;references:ID" json:"-"`
	Categoria Categoria `gorm:"foreignKey:CategoriaID;references:ID" json:"-"`
}

func (LibroCategoria) TableName() string {
	return "libro_categoria"
}

// EJEMPLAR_LIBRO
type EjemplarLibro struct {
	ID             uint      `gorm:"column:id;primaryKey" json:"id_ejemplar"`
	LibroID        uint      `gorm:"column:libro_id;not null" json:"libro_id"`
	UsuarioID      uint      `gorm:"column:usuario_id;not null" json:"usuario_id"`
	EstadoFisico   string    `gorm:"size:30" json:"estado_fisico"`
	Disponibilidad string    `gorm:"size:30" json:"disponibilidad"`
	FechaRegistro  time.Time `json:"fecha_registro"`
	Observaciones  string    `gorm:"type:text" json:"observaciones"`

	Libro   Libro   `gorm:"foreignKey:LibroID;references:ID" json:"-"`
	Usuario Usuario `gorm:"foreignKey:UsuarioID;references:ID" json:"-"`
}

func (EjemplarLibro) TableName() string {
	return "ejemplar_libros"
}

// SOLICITUD_INTERCAMBIO
type SolicitudIntercambio struct {
	ID                   uint      `gorm:"column:id;primaryKey" json:"id_solicitud"`
	UsuarioSolicitanteID uint      `gorm:"column:usuario_solicitante_id" json:"usuario_solicitante_id"`
	EjemplarOfrecidoID   uint      `gorm:"column:ejemplar_ofrecido_id" json:"ejemplar_ofrecido_id"`
	EjemplarDeseadoID    uint      `gorm:"column:ejemplar_deseado_id" json:"ejemplar_deseado_id"`
	FechaSolicitud       time.Time `json:"fecha_solicitud"`
	Estado               string    `gorm:"size:30" json:"estado"`
	Mensaje              string    `gorm:"type:text" json:"mensaje"`

	UsuarioSolicitante Usuario       `gorm:"foreignKey:UsuarioSolicitanteID;references:ID" json:"-"`
	EjemplarOfrecido   EjemplarLibro `gorm:"foreignKey:EjemplarOfrecidoID;references:ID" json:"-"`
	EjemplarDeseado    EjemplarLibro `gorm:"foreignKey:EjemplarDeseadoID;references:ID" json:"-"`
}

func (SolicitudIntercambio) TableName() string {
	return "solicitud_intercambios"
}

// INTERCAMBIO
type Intercambio struct {
	ID                uint       `gorm:"column:id;primaryKey" json:"id_intercambio"`
	LibroID           *uint      `gorm:"column:libro_id" json:"libro_id,omitempty"`
	Persona           string     `gorm:"column:persona;type:text" json:"persona,omitempty"`
	Estado            string     `gorm:"size:30;not null;default:pendiente" json:"estado"`
	UsuarioOfreceID   *uint      `gorm:"column:usuario_ofrece_id" json:"usuario_ofrece_id,omitempty"`
	UsuarioRecibeID   *uint      `gorm:"column:usuario_recibe_id" json:"usuario_recibe_id,omitempty"`
	FechaSolicitud    *time.Time `json:"fecha_solicitud"`
	FechaConfirmacion *time.Time `json:"fecha_confirmacion"`
	FechaFinalizacion *time.Time `json:"fecha_finalizacion"`
	Observaciones     string     `gorm:"type:text" json:"observaciones"`
	SolicitanteID     *uint      `gorm:"column:solicitante_id" json:"solicitante_id,omitempty"`
	Creado            *time.Time `gorm:"column:creado" json:"creado,omitempty"`

	Libro         Libro    `gorm:"foreignKey:LibroID;references:ID" json:"-"`
	UsuarioOfrece *Usuario `gorm:"foreignKey:UsuarioOfreceID;references:ID" json:"-"`
	UsuarioRecibe *Usuario `gorm:"foreignKey:UsuarioRecibeID;references:ID" json:"-"`
}

func (Intercambio) TableName() string {
	return "intercambios"
}

// DETALLE_INTERCAMBIO
type DetalleIntercambio struct {
	ID                  uint `gorm:"column:id;primaryKey" json:"id_detalle"`
	IntercambioID       uint `gorm:"column:intercambio_id" json:"intercambio_id"`
	EjemplarEntregadoID uint `gorm:"column:ejemplar_entregado_id" json:"ejemplar_entregado_id"`
	EjemplarRecibidoID  uint `gorm:"column:ejemplar_recibido_id" json:"ejemplar_recibido_id"`
	Cantidad            int  `json:"cantidad"`

	Intercambio       Intercambio   `gorm:"foreignKey:IntercambioID;references:ID" json:"-"`
	EjemplarEntregado EjemplarLibro `gorm:"foreignKey:EjemplarEntregadoID;references:ID" json:"-"`
	EjemplarRecibido  EjemplarLibro `gorm:"foreignKey:EjemplarRecibidoID;references:ID" json:"-"`
}

func (DetalleIntercambio) TableName() string {
	return "detalle_intercambios"
}

// ENTREGA
type Entrega struct {
	ID              uint       `gorm:"column:id;primaryKey" json:"id_entrega"`
	IntercambioID   uint       `gorm:"column:intercambio_id;uniqueIndex" json:"intercambio_id"`
	LugarEntrega    string     `gorm:"size:150" json:"lugar_entrega"`
	FechaProgramada *time.Time `json:"fecha_programada"`
	FechaReal       *time.Time `json:"fecha_real"`
	Estado          string     `gorm:"size:30" json:"estado"`

	Intercambio Intercambio `gorm:"foreignKey:IntercambioID;references:ID" json:"-"`
}

func (Entrega) TableName() string {
	return "entregas"
}

// RESENA
type Resena struct {
	ID                uint      `gorm:"column:id;primaryKey" json:"id_resena"`
	IntercambioID     uint      `gorm:"column:intercambio_id" json:"intercambio_id"`
	UsuarioAutorID    uint      `gorm:"column:usuario_autor_id" json:"usuario_autor_id"`
	UsuarioEvaluadoID uint      `gorm:"column:usuario_evaluado_id" json:"usuario_evaluado_id"`
	Puntuacion        int       `json:"puntuacion"`
	Comentario        string    `gorm:"type:text" json:"comentario"`
	Fecha             time.Time `json:"fecha"`

	Intercambio     Intercambio `gorm:"foreignKey:IntercambioID;references:ID" json:"-"`
	UsuarioAutor    Usuario     `gorm:"foreignKey:UsuarioAutorID;references:ID" json:"-"`
	UsuarioEvaluado Usuario     `gorm:"foreignKey:UsuarioEvaluadoID;references:ID" json:"-"`
}

func (Resena) TableName() string {
	return "resenas"
}

// Estados del CRUD de intercambios.
const (
	EstadoPendiente = "pendiente"
	EstadoAceptado  = "aceptado"
	EstadoEntregado = "entregado"
	EstadoCancelado = "cancelado"
)

func EstadoValido(estado string) bool {
	switch estado {
	case EstadoPendiente, EstadoAceptado, EstadoEntregado, EstadoCancelado:
		return true
	default:
		return false
	}
}
