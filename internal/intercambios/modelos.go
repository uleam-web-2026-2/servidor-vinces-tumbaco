package intercambios

import "time"

// USUARIO
type Usuario struct {
	ID            uint   `gorm:"column:id;primaryKey"`
	Nombre        string `gorm:"size:50;not null"`
	Apellido      string `gorm:"size:50"`
	Correo        string `gorm:"size:100;unique;not null"`
	Contrasena    string `gorm:"column:contrasena;size:255"`
	Telefono      string `gorm:"size:20"`
	Direccion     string `gorm:"size:150"`
	Ciudad        string `gorm:"size:50"`
	FechaRegistro time.Time
	Estado        string `gorm:"size:20"`
	Rol           string `gorm:"size:20;not null"`
}

func (Usuario) TableName() string {
	return "usuarios"
}

// LIBRO
type Libro struct {
	ID              uint   `gorm:"column:id;primaryKey"`
	Titulo          string `gorm:"size:150;not null"`
	Autor           string `gorm:"type:text"`
	Genero          string `gorm:"type:text"`
	Estado          string `gorm:"type:text"`
	ISBN            string `gorm:"size:20;uniqueIndex"`
	Descripcion     string `gorm:"type:text"`
	AnioPublicacion int    `gorm:"column:anio_publicacion"`
	Editorial       string `gorm:"size:100"`
	Idioma          string `gorm:"size:50"`
	Paginas         int
	Disponible      bool  `gorm:"not null;default:true"`
	UsuarioID       *uint `gorm:"column:usuario_id"`
}

func (Libro) TableName() string {
	return "libros"
}

// AUTOR
type Autor struct {
	ID              uint   `gorm:"column:id;primaryKey"`
	Nombre          string `gorm:"size:50"`
	Apellido        string `gorm:"size:50"`
	Nacionalidad    string `gorm:"size:50"`
	FechaNacimiento *time.Time
}

func (Autor) TableName() string {
	return "autors"
}

// CATEGORIA
type Categoria struct {
	ID          uint   `gorm:"column:id;primaryKey"`
	Nombre      string `gorm:"size:50"`
	Descripcion string `gorm:"size:200"`
}

func (Categoria) TableName() string {
	return "categoria"
}

// LIBRO_AUTOR
type LibroAutor struct {
	LibroID uint   `gorm:"column:libro_id;primaryKey"`
	AutorID uint   `gorm:"column:autor_id;primaryKey"`
	Rol     string `gorm:"size:30"`

	Libro Libro `gorm:"foreignKey:LibroID;references:ID"`
	Autor Autor `gorm:"foreignKey:AutorID;references:ID"`
}

func (LibroAutor) TableName() string {
	return "libro_autors"
}

// LIBRO_CATEGORIA
type LibroCategoria struct {
	LibroID     uint `gorm:"column:libro_id;primaryKey"`
	CategoriaID uint `gorm:"column:categoria_id;primaryKey"`

	Libro     Libro     `gorm:"foreignKey:LibroID;references:ID"`
	Categoria Categoria `gorm:"foreignKey:CategoriaID;references:ID"`
}

func (LibroCategoria) TableName() string {
	return "libro_categoria"
}

// EJEMPLAR_LIBRO
type EjemplarLibro struct {
	ID             uint   `gorm:"column:id;primaryKey"`
	LibroID        uint   `gorm:"column:libro_id;not null"`
	UsuarioID      uint   `gorm:"column:usuario_id;not null"`
	EstadoFisico   string `gorm:"size:30"`
	Disponibilidad string `gorm:"size:30"`
	FechaRegistro  time.Time
	Observaciones  string `gorm:"type:text"`

	Libro   Libro   `gorm:"foreignKey:LibroID;references:ID"`
	Usuario Usuario `gorm:"foreignKey:UsuarioID;references:ID"`
}

func (EjemplarLibro) TableName() string {
	return "ejemplar_libros"
}

