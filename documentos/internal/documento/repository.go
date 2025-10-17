package documento

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, documento *Documento) error
	GetAll(ctx context.Context, filters GetAllReq) ([]Documento, error)
	GetByID(ctx context.Context, id uint) (*Documento, error)
	Update(ctx context.Context, id uint, req UpdateReq) error
	Delete(ctx context.Context, id uint) error
	DeleteBySolicitudID(ctx context.Context, solicitudID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, documento *Documento) error {
	return r.db.WithContext(ctx).Create(documento).Error
}

func (r *repository) GetAll(ctx context.Context, filters GetAllReq) ([]Documento, error) {
	var documentos []Documento
	query := r.db.WithContext(ctx).Model(&Documento{})

	//Aplicar filtros
	if filters.Extension != "" {
		query = query.Where("extension LIKE ?", "%"+filters.Extension+"%")
	}
	if filters.NombreArchivo != "" {
		query = query.Where("nombre_archivo LIKE ?", "%"+filters.NombreArchivo+"%")
	}
	// NUEVO: Filtrar por solicitud_id
	if filters.SolicitudID > 0 {
		query = query.Where("solicitud_id = ?", filters.SolicitudID)
	}

	//Paginacion
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Page > 0 {
		offset := (filters.Page - 1) * filters.Limit
		query = query.Offset(offset)
	}

	err := query.Find(&documentos).Error
	return documentos, err
}

func (r *repository) GetByID(ctx context.Context, id uint) (*Documento, error) {
	var documento Documento
	err := r.db.WithContext(ctx).First(&documento, id).Error
	if err != nil {
		return nil, err
	}
	return &documento, nil
}

func (r *repository) Update(ctx context.Context, id uint, req UpdateReq) error {
	updates := make(map[string]interface{})

	if req.Extension != nil {
		updates["extension"] = *req.Extension
	}
	if req.NombreArchivo != nil {
		updates["nombre_archivo"] = *req.NombreArchivo
	}
	return r.db.WithContext(ctx).Model(&Documento{}).Where("id = ?", id).Updates(updates).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	// Soft delete: GORM automáticamente establece deleted_at
	return r.db.WithContext(ctx).Delete(&Documento{}, id).Error
}

func (r *repository) DeleteBySolicitudID(ctx context.Context, solicitudID uint) error {
	// Soft delete de todos los documentos asociados a una solicitud
	return r.db.WithContext(ctx).Where("solicitud_id = ?", solicitudID).Delete(&Documento{}).Error
}
