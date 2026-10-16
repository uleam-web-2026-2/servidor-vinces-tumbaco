package intercambios

import "gorm.io/gorm"

// Sembrar inserta los datos iniciales del proyecto.
func Sembrar(db *gorm.DB) error {

	return db.Transaction(func(tx *gorm.DB) error {

		// Comprobar si ya existen usuarios.
		var cantidad int64

		if err := tx.Model(&Usuario{}).Count(&cantidad).Error; err != nil {
			return err
		}

		if cantidad > 0 {
			return nil
		}

		// Usuarios iniciales.
		usuarios := []Usuario{
			{
				Nombre: "Ana",
				Correo: "ana@email.com",
				Rol:    "usuario",
			},
			{
				Nombre: "Carlos",
				Correo: "carlos@email.com",
				Rol:    "usuario",
			},
			{
				Nombre: "Administrador",
				Correo: "admin@email.com",
				Rol:    "administrador",
			},
		}

		if err := tx.Create(&usuarios).Error; err != nil {
			return err
		}

		// Libros iniciales con ISBN diferentes.
		libros := []Libro{
			{
				ISBN:       "9780307474728",
				Titulo:     "Cien años de soledad",
				Autor:      "Gabriel García Márquez",
				Genero:     "Novela",
				Disponible: true,
				UsuarioID:  usuarios[0].ID,
			},
			{
				ISBN:       "9780156012195",
				Titulo:     "El Principito",
				Autor:      "Antoine de Saint-Exupéry",
				Genero:     "Fantasía",
				Disponible: true,
				UsuarioID:  usuarios[1].ID,
			},
			{
				ISBN:       "9780060934347",
				Titulo:     "Don Quijote de la Mancha",
				Autor:      "Miguel de Cervantes",
				Genero:     "Novela",
				Disponible: true,
				UsuarioID:  usuarios[0].ID,
			},
		}

		if err := tx.Create(&libros).Error; err != nil {
			return err
		}

		// Intercambios iniciales.
		intercambiosIniciales := []Intercambio{
			{
				LibroID:       libros[0].ID,
				SolicitanteID: usuarios[1].ID,
				Estado:        EstadoPendiente,
			},
			{
				LibroID:       libros[1].ID,
				SolicitanteID: usuarios[0].ID,
				Estado:        EstadoAceptado,
			},
		}

		if err := tx.Create(&intercambiosIniciales).Error; err != nil {
			return err
		}

		return nil
	})
}