// SOLICITUD_INTERCAMBIO
type SolicitudIntercambio struct {
	ID                   uint `gorm:"column:id;primaryKey"`
	UsuarioSolicitanteID uint `gorm:"column:usuario_solicitante_id"`
	EjemplarOfrecidoID   uint `gorm:"column:ejemplar_ofrecido_id"`
	EjemplarDeseadoID    uint `gorm:"column:ejemplar_deseado_id"`
	FechaSolicitud       time.Time
	Estado               string `gorm:"size:30"`
	Mensaje              string `gorm:"type:text"`

	UsuarioSolicitante Usuario       `gorm:"foreignKey:UsuarioSolicitanteID;references:ID"`
	EjemplarOfrecido   EjemplarLibro `gorm:"foreignKey:EjemplarOfrecidoID;references:ID"`
	EjemplarDeseado    EjemplarLibro `gorm:"foreignKey:EjemplarDeseadoID;references:ID"`
}

func (SolicitudIntercambio) TableName() string {
	return "solicitud_intercambios"
}

// INTERCAMBIO
type Intercambio struct {
	ID                uint   `gorm:"column:id;primaryKey"`
	LibroID           *uint  `gorm:"column:libro_id"`
	Persona           string `gorm:"column:persona;type:text"`
	Estado            string `gorm:"size:30;not null;default:pendiente"`
	UsuarioOfreceID   *uint  `gorm:"column:usuario_ofrece_id"`
	UsuarioRecibeID   *uint  `gorm:"column:usuario_recibe_id"`
	FechaSolicitud    *time.Time
	FechaConfirmacion *time.Time
	FechaFinalizacion *time.Time
	Observaciones     string     `gorm:"type:text"`
	SolicitanteID     *uint      `gorm:"column:solicitante_id"`
	Creado            *time.Time `gorm:"column:creado"`

	Libro         Libro    `gorm:"foreignKey:LibroID;references:ID"`
	UsuarioOfrece *Usuario `gorm:"foreignKey:UsuarioOfreceID;references:ID"`
	UsuarioRecibe *Usuario `gorm:"foreignKey:UsuarioRecibeID;references:ID"`
}

func (Intercambio) TableName() string {
	return "intercambios"
}

// DETALLE_INTERCAMBIO
type DetalleIntercambio struct {
	ID                  uint `gorm:"column:id;primaryKey"`
	IntercambioID       uint `gorm:"column:intercambio_id"`
	EjemplarEntregadoID uint `gorm:"column:ejemplar_entregado_id"`
	EjemplarRecibidoID  uint `gorm:"column:ejemplar_recibido_id"`
	Cantidad            int

	Intercambio       Intercambio   `gorm:"foreignKey:IntercambioID;references:ID"`
	EjemplarEntregado EjemplarLibro `gorm:"foreignKey:EjemplarEntregadoID;references:ID"`
	EjemplarRecibido  EjemplarLibro `gorm:"foreignKey:EjemplarRecibidoID;references:ID"`
}

func (DetalleIntercambio) TableName() string {
	return "detalle_intercambios"
}

// ENTREGA
type Entrega struct {
	ID              uint   `gorm:"column:id;primaryKey"`
	IntercambioID   uint   `gorm:"column:intercambio_id;uniqueIndex"`
	LugarEntrega    string `gorm:"size:150"`
	FechaProgramada *time.Time
	FechaReal       *time.Time
	Estado          string `gorm:"size:30"`

	Intercambio Intercambio `gorm:"foreignKey:IntercambioID;references:ID"`
}

func (Entrega) TableName() string {
	return "entregas"
}

// RESENA
type Resena struct {
	ID                uint `gorm:"column:id;primaryKey"`
	IntercambioID     uint `gorm:"column:intercambio_id"`
	UsuarioAutorID    uint `gorm:"column:usuario_autor_id"`
	UsuarioEvaluadoID uint `gorm:"column:usuario_evaluado_id"`
	Puntuacion        int
	Comentario        string `gorm:"type:text"`
	Fecha             time.Time

	Intercambio     Intercambio `gorm:"foreignKey:IntercambioID;references:ID"`
	UsuarioAutor    Usuario     `gorm:"foreignKey:UsuarioAutorID;references:ID"`
	UsuarioEvaluado Usuario     `gorm:"foreignKey:UsuarioEvaluadoID;references:ID"`
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
