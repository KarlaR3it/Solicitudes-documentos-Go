package solicitud

import "time"

// Documento representa un documento asociado a una solicitud
type Documento struct {
	ID            uint   `json:"id"`
	NombreArchivo string `json:"nombre_archivo"`
	Extension     string `json:"extension"`
}

// DocumentoResponse representa un documento en las respuestas de la API
type DocumentoResponse struct {
	ID            uint   `json:"id"`
	NombreArchivo string `json:"nombre_archivo"`
	Extension     string `json:"extension"`
}

// SolicitudResponse representa la respuesta de una solicitud
type SolicitudResponse struct {
	ID                       uint                    `json:"id"`
	Titulo                   string                  `json:"titulo"`
	Estado                   string                  `json:"estado"`
	Area                     string                  `json:"area"`
	Pais                     string                  `json:"pais"`
	Localizacion             string                  `json:"localizacion"`
	NumeroVacantes           int                     `json:"numero_vacantes"`
	Descripcion              string                  `json:"descripcion"`
	BaseEducacional          string                  `json:"base_educacional"`
	ConocimientosExcluyentes string                  `json:"conocimientos_excluyentes"`
	RentaDesde               int                     `json:"renta_desde"`
	RentaHasta               int                     `json:"renta_hasta"`
	ModalidadTrabajo         string                  `json:"modalidad_trabajo"`
	TipoServicio             string                  `json:"tipo_servicio"`
	NivelExperiencia         string                  `json:"nivel_experiencia"`
	FechaInicioProyecto      time.Time               `json:"fecha_inicio_proyecto"`
	CreatedAt                time.Time               `json:"created_at"`
	UpdatedAt                time.Time               `json:"updated_at"`
	UsuarioID                *uint                   `json:"usuario_id,omitempty"`
	Documentos               []DocumentoResponse `json:"documentos,omitempty"`
}

// ToResponse convierte una Solicitud a SolicitudResponse
func (s *Solicitud) ToResponse() SolicitudResponse {
	return SolicitudResponse{
		ID:                       s.ID,
		Titulo:                   s.Titulo,
		Estado:                   s.Estado,
		Area:                     s.Area,
		Pais:                     s.Pais,
		Localizacion:             s.Localizacion,
		NumeroVacantes:           s.NumeroVacantes,
		Descripcion:              s.Descripcion,
		BaseEducacional:          s.BaseEducacional,
		ConocimientosExcluyentes: s.ConocimientosExcluyentes,
		RentaDesde:               s.RentaDesde,
		RentaHasta:               s.RentaHasta,
		ModalidadTrabajo:         s.ModalidadTrabajo,
		TipoServicio:             s.TipoServicio,
		NivelExperiencia:         s.NivelExperiencia,
		FechaInicioProyecto:      s.FechaInicioProyecto,
		CreatedAt:                s.CreatedAt,
		UpdatedAt:                s.UpdatedAt,
		UsuarioID:                s.UsuarioID,
		Documentos:               []DocumentoResponse{}, // Se llenará después si es necesario
	}
}

type Solicitud struct {
	ID                       uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Titulo                   string    `gorm:"type:varchar(200);not null" json:"titulo"`
	Estado                   string    `gorm:"type:varchar(50);not null" json:"estado"`
	Area                     string    `gorm:"type:varchar(50);not null" json:"area"`
	Pais                     string    `gorm:"type:varchar(50);not null" json:"pais"`
	Localizacion             string    `gorm:"type:varchar(50);not null" json:"localizacion"`
	NumeroVacantes           int       `gorm:"type:int;not null" json:"numero_vacantes"`
	Descripcion              string    `gorm:"type:longtext;not null" json:"descripcion"`
	BaseEducacional          string    `gorm:"type:longtext;not null" json:"base_educacional"`
	ConocimientosExcluyentes string    `gorm:"type:longtext;not null" json:"conocimientos_excluyentes"`
	RentaDesde               int       `gorm:"type:int;not null" json:"renta_desde"`
	RentaHasta               int       `gorm:"type:int;not null" json:"renta_hasta"`
	ModalidadTrabajo         string    `gorm:"type:varchar(50);not null" json:"modalidad_trabajo"`
	TipoServicio             string    `gorm:"type:varchar(30);not null" json:"tipo_servicio"`
	NivelExperiencia         string    `gorm:"type:varchar(30);not null" json:"nivel_experiencia"`
	FechaInicioProyecto      time.Time `gorm:"type:date;not null" json:"fecha_inicio_proyecto"`
	CreatedAt                time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt                time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	// UsuarioID es opcional desde que se eliminó la integración con el servicio de usuarios
	UsuarioID  *uint      `gorm:"constraint:OnDelete:SET NULL" json:"usuario_id,omitempty"`
	// Documentos es una lista de documentos asociados a la solicitud
	Documentos []Documento `gorm:"-" json:"documentos,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Solicitud) TableName() string {
	return "solicitudes"
}

// CreateReq representa la petición para crear una solicitud
type CreateReq struct {
	Titulo                   string `json:"titulo" binding:"required"`
	Estado                   string `json:"estado" binding:"required"`
	Area                     string `json:"area" binding:"required"`
	Pais                     string `json:"pais" binding:"required"`
	Localizacion             string `json:"localizacion" binding:"required"`
	NumeroVacantes           int    `json:"numero_vacantes" binding:"required"`
	Descripcion              string `json:"descripcion" binding:"required"`
	BaseEducacional          string `json:"base_educacional" binding:"required"`
	ConocimientosExcluyentes string `json:"conocimientos_excluyentes" binding:"required"`
	RentaDesde               int    `json:"renta_desde" binding:"required"`
	RentaHasta               int    `json:"renta_hasta" binding:"required"`
	ModalidadTrabajo         string `json:"modalidad_trabajo" binding:"required"`
	TipoServicio             string `json:"tipo_servicio" binding:"required"`
	NivelExperiencia         string `json:"nivel_experiencia" binding:"required"`
	FechaInicioProyecto      string `json:"fecha_inicio_proyecto" binding:"required"`
	UsuarioID                *uint  `json:"usuario_id,omitempty"`
}

// UpdateReq representa la petición para actualizar una solicitud
type UpdateReq struct {
	Titulo                   *string `json:"titulo"`
	Estado                   *string `json:"estado"`
	Area                     *string `json:"area"`
	Pais                     *string `json:"pais"`
	Localizacion             *string `json:"localizacion"`
	NumeroVacantes           *int    `json:"numero_vacantes"`
	Descripcion              *string `json:"descripcion"`
	BaseEducacional          *string `json:"base_educacional"`
	ConocimientosExcluyentes *string `json:"conocimientos_excluyentes"`
	RentaDesde               *int    `json:"renta_desde"`
	RentaHasta               *int    `json:"renta_hasta"`
	ModalidadTrabajo         *string `json:"modalidad_trabajo"`
	TipoServicio             *string `json:"tipo_servicio"`
	NivelExperiencia         *string `json:"nivel_experiencia"`
	FechaInicioProyecto      *string `json:"fecha_inicio_proyecto"`
}

//GetAll Req representa los filtros para obtener solicitudes
type GetAllReq struct {
	Titulo              string
	Estado              string
	Area                string
	Pais                string
	NumeroVacantes      int
	RentaDesde          int
	RentaHasta          int
	ModalidadTrabajo    string
	TipoServicio        string
	NivelExperiencia    string
	FechaInicioProyecto string
	Limit               int
	Page                int
}
