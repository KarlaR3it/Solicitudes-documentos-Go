# 📋 Microservicio de Solicitudes

Este microservicio se encarga de la gestión de solicitudes en el sistema. Es parte de la arquitectura de microservicios del proyecto.

## 🚀 Empezando

### Requisitos Previos

Asegúrate de tener instalados los requisitos generales del proyecto principal. Consulta el [README principal](../../README.md) para más detalles.

### Configuración del Entorno

1. **Configura las variables de entorno**:
   - Copia el archivo `.env.example` a `.env`
   - Edita el archivo `.env` según tu configuración

2. **Variables de entorno principales**:
   ```
   DB_HOST=db
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=root
   DB_NAME=solicitudes_db
   SERVICE_PORT=8082
   ```

### 3. Iniciar con Docker Compose (Recomendado)

```bash
docker-compose up --build
```

### 4. O Iniciar Localmente

1. **Instalar dependencias**:
   ```bash
   go mod download
   ```

2. **Iniciar el servidor**:
   ```bash
   go run cmd/main.go
   ```

## 🔐 Autenticación

> **Nota**: Actualmente, la autenticación está en desarrollo. Para probar los endpoints, utiliza el siguiente ID de usuario hardcodeado en tus peticiones:
> 
> ```
> usuario_id: 1
> ```

## 📚 Endpoints Disponibles

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET`  | `/solicitudes` | Lista todas las solicitudes |
| `GET`  | `/solicitudes/{id}` | Obtiene una solicitud por ID |
| `POST` | `/solicitudes` | Crea una nueva solicitud |
| `PATCH`| `/solicitudes/{id}` | Actualiza una solicitud existente |
| `DELETE`| `/solicitudes/{id}` | Elimina una solicitud |

## 🧪 Probar la API

Puedes probar la API usando [Postman](https://www.postman.com/) o `curl`:

```bash
# Obtener todas las solicitudes
curl http://localhost:8082/solicitudes

# Crear una nueva solicitud
curl -X POST http://localhost:8082/solicitudes \
  -H "Content-Type: application/json" \
  -d '{"titulo":"Nueva solicitud","descripcion":"Descripción de la solicitud"}'
```

## 🏗 Estructura del Proyecto

```
solicitudes/
├── cmd/           # Punto de entrada de la aplicación
├── internal/      # Código interno del servicio
│   └── solicitud/ # Lógica de negocio
├── pkg/           # Bibliotecas compartidas
├── .env.example   # Plantilla de configuración
└── docker-compose.yml  # Configuración de Docker
```

## 🚨 Solución de Problemas Comunes

### No se puede conectar a la base de datos
- Verifica que el contenedor de MySQL esté en ejecución:
  ```bash
  docker ps  # Deberías ver el contenedor mysql-solicitudes
  ```
- Revisa los logs del contenedor:
  ```bash
  docker logs mysql-solicitudes
  ```

### El puerto está en uso
- Cambia el puerto en el archivo `.env`:
  ```
  SERVICE_PORT=8083
  ```

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo [LICENSE](LICENSE) para más detalles.
