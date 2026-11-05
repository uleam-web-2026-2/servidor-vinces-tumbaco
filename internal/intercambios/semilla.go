package intercambios

import (
	"time"

	"gorm.io/gorm"
)

// Sembrar inserta los datos iniciales del proyecto.
func Sembrar(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var cantidad int64

		if err := tx.Model(&Usuario{}).Count(&cantidad).Error; err != nil {
			return err
		}

		// Evitar duplicar los datos iniciales.
		if cantidad > 0 {
			return nil
		}

		ahora := time.Now()

		// Usuarios iniciales.
		usuarios := []Usuario{
			{
				Nombre:        "Ana",
				Apellido:      "Lopez",
				Correo:        "ana@email.com",
				FechaRegistro: ahora,
				Estado:        "activo",
				Rol:           "usuario",
			},
			{
				Nombre:        "Carlos",
				Apellido:      "Perez",
				Correo:        "carlos@email.com",
				FechaRegistro: ahora,
				Estado:        "activo",
				Rol:           "usuario",
			},
			{
				Nombre:        "Administrador",
				Correo:        "admin@email.com",
				FechaRegistro: ahora,
				Estado:        "activo",
				Rol:           "administrador",
			},
		}

		if err := tx.Create(&usuarios).Error; err != nil {
			return err
		}

		// Autores iniciales.
		autores := []Autor{
			{
				Nombre:       "Gabriel",
				Apellido:     "Garcia Marquez",
				Nacionalidad: "Colombiana",
			},
			{
				Nombre:       "Antoine",
				Apellido:     "de Saint-Exupery",
				Nacionalidad: "Francesa",
			},
			{
				Nombre:       "Miguel",
				Apellido:     "de Cervantes",
				Nacionalidad: "Espanola",
			},
		}

		if err := tx.Create(&autores).Error; err != nil {
			return err
		}

		// Categorias iniciales.
		categorias := []Categoria{
			{
				Nombre:      "Novela",
				Descripcion: "Obras literarias narrativas",
			},
			{
				Nombre:      "Literatura infantil",
				Descripcion: "Libros dirigidos al publico infantil",
			},
		}

		if err := tx.Create(&categorias).Error; err != nil {
			return err
		}

		// Libros iniciales.
		libros := []Libro{
			{
				ISBN:            "9780307474728",
				Titulo:          "Cien anos de soledad",
				Descripcion:     "Novela de Gabriel Garcia Marquez",
				AnioPublicacion: 1967,
				Editorial:       "Sudamericana",
				Idioma:          "Espanol",
				Paginas:         471,
				Disponible:      true,
			},
			{
				ISBN:            "9780156012195",
				Titulo:          "El Principito",
				Descripcion:     "Obra literaria de Antoine de Saint-Exupery",
				AnioPublicacion: 1943,
				Editorial:       "Reynal and Hitchcock",
				Idioma:          "Espanol",
				Paginas:         96,
				Disponible:      true,
			},
			{
				ISBN:            "9780060934347",
				Titulo:          "Don Quijote de la Mancha",
				Descripcion:     "Obra literaria de Miguel de Cervantes",
				AnioPublicacion: 1605,
				Editorial:       "Edicion clasica",
				Idioma:          "Espanol",
				Paginas:         863,
				Disponible:      true,
			},
		}

		if err := tx.Create(&libros).Error; err != nil {
			return err
		}

		// Relacionar libros con autores.
		librosAutores := []LibroAutor{
			{LibroID: libros[0].ID, AutorID: autores[0].ID, Rol: "autor"},
			{LibroID: libros[1].ID, AutorID: autores[1].ID, Rol: "autor"},
			{LibroID: libros[2].ID, AutorID: autores[2].ID, Rol: "autor"},
		}

		if err := tx.Create(&librosAutores).Error; err != nil {
			return err
		}

		// Relacionar libros con categorias.
		librosCategorias := []LibroCategoria{
			{LibroID: libros[0].ID, CategoriaID: categorias[0].ID},
			{LibroID: libros[1].ID, CategoriaID: categorias[1].ID},
			{LibroID: libros[2].ID, CategoriaID: categorias[0].ID},
		}

		if err := tx.Create(&librosCategorias).Error; err != nil {
			return err
		}

		// Ejemplares fisicos pertenecientes a los usuarios.
		ejemplares := []EjemplarLibro{
			{
				LibroID:        libros[0].ID,
				UsuarioID:      usuarios[0].ID,
				EstadoFisico:   "bueno",
				Disponibilidad: "disponible",
				FechaRegistro:  ahora,
			},
			{
				LibroID:        libros[1].ID,
				UsuarioID:      usuarios[1].ID,
				EstadoFisico:   "muy bueno",
				Disponibilidad: "disponible",
				FechaRegistro:  ahora,
			},
			{
				LibroID:        libros[2].ID,
				UsuarioID:      usuarios[0].ID,
				EstadoFisico:   "bueno",
				Disponibilidad: "disponible",
				FechaRegistro:  ahora,
			},
		}

		if err := tx.Create(&ejemplares).Error; err != nil {
			return err
		}

		// Solicitud inicial de intercambio.
		solicitud := SolicitudIntercambio{
			UsuarioSolicitanteID: usuarios[1].ID,
			EjemplarOfrecidoID:   ejemplares[1].ID,
			EjemplarDeseadoID:    ejemplares[0].ID,
			FechaSolicitud:       ahora,
			Estado:               EstadoPendiente,
			Mensaje:              "Solicitud inicial de intercambio",
		}

		if err := tx.Create(&solicitud).Error; err != nil {
			return err
		}

		// Intercambio inicial.
		usuarioOfreceID := usuarios[0].ID
		usuarioRecibeID := usuarios[1].ID

		intercambio := Intercambio{
			UsuarioOfreceID: &usuarioOfreceID,
			UsuarioRecibeID: &usuarioRecibeID,
			FechaSolicitud:  &ahora,
			Estado:          EstadoPendiente,
			Observaciones:   "Intercambio de libros de demostracion",
		}

		if err := tx.Create(&intercambio).Error; err != nil {
			return err
		}

		// Detalle del intercambio.
		detalle := DetalleIntercambio{
			IntercambioID:       intercambio.ID,
			EjemplarEntregadoID: ejemplares[0].ID,
			EjemplarRecibidoID:  ejemplares[1].ID,
			Cantidad:            1,
		}

		if err := tx.Create(&detalle).Error; err != nil {
			return err
		}

		return nil
	})
}
