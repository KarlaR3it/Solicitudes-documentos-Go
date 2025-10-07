package solicitud

import "time"

//Entidad de solicitud en la base de datos
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
	RentaDesde               int       `gorm:"type:int;not null" json:"rentaDesde"`
	RentaHasta               int       `gorm:"type:int;not null" json:"rentaHasta"`
	ModalidadTrabajo         string    `gorm:"type:varchar(50);not null" json:"modalidad_trabajo"`
	TipoServicio             string    `gorm:"type:varchar(30);not null" json:"tipo_servicio"`
	NivelExperiencia         string    `gorm:"type:varchar(30);not null" json:"nivel_experiencia"`
	FechaInicioProyecto      time.Time `gorm:"type:date;not null" json:"fecha_inicio_proyecto"`
	CreatedAt                time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt                time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	// CASCADE: Cuando se elimine un usuario, se eliminarán automáticamente todas sus solicitudes
	UsuarioID                uint      `gorm:"not null;constraint:OnDelete:CASCADE" json:"usuario_id"`
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
	UsuarioID                uint   `json:"usuario_id" binding:"required"`
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
