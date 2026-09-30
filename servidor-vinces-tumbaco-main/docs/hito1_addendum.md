# Hito 1 — Addendum

## 1. Alcance del proyecto

El proyecto consiste en una **Plataforma de Intercambio de Libros Usados** orientada principalmente a estudiantes.

La primera versión permitirá:

* Registrar usuarios.
* Publicar libros usados.
* Consultar libros disponibles.
* Solicitar intercambios.
* Aceptar o rechazar solicitudes.
* Registrar el estado del intercambio.

## 2. Fuera del alcance

En esta primera versión no se incluye:

* Pagos en línea.
* Envíos a domicilio.
* Integración con bancos.
* Sistema de calificaciones.
* Chat en tiempo real.
* Aplicación móvil.

Estas funciones podrían considerarse en futuras versiones.

## 3. Entidades principales

El modelo utilizará tres entidades principales:

### Usuario

Representa a la persona que utiliza la plataforma.

### Libro

Representa el libro usado que un usuario publica para intercambiar.

### Intercambio

Representa la solicitud de intercambio entre dos usuarios.

## 4. Estado principal

La entidad principal es **Intercambio**.

Sus estados son:

* **PENDIENTE:** se ha solicitado un intercambio.
* **ACEPTADO:** el propietario acepta la solicitud.
* **RECHAZADO:** el propietario rechaza la solicitud.
* **COMPLETADO:** se realizó la entrega del libro.

El flujo permitido es:

**PENDIENTE → ACEPTADO → COMPLETADO**

También se permite:

**PENDIENTE → RECHAZADO**

Un intercambio rechazado no puede volver a aceptarse.

## 5. Reglas principales

1. Solo usuarios registrados pueden publicar libros.
2. Un libro pertenece a un usuario.
3. Solo los libros disponibles pueden participar en nuevos intercambios.
4. Todo intercambio comienza como PENDIENTE.
5. El propietario puede aceptar o rechazar una solicitud.
6. Un intercambio aceptado puede finalizar como COMPLETADO.
7. Un intercambio RECHAZADO no puede volver a ACEPTADO.
8. Los estados solamente pueden utilizar los valores definidos por el sistema.

## 6. Roles

### Usuario

Puede:

* Registrar y consultar libros.
* Solicitar intercambios.
* Consultar sus intercambios.

### Administrador

Puede:

* Gestionar usuarios.
* Gestionar libros.
* Supervisar intercambios.
* Revisar información del sistema.

## 7. Decisiones técnicas

Para los identificadores se utilizará el tipo entero.

Los nombres, títulos, autores, correos y estados se manejarán como texto.

Las relaciones entre entidades se realizarán mediante identificadores de los registros relacionados.

Los estados se manejarán mediante valores definidos para evitar información diferente o incorrecta.

## 8. Criterios de la primera versión

La primera versión se considera funcional cuando:

* Se puedan registrar libros.
* Se puedan consultar libros.
* Se pueda crear una solicitud de intercambio.
* Se pueda aceptar o rechazar una solicitud.
* Se pueda completar un intercambio.
* Se respeten las reglas y estados definidos.

## 9. Uso de IA

Se utilizó inteligencia artificial como apoyo para organizar y redactar el documento.

El equipo revisó y adaptó el contenido final de acuerdo con los requerimientos del Hito 1 y el funcionamiento planteado para el proyecto.
