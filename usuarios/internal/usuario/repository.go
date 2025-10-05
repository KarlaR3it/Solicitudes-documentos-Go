package usuario

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(usuario *Usuario) error
	GetAll(filters GetAllReq) ([]Usuario, error)
	GetByID(id uint) (*Usuario, error)
	Update(id uint, req UpdateReq) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(usuario *Usuario) error {
	return r.db.Create(usuario).Error
}

func (r *repository) GetAll(filters GetAllReq) ([]Usuario, error) {
	var usuarios []Usuario
	query := r.db.Model(&Usuario{})

	// Aplicar filtros
	if filters.NombreUsuario != "" {
		query = query.Where("nombre_usuario LIKE ?", "%"+filters.NombreUsuario+"%")
	}
	if filters.EmailUsuario != "" {
		query = query.Where("email_usuario LIKE ?", "%"+filters.EmailUsuario+"%")
	}
	if filters.Nombre != "" {
		query = query.Where("nombre LIKE ?", "%"+filters.Nombre+"%")
	}
	if filters.Apellidos != "" {
		query = query.Where("apellidos LIKE ?", "%"+filters.Apellidos+"%")
	}

	// Paginación
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Page > 0 {
		offset := (filters.Page - 1) * filters.Limit
		query = query.Offset(offset)
	}

	err := query.Find(&usuarios).Error
	return usuarios, err
}

func (r *repository) GetByID(id uint) (*Usuario, error) {
	var usuario Usuario
	err := r.db.First(&usuario, id).Error
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}

func (r *repository) Update(id uint, req UpdateReq) error {
	updates := make(map[string]interface{})

	if req.NombreUsuario != nil {
		updates["nombre_usuario"] = *req.NombreUsuario
	}
	if req.EmailUsuario != nil {
		updates["email_usuario"] = *req.EmailUsuario
	}
	if req.Nombre != nil {
		updates["nombre"] = *req.Nombre
	}
	if req.Apellidos != nil {
		updates["apellidos"] = *req.Apellidos
	}
	if req.Idioma != nil {
		updates["idioma"] = *req.Idioma
	}
	if req.ZonaHoraria != nil {
		updates["zona_horaria"] = *req.ZonaHoraria
	}

	return r.db.Model(&Usuario{}).Where("id = ?", id).Updates(updates).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Usuario{}, id).Error
}
