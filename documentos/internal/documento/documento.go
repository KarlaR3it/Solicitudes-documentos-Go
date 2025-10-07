package documento

import "time"

type Documento struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	Extension string `gorm:"type:varchar(5);not null" json:"extension"`
	NombreArchivo string `gorm:"type:varchar(255);not null" json:"nombre_archivo"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	// CASCADE: Cuando se elimine una solicitud, se eliminarán automáticamente todos sus documents
	SolicitudID uint `gorm:"not null;constraint:OnDelete:CASCADE" json:"solicitud_id"`
}

// TableName especifica el nombre de la tabla
func (Documento) TableName() string {
	return "documentos"
}

//CreateReq representa la petición para crear un documento
type CreateReq struct {
	Extension string `json:"extension" binding:"required"`
	NombreArchivo string `json:"nombre_archivo" binding:"required"`
	SolicitudID uint `json:"solicitud_id" binding:"required"`
}

//UpdateReq representa la petición para actualizar un documento
type UpdateReq struct {
	Extension *string `json:"extension"`
	NombreArchivo *string `json:"nombre_archivo"`
}

//GetAllReq representa los filtros para obtener documentos
type GetAllReq struct {
	Extension string
	NombreArchivo string
	Limit int
	Page int
}
