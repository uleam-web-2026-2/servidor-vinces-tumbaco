# Hito 1 — Addendum Técnico

## A. Estructura del repositorio

La estructura principal del proyecto es:

```text
servidor-vinces-tumbaco-main/
│
├── docs/
│   ├── decisiones.md
│   ├── hito1_addendum.md
│   └── hito1_ficha_del_negocio.md
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   └── libros/
│       ├── manejadores.go
│       └── manejadores_test.go
│
├── .env
├── go.mod
├── go.sum
├── main.go
└── README.md
```

### Descripción

- `main.go`: inicia el servidor, configura las rutas y la conexión con la aplicación.
- `internal/config/`: contiene la configuración del proyecto y lectura de variables de entorno.
- `internal/libros/`: contiene la lógica, manejadores y pruebas relacionadas con los libros.
- `docs/`: contiene la documentación correspondiente al proyecto.
- `go.mod` y `go.sum`: administran las dependencias utilizadas por Go.
- `.env`: contiene las variables de entorno utilizadas de forma local.

---

## B. Variables de entorno

El proyecto utiliza variables de entorno para configurar el servidor y la conexión con PostgreSQL.

La variable principal utilizada es:

| Variable | Propósito |
|---|---|
| `DATABASE_URL` | Define la conexión con la base de datos PostgreSQL. |
| `PUERTO` | Permite definir el puerto utilizado por el servidor. |
| `TIEMPO_ESPERA_SEGUNDOS` | Define el tiempo de espera utilizado por el servidor. |

Ejemplo de configuración:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/intercambio_libros?sslmode=disable
PUERTO=8080
TIEMPO_ESPERA_SEGUNDOS=5
```

Las credenciales reales no deben publicarse en el repositorio. El archivo `.env` se utiliza únicamente para la configuración local.

---

## C. Pruebas

Las pruebas automáticas se ejecutan mediante:

```powershell
go test ./...
```

El resultado obtenido fue:

```text
?    github.com/uleam-web-2026-2/KEMA                  [no test files]
?    github.com/uleam-web-2026-2/KEMA/internal/config  [no test files]
ok   github.com/uleam-web-2026-2/KEMA/internal/libros
```

También se ejecutaron las pruebas detalladas mediante:

```powershell
go test ./internal/libros -v
```

Se comprobaron los siguientes casos:

1. JSON incorrecto responde con código `400`.
2. Título vacío responde con código `422`.
3. Estado inválido responde con código `422`.
4. Autor vacío responde con código `422`.
5. Título y autor iguales responden con código `422`.
6. ID inválido es rechazado correctamente.

El resultado final de las pruebas fue:

```text
PASS
ok github.com/uleam-web-2026-2/KEMA/internal/libros
```

**Evidencia:** insertar aquí una captura de pantalla donde se observe la ejecución de:

```powershell
go test ./internal/libros -v
```

---

## D. Boceto de la pantalla principal

La pantalla principal propuesta permitirá consultar los libros registrados en la plataforma.

```text
+--------------------------------------------------+
|       PLATAFORMA DE INTERCAMBIO DE LIBROS        |
+--------------------------------------------------+
| Buscar libro: [________________________] [Buscar] |
+--------------------------------------------------+
| ID | Título                 | Autor     | Estado  |
|----|------------------------|-----------|---------|
| 1  | Cien años de soledad   | G. Márquez|Disponible|
| 2  | El principito          | Saint-Ex. |Disponible|
| 3  | Don Quijote            | Cervantes |Disponible|
+--------------------------------------------------+
|               [ Publicar libro ]                 |
+--------------------------------------------------+
```

La pantalla consumiría principalmente el endpoint:

```text
GET /libros
```

Los campos principales visibles serían:

- ID del libro.
- Título.
- Autor.
- Estado.

El estado permite conocer la situación actual del libro dentro de la plataforma.

---

## E. Diagrama de secuencia

El siguiente ejemplo representa la consulta de los libros disponibles.

```mermaid
sequenceDiagram
    actor Usuario
    participant Pantalla
    participant API
    participant BD as PostgreSQL

    Usuario->>Pantalla: Ingresa a la lista de libros
    Pantalla->>API: GET /libros
    API->>BD: Consultar libros
    BD-->>API: Lista de libros
    API-->>Pantalla: HTTP 200 + JSON
    Pantalla-->>Usuario: Muestra los libros
```

El flujo comienza cuando el usuario consulta la pantalla de libros. La pantalla realiza una petición `GET /libros`, el servidor consulta PostgreSQL y devuelve la información en formato JSON.

---

## F. Evidencias de respuestas del servidor

### Respuesta correcta

Se realizó una petición:

```text
POST /libros
```

con los siguientes datos:

```json
{
  "titulo": "El Alquimista",
  "autor": "Paulo Coelho",
  "estado": "disponible"
}
```

El servidor respondió correctamente con:

```text
StatusCode: 201
StatusDescription: Created
```

Respuesta:

```json
{
  "ID": 5,
  "Titulo": "El Alquimista",
  "Autor": "Paulo Coelho",
  "Estado": "disponible",
  "Prestamos": null
}
```

**Evidencia:** insertar aquí la captura de pantalla de PowerShell donde se observe el código `201` y la respuesta del servidor.

### Respuesta con error de validación

También se realizó una petición `POST /libros` enviando el título vacío:

```json
{
  "titulo": "",
  "autor": "Gabriel Garcia Marquez",
  "estado": "disponible"
}
```

El servidor rechazó correctamente la petición porque el título es obligatorio.

Respuesta obtenida:

```text
Titulo y autor son obligatorios
```

Este caso corresponde al código HTTP:

```text
422 Unprocessable Entity
```

**Evidencia:** insertar aquí la captura de pantalla donde se observe la respuesta de validación.

---

Con estas evidencias se comprueba el funcionamiento del servidor, la conexión con PostgreSQL, las validaciones implementadas y las pruebas automáticas del proyecto.