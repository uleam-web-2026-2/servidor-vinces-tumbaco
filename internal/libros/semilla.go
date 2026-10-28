package libros

import "gorm.io/gorm"

func Sembrar(db *gorm.DB) error {
	libro1 := Libro{
		Titulo: "Cien años de soledad",
		Autor:  "Gabriel García Márquez",
		Estado: "disponible",
	}

	libro2 := Libro{
		Titulo: "Don Quijote de la Mancha",
		Autor:  "Miguel de Cervantes",
		Estado: "prestado",
	}

	libro3 := Libro{
		Titulo: "El principito",
		Autor:  "Antoine de Saint-Exupéry",
		Estado: "reservado",
	}

	if err := db.Create(&libro1).Error; err != nil {
		return err
	}

	if err := db.Create(&libro2).Error; err != nil {
		return err
	}

	if err := db.Create(&libro3).Error; err != nil {
		return err
	}

	prestamos := []Prestamo{
		{
			LibroID: libro2.ID,
			Usuario: "María",
			Estado:  "activo",
		},
		{
			LibroID: libro3.ID,
			Usuario: "Carlos",
			Estado:  "pendiente",
		},
	}

	return db.Create(&prestamos).Error
}
