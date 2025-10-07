# 📄 Microservicio de Documentos

Este microservicio se encarga de la gestión de documentos adjuntos en el sistema. Es parte de la arquitectura de microservicios del proyecto.

## 🚀 Empezando

### Requisitos Previos

Asegúrate de tener instalados los requisitos generales del proyecto principal. Consulta el [README principal](../../README.md) para más detalles.

### Configuración del Proyecto

1. **Configura las variables de entorno**:
   - Copia el archivo `.env.example` a `.env`
   - Edita el archivo `.env` según tu configuración

2. **Variables de entorno principales**:
   ```
   DB_HOST=db
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=root
   DB_NAME=documentos_db
   SERVICE_PORT=8083
   ```

### 1. Iniciar la Aplicación

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

### 3. Configuración del Entorno

1. **Configura las variables de entorno** (antes de iniciar los contenedores):
   - Asegúrate de que el archivo `.env` tenga estos valores:
     ```
     DB_HOST=db
     DB_PORT=3306
     DB_USER=root
     DB_PASSWORD=root
     DB_NAME=documentos_db
     SERVICE_PORT=8083
     ```
   - Estos valores son los predeterminados y deberían funcionar con la configuración de Docker Compose

### 4. Iniciar la Aplicación

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

### Usando Postman (Recomendado)

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
   - Deberías recibir una respuesta con el documento creado

3. **Otras operaciones que puedes probar**:
   - **Listar todos los documentos**: `GET http://localhost:8083/documentos`
   - **Obtener un documento por ID**: `GET http://localhost:8083/documentos/1`
   - **Actualizar un documento**:
     ```http
     PATCH http://localhost:8083/documentos/1
     Content-Type: application/json
     
     {
         "extension": "doc",
         "nombre_archivo": "CV_prof"
     }
     ```
     > 💡 Puedes actualizar solo los campos que necesites, no es necesario enviar todos los campos.

   - **Eliminar un documento**:
     ```http
     DELETE http://localhost:8083/documentos/1
     ```

   Recuerda reemplazar `1` por el ID real del documento que quieras consultar, actualizar o eliminar.

## 🔐 Autenticación

> **Nota**: Actualmente, la autenticación está en desarrollo. Para probar los endpoints, asegúrate de incluir el siguiente ID de usuario en el cuerpo de tus peticiones:
> 
> ```json
> "usuario_id": 1
> ```

## 📚 Endpoints Disponibles

{{ ... }}
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
