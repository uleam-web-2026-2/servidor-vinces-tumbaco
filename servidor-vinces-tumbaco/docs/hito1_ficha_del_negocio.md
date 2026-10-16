# Hito 1 · Ficha del negocio

**Pareja:** Kenyis Yusley Tumbaco Pillasagua · [NOMBRE DEL COMPAÑERO/A]

**Paralelo:** [ESCRIBIR PARALELO]

**Negocio en una línea:** Plataforma que permite a los usuarios publicar libros usados y solicitar intercambios de libros con otros usuarios.

---

## 1. Negocio de referencia

Como negocio de referencia se tomó a **Clankart**, una plataforma creada para facilitar la compra y venta de libros usados entre estudiantes. La idea surgió de la necesidad de conectar a estudiantes que tenían libros que ya no utilizaban con otros interesados en conseguirlos.

La plataforma permite publicar libros usados y contactar con otros estudiantes interesados. Su público principal son estudiantes que buscan una alternativa para aprovechar libros que ya poseen o conseguir libros de segunda mano. De acuerdo con las cifras declaradas por su fundador en Starter Story, Clankart llegó a contar con más de **25.000 estudiantes registrados** y reportó aproximadamente **1.200 dólares de ingresos mensuales**.

Este caso sirve como referencia para nuestro proyecto porque demuestra que una plataforma especializada puede conectar personas que tienen libros disponibles con otras que los necesitan. Nuestra propuesta adapta esta idea hacia el intercambio directo de libros usados.

**Enlace:** https://www.starterstory.com/stories/clankart

---

## 2. Caso de contraste

Como caso de contraste se tomó a **BIGWORDS**, una empresa relacionada con el mercado de libros de texto. El negocio tuvo un crecimiento rápido y consiguió una cantidad importante de financiamiento, pero posteriormente entró en bancarrota y cerró durante la crisis de las empresas puntocom.

Nuestra hipótesis es que una diferencia importante se encuentra en la **estructura de costos**. Un negocio que compra, almacena o distribuye productos necesita mantener una operación más costosa. Nuestra Plataforma de Intercambio de Libros Usados busca conectar directamente a los usuarios, sin comprar ni mantener un inventario propio de libros.

Esto permite plantear un modelo más sencillo, donde los propios usuarios publican los libros que tienen disponibles y gestionan sus solicitudes de intercambio.

**Fuente:** https://techcrunch.com/2012/08/16/the-bigwords-story/

---

## 3. Adaptación al Ecuador

Para adaptar la plataforma al contexto ecuatoriano se identificaron las siguientes restricciones:

1. **Entregas entre usuarios:** los usuarios pueden encontrarse en diferentes sectores o ciudades. Por esta razón, la plataforma debe registrar información que permita coordinar la entrega del libro y conocer cuándo el intercambio realmente fue completado.

2. **Acceso y confianza entre desconocidos:** un usuario puede solicitar un libro y después no continuar con el intercambio. Por esto, la solicitud necesita diferentes estados que permitan saber si está pendiente, aceptada, entregada o cancelada.

3. **Costo del intercambio:** enviar un libro a otra ciudad puede tener un costo que haga poco conveniente el intercambio. Por ello, inicialmente la plataforma estará orientada principalmente a facilitar intercambios entre usuarios que puedan coordinar directamente la entrega.

**Cambio realizado en el modelo:**  
La segunda restricción cambió directamente nuestro modelo. Se agregó el atributo **estado** a la entidad `Intercambio`, con los valores `pendiente`, `aceptado`, `entregado` y `cancelado`. Esto permite controlar el proceso y evitar que una solicitud sea considerada completada antes de realizar la entrega.

---

## 4. Modelo de datos

### Entidad: Usuario

| Atributo | Tipo | Obligatorio | Ejemplo |
|---|---|---|---|
| id | número entero | sí | 1 |
| nombre | texto | sí | Ana López |
| correo | texto | sí | ana@email.com |
| rol | uno de: usuario, administrador | sí | usuario |

### Entidad: Libro

| Atributo | Tipo | Obligatorio | Ejemplo |
|---|---|---|---|
| id | número entero | sí | 1 |
| titulo | texto | sí | El Principito |
| autor | texto | sí | Antoine de Saint-Exupéry |
| genero | texto | sí | Literatura |
| disponible | sí/no | sí | sí |
| usuario_id | referencia a Usuario | sí | 1 |

### Entidad: Intercambio

| Atributo | Tipo | Obligatorio | Ejemplo |
|---|---|---|---|
| id | número entero | sí | 1 |
| libro_id | referencia a Libro | sí | 1 |
| solicitante_id | referencia a Usuario | sí | 2 |
| estado | uno de: pendiente, aceptado, entregado, cancelado | sí | pendiente |
| creado | fecha y hora | sí | 30-09-2026 14:30 |

### Relaciones

| Entidades | Cardinalidad | Descripción |
|---|---|---|
| Usuario - Libro | uno a muchos | Un usuario puede publicar muchos libros; cada libro pertenece a un usuario. |
| Libro - Intercambio | uno a muchos | Un libro puede recibir varias solicitudes de intercambio; cada intercambio corresponde a un libro. |
| Usuario - Intercambio | uno a muchos | Un usuario puede realizar varias solicitudes de intercambio; cada solicitud pertenece a un usuario solicitante. |

### Structs en Go

```go
type Usuario struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Nombre string `gorm:"not null" json:"nombre"`
	Correo string `gorm:"unique;not null" json:"correo"`
	Rol    string `gorm:"not null" json:"rol"`
}

type Libro struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Titulo     string `gorm:"not null" json:"titulo"`
	Autor      string `gorm:"not null" json:"autor"`
	Genero     string `gorm:"not null" json:"genero"`
	Disponible bool   `gorm:"not null" json:"disponible"`
	UsuarioID  uint   `gorm:"not null" json:"usuario_id"`
	Usuario    Usuario `gorm:"foreignKey:UsuarioID" json:"usuario,omitempty"`
}

type Intercambio struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	LibroID       uint      `gorm:"not null" json:"libro_id"`
	SolicitanteID uint      `gorm:"not null" json:"solicitante_id"`
	Estado        string    `gorm:"not null" json:"estado"`
	Creado        time.Time `gorm:"autoCreateTime" json:"creado"`
}