# Sistema de Gestión de Solicitudes y Documentos

> **Arquitectura de microservicios independientes para gestión de solicitudes de trabajo y documentos adjuntos**

## 📋 ¿Qué es este proyecto?

Este repositorio contiene **dos microservicios independientes** desarrollados en **Go** que trabajan en conjunto para gestionar solicitudes de trabajo y sus documentos adjuntos:

### 🎯 Microservicios

| Microservicio | Puerto | Descripción | Base de Datos |
|--------------|--------|-------------|---------------|
| **Solicitudes** | 8082 | Gestiona solicitudes de trabajo (CRUD completo) | MySQL (puerto 3009) |
| **Documentos** | 8083 | Gestiona documentos adjuntos a las solicitudes | MySQL (puerto 3010) |

### 🔗 Características principales

- ✅ **Arquitectura de microservicios** - Servicios independientes y desacoplados
- ✅ **Comunicación HTTP REST** - Los microservicios se comunican mediante APIs REST
- ✅ **Soft Delete** - Los registros no se eliminan físicamente, se marcan con `deleted_at`
- ✅ **Eliminación en cascada** - Al eliminar una solicitud, sus documentos también se marcan como eliminados
- ✅ **CORS habilitado** - Listo para consumir desde frontend (Vue.js, React, etc.)
- ✅ **Validación de datos** - Validación estricta de campos y tipos de datos
- ✅ **Tests unitarios** - Cobertura del 87.0% en el microservicio de Solicitudes

## Arquitectura del Sistema

```
Solicitudes (Puerto 8082)  -----HTTP-----> Documentos (Puerto 8083)
        |                                          |
        |                                          |
        v                                          v
MySQL DB (solicitudes)                    MySQL DB (documentos)
```

## 🚀 Cómo ejecutar el proyecto

### 📋 Prerrequisitos

Asegúrate de tener instalado lo siguiente en tu máquina:

