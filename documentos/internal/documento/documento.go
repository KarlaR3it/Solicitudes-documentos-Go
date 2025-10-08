package documento

import "time"

type Documento struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Extension     string    `gorm:"type:varchar(5);not null" json:"extension"`
	NombreArchivo string    `gorm:"type:varchar(255);not null" json:"nombre_archivo"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	// CASCADE: Cuando se elimine una solicitud, se eliminarán automáticamente todos sus documentos
	SolicitudID uint `gorm:"not null;constraint:OnDelete:CASCADE" json:"-"`
}

// DocumentoResponse es la estructura de respuesta para los documentos
type DocumentoResponse struct {
	ID            uint      `json:"id"`
	Extension     string    `json:"extension"`
	NombreArchivo string    `json:"nombre_archivo"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Solicitud     struct {
		ID     uint   `json:"solicitud_id"`
		Titulo string `json:"titulo"`
		Area   string `json:"area"`
	} `json:"solicitud"`
}

// TableName especifica el nombre de la tabla
func (Documento) TableName() string {
	return "documentos"
}

//CreateReq representa la petición para crear un documento
type CreateReq struct {
	Extension     string `json:"extension" binding:"required"`
	NombreArchivo string `json:"nombre_archivo" binding:"required"`
	SolicitudID   uint   `json:"solicitud_id" binding:"required"`
}

//UpdateReq representa la petición para actualizar un documento
type UpdateReq struct {
	Extension     *string `json:"extension"`
	NombreArchivo *string `json:"nombre_archivo"`
}

//GetAllReq representa los filtros para obtener documentos
type GetAllReq struct {
	Extension     string
	NombreArchivo string
	SolicitudID   uint // para ver los documentos de una solicitud en particular
	Limit         int
	Page          int
}
