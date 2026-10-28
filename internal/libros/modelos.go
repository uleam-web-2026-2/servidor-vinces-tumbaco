package libros

type Libro struct {
	ID        uint
	Titulo    string
	Autor     string
	Estado    string
	Prestamos []Prestamo
}

type Prestamo struct {
	ID      uint
	LibroID uint
	Usuario string
	Estado  string
	Libro   Libro
}

var EstadosValidos = map[string]bool{
	"disponible": true,
	"prestado":   true,
	"reservado":  true,
}
