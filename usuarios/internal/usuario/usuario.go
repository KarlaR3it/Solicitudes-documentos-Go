package usuario

import "time"

// Entidad de usuario en la base de datos
type Usuario struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	NombreUsuario string    `gorm:"type:varchar(30);not null;uniqueIndex" json:"nombre_usuario"`
	EmailUsuario  string    `gorm:"type:varchar(150);not null;uniqueIndex" json:"email_usuario"`
	Nombre        string    `gorm:"type:varchar(50);not null" json:"nombre"`
	Apellidos     string    `gorm:"type:varchar(50);not null" json:"apellidos"`
	Idioma        string    `gorm:"type:varchar(30)" json:"idioma"`
	ZonaHoraria   string    `gorm:"type:varchar(30);column:zona_horaria" json:"zona_horaria"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName especifica el nombre de la tabla
func (Usuario) TableName() string {
	return "usuarios"
}

// CreateReq representa la petición para crear un usuario
type CreateReq struct {
	NombreUsuario string `json:"nombre_usuario" binding:"required"`
	EmailUsuario  string `json:"email_usuario" binding:"required,email"`
	Nombre        string `json:"nombre" binding:"required"`
	Apellidos     string `json:"apellidos" binding:"required"`
	Idioma        string `json:"idioma" `
	ZonaHoraria   string `json:"zona_horaria"`
}

// UpdateReq representa la petición para actualizar un usuario
type UpdateReq struct {
	NombreUsuario *string `json:"nombre_usuario"`
	EmailUsuario  *string `json:"email_usuario"`
	Nombre        *string `json:"nombre"`
	Apellidos     *string `json:"apellidos"`
	Idioma        *string `json:"idioma"`
	ZonaHoraria   *string `json:"zona_horaria"`
}

// GetAllReq representa los filtros para obtener usuarios
type GetAllReq struct {
	Nombre        string
	Apellidos     string
	EmailUsuario  string
	NombreUsuario string
	Limit         int
	Page          int
}
