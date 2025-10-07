package solicitud

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, solicitud *Solicitud) error
	GetAll(ctx context.Context, filters GetAllReq) ([]Solicitud, error)
	GetByID(ctx context.Context, id uint) (*Solicitud, error)
	Update(ctx context.Context, id uint, req UpdateReq) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, solicitud *Solicitud) error {
	return r.db.WithContext(ctx).Create(solicitud).Error
}

func (r *repository) GetAll(ctx context.Context, filters GetAllReq) ([]Solicitud, error) {
	var solicitudes []Solicitud
	query := r.db.WithContext(ctx).Model(&Solicitud{})

	//Aplicar filtros
	if filters.Titulo != "" {
		query = query.Where("titulo LIKE ?", "%"+filters.Titulo+"%")
	}
	if filters.Estado != "" {
		query = query.Where("estado LIKE ?", "%"+filters.Estado+"%")
	}
	if filters.Area != "" {
		query = query.Where("area LIKE ?", "%"+filters.Area+"%")
	}
	if filters.Pais != "" {
		query = query.Where("pais LIKE ?", "%"+filters.Pais+"%")
	}
	if filters.NumeroVacantes != 0 {
		query = query.Where("numero_vacantes = ?", filters.NumeroVacantes)
	}
	if filters.RentaDesde != 0 {
		query = query.Where("renta_desde = ?", filters.RentaDesde)
	}
	if filters.RentaHasta != 0 {
		query = query.Where("renta_hasta = ?", filters.RentaHasta)
	}
	if filters.ModalidadTrabajo != "" {
		query = query.Where("modalidad_trabajo LIKE ?", "%"+filters.ModalidadTrabajo+"%")
	}
	if filters.TipoServicio != "" {
		query = query.Where("tipo_servicio LIKE ?", "%"+filters.TipoServicio+"%")
	}
	if filters.NivelExperiencia != "" {
		query = query.Where("nivel_experiencia LIKE ?", "%"+filters.NivelExperiencia+"%")
	}
	if filters.FechaInicioProyecto != "" {
		query = query.Where("fecha_inicio_proyecto = ?", filters.FechaInicioProyecto)
	}

	//Paginacion
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Page > 0 {
		offset := (filters.Page - 1) * filters.Limit
		query = query.Offset(offset)
	}

	err := query.Find(&solicitudes).Error
	return solicitudes, err

}

func (r *repository) GetByID(ctx context.Context, id uint) (*Solicitud, error) {
	var solicitud Solicitud
	err := r.db.WithContext(ctx).First(&solicitud, id).Error
	if err != nil {
		return nil, err
	}
	return &solicitud, nil
}

func (r *repository) Update(ctx context.Context, id uint, req UpdateReq) error {
	updates := make(map[string]interface{})

	if req.Titulo != nil {
		updates["titulo"] = *req.Titulo
	}
	if req.Estado != nil {
		updates["estado"] = *req.Estado
	}
	if req.Area != nil {
		updates["area"] = *req.Area
	}
	if req.Pais != nil {
		updates["pais"] = *req.Pais
	}
	if req.NumeroVacantes != nil {
		updates["numero_vacantes"] = *req.NumeroVacantes
	}
	if req.RentaDesde != nil {
		updates["renta_desde"] = *req.RentaDesde
	}
	if req.RentaHasta != nil {
		updates["renta_hasta"] = *req.RentaHasta
	}
	if req.ModalidadTrabajo != nil {
		updates["modalidad_trabajo"] = *req.ModalidadTrabajo
	}
	if req.TipoServicio != nil {
		updates["tipo_servicio"] = *req.TipoServicio
	}
	if req.NivelExperiencia != nil {
		updates["nivel_experiencia"] = *req.NivelExperiencia
	}
	if req.FechaInicioProyecto != nil {
		updates["fecha_inicio_proyecto"] = *req.FechaInicioProyecto
	}
	return r.db.WithContext(ctx).Model(&Solicitud{}).Where("id = ?", id).Updates(updates).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Solicitud{}, id).Error
}
