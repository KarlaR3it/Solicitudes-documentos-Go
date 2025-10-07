# 📄 Documentos Service

## 👋 ¡Bienvenid@!
Este documento te guiará paso a paso para configurar y ejecutar el servicio de gestión de documentos.

## 📋 ¿Qué es este proyecto?

Este es un **microservicio** diseñado para la gestión de documentos. Proporciona una **API REST** que permite:

- Crear, leer, actualizar y eliminar documentos
- Almacenar información en una base de datos MySQL
- Integrarse con otros servicios mediante HTTP

### Características principales:
- **Arquitectura de microservicios**: Desplegable de forma independiente
- **API RESTful**: Interfaz estándar para integración
- **Base de datos MySQL**: Almacenamiento persistente de documentos
- **Configuración mediante variables de entorno**: Fácil despliegue en diferentes entornos

## 🛠️ Requisitos Previos

Antes de comenzar, necesitarás instalar estas herramientas en tu computadora:

1. **Git** - Para descargar el código
   - [Descargar Git para Windows](https://git-scm.com/download/win)
   - Al instalar, selecciona "Git from the command line and also from 3rd-party software"
   - Para verificar que se instaló correctamente, abre una nueva ventana de PowerShell y escribe:
     ```bash
     git --version
     ```
     Deberías ver un número de versión (ejemplo: git version 2.40.0).

2. **Docker Desktop** - Para ejecutar la base de datos y la aplicación
   - [Descargar Docker Desktop para Windows](https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe)
   - Sigue las instrucciones del instalador
   - Después de instalar, reinicia tu computadora
   - Abre Docker Desktop para iniciar Docker
   - Para verificar que Docker está funcionando, abre PowerShell y escribe:
     ```bash
     docker --version
     docker-compose --version
     ```
     Deberías ver números de versión para ambos comandos.

## 🚀 Empecemos: Guía Paso a Paso

### 1. Obtener el Código

1. **Abre PowerShell** (presiona `Windows + X` y selecciona "Windows PowerShell")
2. **Navega a la carpeta donde quieres guardar el proyecto**, por ejemplo:
   ```bash
   cd Documentos  # O la carpeta que prefieras
   ```
3. **Clona el repositorio** (copia y pega este comando):
   ```bash
   git clone https://github.com/KarlaR3it/Solicitudes-documentos-Go.git
   ```
4. **Entra a la carpeta del proyecto**:
   ```bash
   cd Solicitudes-documentos-Go
   ```
   
   💡 *Verás que se creó una carpeta llamada "Solicitudes-documentos-Go" con todos los archivos del proyecto.*

### 2. Configurar el Proyecto

1. **Configura las variables de entorno**:
   - Abre el archivo `.env.example` que está en la raíz del proyecto
   - Copia todo su contenido (Ctrl+A, Ctrl+C)
   - Crea un nuevo archivo llamado `.env` en la misma carpeta
   - Pega el contenido copiado (Ctrl+V)
   - Guarda el archivo (Ctrl+S)

2. **Edita el archivo `.env`** con tus credenciales:
   - Abre el archivo `.env`
   - Actualiza los valores según tu configuración (usuario, contraseña, etc.)
   - Asegúrate de que los puertos no estén en uso por otras aplicaciones

### 3. Iniciar la Aplicación

1. **Abre Docker Desktop** (si no lo has hecho ya)
   - Busca "Docker Desktop" en el menú de inicio y ábrelo
   - Espera a que el ícono de Docker en la barra de tareas deje de animarse (puede tomar unos minutos la primera vez)

2. **En PowerShell**, asegúrate de estar en la carpeta del proyecto:
   ```bash
   cd ruta\a\documentos  # Ajusta la ruta según donde hayas clonado el proyecto
   ```

3. **Inicia la aplicación con Docker** (esto puede tomar unos minutos la primera vez):
   ```bash
   docker-compose up --build
   ```

   🔍 *La aplicación comenzará a iniciarse.*

4. **Listo!** La aplicación está funcionando en:
   - [http://localhost:8083](http://localhost:8083)

   Para detener la aplicación, presiona `Ctrl + C` en la ventana de PowerShell.

### Opción Alternativa: Instalación Manual

1. **Instalar Go**
   - Descarga Go desde [golang.org](https://golang.org/dl/)
   - Verifica la instalación:
     ```bash
     go version  # Debe mostrar la versión instalada
     ```

2. **Base de Datos MySQL**
   ```bash
   docker run --name mysql-documentos \
     -e MYSQL_ROOT_PASSWORD=root \
     -e MYSQL_DATABASE=documentos_db \
     -p 3306:3306 \
     -d mysql:8.0
   ```

3. **Configuración**
   - Asegúrate de que el archivo `.env` tenga estos valores:
     ```env
     DB_HOST=localhost
     DB_PORT=3306
     DB_USER=root
     DB_PASSWORD=root
     DB_NAME=documentos_db
     SERVICE_PORT=8083
     ```

4. **Instalar dependencias e iniciar**
   ```bash
   go mod download
   go run cmd/main.go
   ```

## 🧪 Probar que todo funciona

¡Perfecto! Si has llegado hasta aquí, la aplicación debería estar funcionando. 

### Importante antes de comenzar
Antes de probar la creación de documentos, asegúrate de que exista una solicitud con el ID que vas a utilizar. Puedes verificar las solicitudes existentes en el servicio de solicitudes.

Vamos a probar los endpoints principales del servicio de documentos.

### Opción 1: Usando el Navegador (Solo para GET)

1. **Verificar que el servicio está activo**:
   - Abre tu navegador web favorito
   - Ve a: [http://localhost:8083/health](http://localhost:8083/health)
   - Deberías ver un mensaje que dice "OK"

2. **Ver documentos existentes**:
   - Ve a: [http://localhost:8083/documentos](http://localhost:8083/documentos)
   - Al principio verás una lista vacía `[]`

### Opción 2: Usando Postman (Recomendado para probar todos los métodos)

1. **Descarga e instala Postman** (si no lo tienes):
   - [Descargar Postman](https://www.postman.com/downloads/)
   - Sigue las instrucciones de instalación

2. **Prueba a crear un documento**:
   - Abre Postman
   - Crea una nueva petición con el botón "New" > "HTTP Request"
   - Configura la petición así:
     - **Método**: `POST`
     - **URL**: `http://localhost:8083/documentos`
     - **Headers**:
       ```
       Content-Type: application/json
       ```
     - **Body** (selecciona "raw" y luego "JSON"):
       ```json
       {
           "extension": "pdf",
           "nombre_archivo": "CV_profesional",
           "solicitud_id": 1
       }
       ```
   - Haz clic en "Send"
   - Deberías recibir una respuesta con el documento creado

3. **Otras operaciones que puedes probar**:
   - **Listar todos los documentos**: `GET http://localhost:8083/documentos`
   - **Obtener un documento por ID**: `GET http://localhost:8083/documentos/1`
   - **Actualizar un documento**: `PUT http://localhost:8083/documentos/1` con un body similar al de creación
   - **Eliminar un documento**: `DELETE http://localhost:8083/documentos/1`

   Recuerda reemplazar `1` por el ID real del documento que quieras consultar, actualizar o eliminar.

## 📚 Endpoints Disponibles

| Método | Ruta | Descripción | Ejemplo de Uso |
|--------|------|-------------|----------------|
| `GET`  | `/documentos` | Lista todos los documentos | `GET http://localhost:8083/documentos` |
| `GET`  | `/documentos/1` | Obtiene un documento por ID | `GET http://localhost:8083/documentos/1` |
| `POST` | `/documentos` | Crea un nuevo documento | `POST http://localhost:8083/documentos` con body JSON |
| `PATCH`  | `/documentos/1` | Actualiza un documento | `PUT http://localhost:8083/documentos/1` con body JSON |
| `DELETE` | `/documentos/1` | Elimina un documento | `DELETE http://localhost:8083/documentos/1` |



## 🚨 Solución de Problemas Comunes

### No se puede conectar a la base de datos
- Verifica que el contenedor de MySQL esté en ejecución:
  ```bash
  docker ps  # Deberías ver el contenedor mysql-documentos
  ```
- Revisa los logs de MySQL:
  ```bash
  docker logs mysql-documentos
  ```

### El puerto 8083 está en uso
- Cambia el puerto en el archivo `.env`:
  ```
  SERVICE_PORT=8084
  ```
  Y reinicia los contenedores.

### Error al instalar dependencias
Si ves errores al ejecutar `go mod download`:
```bash
go env -w GOPROXY=https://proxy.golang.org,direct
go mod download
```

## 📁 Estructura del Proyecto

```
documentos/
├── cmd/           # Punto de entrada de la aplicación
├── internal/      # Código interno del servicio
│   └── documento/ # Lógica de negocio principal
├── pkg/           # Bibliotecas compartidas
├── .env.example   # Plantilla de configuración
└── docker-compose.yml  # Configuración de Docker
```
