# Sistema para Gestionar Solicitudes y Documentos

> **Arquitectura de microservicios independientes para gestión de solicitudes de trabajo y documentos adjuntos**

Este repositorio contiene un microservicio para **Solicitudes** y un microservicio para **Documentos**.

Una solicitud puede tener uno o más documentos adjuntos.

## ¿Qué es este proyecto?

Este repositorio contiene **dos microservicios independientes** desarrollados en Go que trabajan en conjunto:

- **Microservicio Solicitudes** - Gestiona solicitudes de trabajo (crear, actualizar, consultar, eliminar) - Puerto: 8082
- **Microservicio Documentos** - Gestiona documentos adjuntos asociados a las solicitudes - Puerto: 8083

Los microservicios se comunican entre sí mediante HTTP REST APIs y cada uno mantiene su propia base de datos MySQL.

## Arquitectura del Sistema

```
Solicitudes (Puerto 8082)  -----HTTP-----> Documentos (Puerto 8083)
        |                                          |
        |                                          |
        v                                          v
MySQL DB (solicitudes)                    MySQL DB (documentos)
```

## Cómo ejecutar el proyecto

### Prerrequisitos

### Prerrequisitos- [Go](https://golang.org/dl/) 1.20 o superior

- **Go 1.20+** - [Descargar aquí](https://golang.org/dl/)- [Docker](https://www.docker.com/products/docker-desktop) 20.10 o superior

- **Docker & Docker Compose** - [Descargar aquí](https://www.docker.com/products/docker-desktop)
- **Git** - [Descargar aquí](https://git-scm.com/downloads)

### Clonar el repositorio

```bash
git clone https://github.com/KarlaR3it/Solicitudes-documentos-Go.git
cd Solicitudes-documentos-Go
```

### Ejecutar el microservicio de Solicitudes

```bash
cd solicitudes
make install    # Instala dependencias y levanta MySQL con Docker
make start      # Inicia el servidor en puerto 8082
```

**Nota:** El microservicio de Solicitudes tiene un `Makefile` configurado, por lo que puedes usar comandos `make`.

### Ejecutar el microservicio de Documentos (en otra terminal)

```bash
cd documentos
docker compose up -d        # Levanta MySQL con Docker
go run cmd/main.go         # Inicia el servidor en puerto 8083
```

**Nota:** El microservicio de Documentos NO tiene `Makefile`, usa comandos directos de Go.

### 4️⃣ Verificar que funciona

- **Solicitudes**: http://localhost:8082/solicitudes
- **Documentos**: http://localhost:8083/documentos

## 🧪 Cómo ejecutar las pruebas

> **⚠️ Importante**: Las pruebas unitarias están implementadas únicamente en el **microservicio de Solicitudes**.

### Ejecutar todas las pruebas

```bash
cd solicitudes
make test-cover
```

### Solo ejecutar pruebas (sin cobertura)

```bash
cd solicitudes
make test
```

### Ver reporte de cobertura

Después de ejecutar `make test-cover`, se genera automáticamente:

- **coverage.out** - Datos de cobertura
- **coverage.html** - Reporte visual (ábrelo en tu navegador)

## 📊 ¿Qué cubren las pruebas?

**Cobertura actual: 87.6%** ✅ (Meta: >75%)

### **🎯 Capas Probadas**

| **Capa**          | **Tests** | **Cobertura** | **Descripción**                       |
| ----------------- | --------- | ------------- | ------------------------------------- |
| **🌐 Endpoints**  | 16 tests  | 77.2%         | Controladores HTTP (handlers) con Gin |
| **⚙️ Services**   | 21 tests  | 94.7%         | Lógica de negocio y validaciones      |
| **🗄️ Repository** | 13 tests  | 90.5%         | Acceso a datos con GORM + sqlmock     |
| **📦 Models**     | 2 tests   | 80.0%         | Métodos de dominio y conversión       |

### **🔍 Operaciones CRUD Probadas**

#### ✅ **CREATE (Crear solicitudes)**

- ✅ Creación exitosa (200)
- ❌ JSON inválido (400)
- ❌ Error de base de datos (500)
- ✅ Validaciones de campos requeridos
- ✅ Valores por defecto

#### ✅ **READ (Consultar solicitudes)**

- ✅ Obtener todas las solicitudes con filtros
- ✅ Obtener solicitud por ID
- ✅ Obtener solicitud con documentos adjuntos
- ❌ ID inválido (400)
- ❌ No encontrado (404)
- ✅ Paginación y filtros avanzados

#### ✅ **UPDATE (Actualizar solicitudes)**

- ✅ Actualización exitosa (200)
- ❌ Campos prohibidos (400)
- ❌ JSON inválido (400)
- ❌ ID inválido (400)
- ❌ No encontrado (500)
- ✅ Actualización parcial de campos

#### ✅ **DELETE (Eliminar solicitudes)**

- ✅ Eliminación exitosa (200)
- ❌ ID inválido (400)
- ❌ No encontrado (500)
- ✅ Soft delete con GORM

### **🛡️ Validaciones de Negocio Probadas**

- ✅ Título requerido
- ✅ Área requerida
- ✅ País requerido
- ✅ Localización requerida
- ✅ Usuario ID requerido
- ✅ Formato de fecha válido (YYYY-MM-DD)
- ✅ Rango de renta válido (desde ≤ hasta)
- ✅ Estado por defecto ("pendiente")

### **🔧 Tipos de Tests Implementados**

| **Tipo**              | **Herramienta**            | **Propósito**                 |
| --------------------- | -------------------------- | ----------------------------- |
| **Unit Tests**        | `testify/assert`           | Verificar lógica de funciones |
| **Integration Tests** | `httptest` + `gin`         | Probar endpoints HTTP         |
| **Mock Tests**        | `testify/mock` + `sqlmock` | Simular dependencias externas |
| **Coverage Tests**    | `go tool cover`            | Medir cobertura de código     |

## 📁 Estructura del proyecto

```
Solicitudes-documentos-Go/
├── solicitudes/              # 📋 Microservicio Solicitudes
│   ├── cmd/main.go          # Punto de entrada
│   ├── internal/solicitud/   # Lógica de negocio
│   │   ├── endpoint.go      # 🌐 HTTP handlers
│   │   ├── service.go       # ⚙️ Lógica de negocio
│   │   ├── repository.go    # 🗄️ Acceso a datos
│   │   ├── solicitud.go     # 📦 Modelo de dominio
│   │   └── *_test.go        # 🧪 Tests unitarios
│   ├── pkg/                 # Utilidades compartidas
│   ├── Makefile            # 🔧 Comandos automatizados
│   ├── docker-compose.yml  # 🐳 MySQL container
│   └── coverage.html       # 📊 Reporte de cobertura
│
├── documentos/              # 📄 Microservicio Documentos
│   ├── cmd/main.go         # Punto de entrada
│   ├── internal/documento/ # Lógica de negocio
│   └── docker-compose.yml # 🐳 MySQL container
│
├── go.work                 # Go workspace
└── README.md              # 📖 Este archivo
```

## 🛠️ Comandos útiles

### Microservicio Solicitudes (con Makefile)

```bash
cd solicitudes
make help           # Ver todos los comandos disponibles
make install        # Instalar dependencias + Docker
make start          # Iniciar servidor
make test           # Ejecutar pruebas
make test-cover     # Ejecutar pruebas con cobertura
make clean          # Limpiar archivos temporales
```

### Microservicio Documentos (sin Makefile)

```bash
cd documentos
docker compose up -d    # Instalar dependencias + Docker
go run cmd/main.go     # Iniciar servidor
go test ./...          # Ejecutar pruebas (cuando estén implementadas)
```

## 🌐 Endpoints disponibles

### 📋 Solicitudes (Puerto 8082)

- `GET /solicitudes` - Listar todas las solicitudes
- `POST /solicitudes` - Crear nueva solicitud
- `GET /solicitudes/{id}` - Obtener solicitud por ID
- `GET /solicitudes/{id}/con-documentos` - Obtener solicitud con documentos
- `PATCH /solicitudes/{id}` - Actualizar solicitud
- `DELETE /solicitudes/{id}` - Eliminar solicitud

### 📄 Documentos (Puerto 8083)

- `GET /documentos` - Listar todos los documentos
- `POST /documentos` - Crear nuevo documento
- `GET /documentos/{id}` - Obtener documento por ID
- `PATCH /documentos/{id}` - Actualizar documento
- `DELETE /documentos/{id}` - Eliminar documento