- **Go 1.20+** - [Descargar aquí](https://golang.org/dl/)
- **Docker & Docker Compose** - [Descargar aquí](https://www.docker.com/products/docker-desktop)
- **Git** - [Descargar aquí](https://git-scm.com/downloads)
- **Make** (opcional, pero recomendado para Windows) - [Descargar aquí](https://gnuwin32.sourceforge.net/packages/make.htm)

### 1️⃣ Clonar el repositorio

```bash
git clone https://github.com/KarlaR3it/Solicitudes-documentos-Go.git
cd Solicitudes-documentos-Go
```

### 2️⃣ Ejecutar el microservicio de Solicitudes

Abre una terminal y ejecuta:

```bash
cd solicitudes
make install    # Instala dependencias y levanta MySQL con Docker
make start      # Inicia el servidor en puerto 8082
```

**✅ Comandos disponibles con Makefile:**
- `make help` - Ver todos los comandos disponibles
- `make install` - Instalar dependencias y levantar Docker
- `make start` - Iniciar el servidor
- `make test` - Ejecutar pruebas unitarias
- `make test-cover` - Ejecutar pruebas con reporte de cobertura

### 3️⃣ Ejecutar el microservicio de Documentos

Abre **otra terminal** y ejecuta:

```bash
cd documentos
make install    # Instala dependencias y levanta MySQL con Docker
make start      # Inicia el servidor en puerto 8083
```

**✅ Comandos disponibles con Makefile:**
- `make help` - Ver todos los comandos disponibles
- `make install` - Instalar dependencias y levantar Docker
- `make start` - Iniciar el servidor

### 4️⃣ Verificar que todo funciona

Abre tu navegador o Postman y verifica:

- **Solicitudes**: http://localhost:8082/solicitudes
- **Documentos**: http://localhost:8083/documentos

Si ves una respuesta JSON (aunque sea vacía `[]`), ¡todo está funcionando correctamente! ✅

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

**Cobertura actual: 87.0%** ✅ (Meta: >75%)

### **🎯 Capas Probadas**

| **Capa**          | **Tests** | **Cobertura** | **Descripción**                       |
| ----------------- | --------- | ------------- | ------------------------------------- |
| **🌐 Endpoints**  | 16 tests  | ~77%          | Controladores HTTP (handlers) con Gin |
| **⚙️ Services**   | 22 tests  | ~95%          | Lógica de negocio y validaciones      |
| **🗄️ Repository** | 13 tests  | ~91%          | Acceso a datos con GORM + sqlmock     |
| **📦 Models**     | 2 tests   | ~80%          | Métodos de dominio y conversión       |

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
- ✅ Eliminación en cascada de documentos
- ✅ Continúa aunque falle eliminación de documentos

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
├── solicitudes/                    # 📋 Microservicio Solicitudes
│   ├── cmd/
│   │   └── main.go                # Punto de entrada
│   ├── internal/solicitud/        # Lógica de negocio
│   │   ├── endpoint.go           # 🌐 HTTP handlers (controladores)
│   │   ├── service.go            # ⚙️ Lógica de negocio y validaciones
│   │   ├── repository.go         # 🗄️ Acceso a datos (GORM)
│   │   ├── solicitud.go          # 📦 Modelos de dominio
│   │   ├── endpoint_test.go      # 🧪 Tests de endpoints
│   │   ├── service_test.go       # 🧪 Tests de servicios
│   │   └── repository_test.go    # 🧪 Tests de repositorio
│   ├── pkg/
│   │   ├── bootstrap/            # Inicialización (DB, Logger, Env)
│   │   ├── handler/              # Configuración de rutas
│   │   └── httpclient/           # Cliente HTTP para Documentos
│   ├── Makefile                  # 🔧 Comandos automatizados
│   ├── docker-compose.yml        # 🐳 MySQL container (puerto 3009)
│   ├── .env                      # Variables de entorno
│   └── coverage.html             # 📊 Reporte de cobertura
│
├── documentos/                    # 📄 Microservicio Documentos
│   ├── cmd/
│   │   └── main.go               # Punto de entrada
│   ├── internal/documento/       # Lógica de negocio
│   │   ├── endpoint.go          # 🌐 HTTP handlers
│   │   ├── service.go           # ⚙️ Lógica de negocio
│   │   ├── repository.go        # 🗄️ Acceso a datos (GORM)
│   │   └── documento.go         # 📦 Modelos de dominio
│   ├── pkg/
│   │   ├── bootstrap/           # Inicialización (DB, Logger, Env)
│   │   ├── handler/             # Configuración de rutas
│   │   └── httpclient/          # Cliente HTTP para Solicitudes
│   ├── Makefile                 # 🔧 Comandos automatizados
│   ├── docker-compose.yml       # 🐳 MySQL container (puerto 3010)
│   └── .env                     # Variables de entorno
│
├── go.work                       # Go workspace (multi-módulo)
└── README.md                     # 📖 Este archivo
```

## 🔧 Tecnologías utilizadas

| Tecnología | Versión | Propósito |
|-----------|---------|-----------|
| **Go** | 1.20+ | Lenguaje de programación |
| **Gin** | v1.9+ | Framework web HTTP |
| **GORM** | v1.25+ | ORM para base de datos |
| **MySQL** | 8.0+ | Base de datos relacional |
| **Docker** | 20.10+ | Contenedores para BD |
| **Testify** | v1.8+ | Framework de testing |

## ❓ Troubleshooting (Solución de problemas)

### Error: "Puerto ya en uso"

Si obtienes un error como `bind: address already in use`:

```bash
# Windows
netstat -ano | findstr :8082
taskkill /PID <PID> /F

# Linux/Mac
lsof -ti:8082 | xargs kill -9
```

### Error: "Cannot connect to MySQL"

Verifica que Docker esté corriendo:

```bash
docker ps
```

Si no ves los contenedores de MySQL, levántalos:

```bash
cd solicitudes
docker compose up -d

cd ../documentos
docker compose up -d
```

### Error: "deleted_at column doesn't exist"

Necesitas agregar la columna `deleted_at` a la tabla `documentos`:

```sql
ALTER TABLE documentos ADD COLUMN deleted_at TIMESTAMP NULL;
CREATE INDEX idx_documentos_deleted_at ON documentos(deleted_at);
```

### Los cambios no se reflejan

Asegúrate de reiniciar ambos servicios después de hacer cambios en el código:

1. Detén el servidor (Ctrl+C)
2. Ejecuta `make start` nuevamente

## 🤝 Contribuir

Si encuentras algún bug o tienes sugerencias:

1. Abre un **Issue** en GitHub
2. Haz un **Fork** del repositorio
3. Crea una **rama** con tu feature (`git checkout -b feature/nueva-funcionalidad`)
4. Haz **commit** de tus cambios (`git commit -m 'Agrega nueva funcionalidad'`)
5. Haz **push** a la rama (`git push origin feature/nueva-funcionalidad`)
6. Abre un **Pull Request**

## 📝 Licencia

Este proyecto es de código abierto y está disponible bajo la licencia MIT.

## 👥 Autores

- **Karla Ramírez** - [GitHub](https://github.com/KarlaR3it)

## 🌐 Endpoints API - Operaciones CRUD

### 📋 Solicitudes (Puerto 8082)

| Método | Endpoint | Descripción | Tipo de Eliminación |
|--------|----------|-------------|---------------------|
| `GET` | `/solicitudes` | Listar todas las solicitudes (con filtros opcionales) | - |
| `POST` | `/solicitudes` | Crear nueva solicitud | - |
| `GET` | `/solicitudes/:id` | Obtener solicitud por ID (sin documentos) | - |
| `GET` | `/solicitudes/:id/con-documentos` | Obtener solicitud con sus documentos adjuntos | - |
| `PATCH` | `/solicitudes/:id` | Actualizar solicitud (parcial) | - |
| `DELETE` | `/solicitudes/:id` | **Eliminar solicitud (Soft Delete)** | ⚠️ **Soft Delete** |

### 📄 Documentos (Puerto 8083)

| Método | Endpoint | Descripción | Tipo de Eliminación |
|--------|----------|-------------|---------------------|
| `GET` | `/documentos` | Listar todos los documentos (con filtros opcionales) | - |
| `POST` | `/documentos` | Crear nuevo documento | - |
| `GET` | `/documentos/:id` | Obtener documento por ID | - |
| `PATCH` | `/documentos/:id` | Actualizar documento (parcial) | - |
| `DELETE` | `/documentos/:id` | **Eliminar documento (Soft Delete)** | ⚠️ **Soft Delete** |

### ⚠️ Importante: Soft Delete

**¿Qué es Soft Delete?**

Cuando eliminas una solicitud o documento usando `DELETE`, **NO se borra físicamente de la base de datos**. En su lugar:

1. ✅ Se marca el registro con una fecha en la columna `deleted_at`
2. ✅ El registro deja de aparecer en las consultas normales (GET)
3. ✅ En el frontend, el elemento desaparecerá de la interfaz
4. ✅ En la base de datos, el registro sigue existiendo pero está marcado como eliminado

**Eliminación en cascada:**
- Al eliminar una **solicitud**, todos sus **documentos asociados** también se marcan como eliminados automáticamente
- Esto mantiene la integridad referencial entre ambos microservicios

**Ejemplo en la BD:**
```sql
-- Antes de DELETE
id | titulo | deleted_at
1  | "Solicitud 1" | NULL

-- Después de DELETE
id | titulo | deleted_at
1  | "Solicitud 1" | 2025-10-17 13:45:00
```

## 🧪 Probar los endpoints

### Opción 1: Postman (Recomendado)

Puedes importar la colección de Postman con todos los endpoints configurados:

📦 **[Descargar colección de Postman](#)** *(Comparte el link de tu colección aquí)*

### Opción 2: cURL

**Crear una solicitud:**
```bash
curl -X POST http://localhost:8082/solicitudes \
  -H "Content-Type: application/json" \
  -d '{
    "titulo": "Solicitud de prueba",
    "estado": "pendiente",
    "area": "Tecnología",
    "pais": "México",
    "localizacion": "CDMX",
    "rango_renta_desde": 50000,
    "rango_renta_hasta": 80000,
    "fecha_inicio_proyecto": "2025-11-01",
    "usuario_id": 1
  }'
```

**Listar solicitudes:**
```bash
curl http://localhost:8082/solicitudes
```

**Eliminar solicitud (Soft Delete):**
```bash
curl -X DELETE http://localhost:8082/solicitudes/1
```

### Opción 3: Frontend Vue.js

Este proyecto está preparado para consumirse desde el frontend desarrollado en Vue 3:

🔗 **[Frontend Vue.js](https://github.com/KarlaR3it/Solicitudes-documentos-Vue.git)**

## 🛠️ Comandos útiles

### Microservicio Solicitudes

```bash
cd solicitudes
make help           # Ver todos los comandos disponibles
make install        # Instalar dependencias + Docker
make start          # Iniciar servidor
make test           # Ejecutar pruebas
make test-cover     # Ejecutar pruebas con cobertura
make clean          # Limpiar archivos temporales
```

### Microservicio Documentos

```bash
cd documentos
make help           # Ver todos los comandos disponibles
make install        # Instalar dependencias + Docker
make start          # Iniciar servidor
```
