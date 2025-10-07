# 🏗️ Arquitectura de Microservicios

Repositorio que contiene los microservicios de Solicitudes y Documentos.
Una solicitud podría tener uno o más documentos adjuntos.

## 📦 Microservicios Disponibles

1. **Solicitudes** (`/solicitudes/`)
   - Gestión de solicitudes
   - Puerto: 8082
   - [Ver documentación](./solicitudes/README.md)

2. **Documentos** (`/documentos/`)
   - Gestión de documentos adjuntos
   - Puerto: 8083
   - [Ver documentación](./documentos/README.md)

<!-- 
3. **Usuarios** (`/usuarios/`)
   - Gestión de usuarios y autenticación
   - Puerto: 8081
   - [Ver documentación](../usuarios/README.md)
-->

> **Nota sobre autenticación**: Actualmente, la autenticación de usuarios está en desarrollo. Mientras tanto, puedes usar un ID de usuario hardcodeado en tus peticiones. Consulta la documentación de cada microservicio para más detalles.

## 🚀 Empezando

### Requisitos Previos

- [Go](https://golang.org/dl/) 1.20 o superior
- [Docker](https://www.docker.com/products/docker-desktop) 20.10 o superior
- [Docker Compose](https://docs.docker.com/compose/install/) 2.0 o superior
- [Git](https://git-scm.com/downloads)

### Clonar el Repositorio

```bash
git clone [URL_DEL_REPOSITORIO]
cd solicitudes-Go
```

### Configuración Inicial

Cada microservicio tiene su propia configuración. Por favor, consulta el README.md dentro de cada carpeta para las instrucciones específicas de configuración e instalación.

## 🏗 Estructura del Repositorio

```
solicitudes-Go/
├── documentos/     # Microservicio de documentos
├── solicitudes/    # Microservicio de solicitudes
├── go.work         # Archivo de espacio de trabajo de Go
└── go.work.sum     # Suma de verificación de dependencias
```