package usuario

import (
	"log"
)

type Service interface {
	Create(req CreateReq) (*Usuario, error)
	GetAll(filters GetAllReq) ([]Usuario, error)
	GetByID(id uint) (*Usuario, error)
	Update(id uint, req UpdateReq) error
	Delete(id uint) error
}

type service struct {
	repo   Repository
	logger *log.Logger
}

func NewService(repo Repository, logger *log.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) Create(req CreateReq) (*Usuario, error) {
	usuario := &Usuario{
		NombreUsuario: req.NombreUsuario,
		EmailUsuario:  req.EmailUsuario,
		Nombre:        req.Nombre,
		Apellidos:     req.Apellidos,
		Idioma:        req.Idioma,
		ZonaHoraria:   req.ZonaHoraria,
	}

	if err := s.repo.Create(usuario); err != nil {
		s.logger.Printf("Error al crear usuario: %v", err)
		return nil, err
	}

	s.logger.Printf("Usuario creado exitosamente: ID=%d", usuario.ID)
	return usuario, nil
}

func (s *service) GetAll(filters GetAllReq) ([]Usuario, error) {
	usuarios, err := s.repo.GetAll(filters)
	if err != nil {
		s.logger.Printf("Error al obtener usuarios: %v", err)
		return nil, err
	}

	s.logger.Printf("Se obtuvieron %d usuarios", len(usuarios))
	return usuarios, nil
}

func (s *service) GetByID(id uint) (*Usuario, error) {
	usuario, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Printf("Error al obtener usuario ID=%d: %v", id, err)
		return nil, err
	}

	return usuario, nil
}

func (s *service) Update(id uint, req UpdateReq) error {
	if err := s.repo.Update(id, req); err != nil {
		s.logger.Printf("Error al actualizar usuario ID=%d: %v", id, err)
		return err
	}

	s.logger.Printf("Usuario actualizado exitosamente: ID=%d", id)
	return nil
}

func (s *service) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		s.logger.Printf("Error al eliminar usuario ID=%d: %v", id, err)
		return err
	}

	s.logger.Printf("Usuario eliminado exitosamente: ID=%d", id)
	return nil
}
