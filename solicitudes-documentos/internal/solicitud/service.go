package solicitud

import (
	"fmt"
	"log"
	"time"
	"github.com/kramirez/solicitudes-documentos/pkg/httpclient"
)

type Service interface {
	Create(req CreateReq) (*Solicitud, error)
	GetAll(filter GetAllReq) ([]Solicitud, error)
	GetByID(id uint) (*Solicitud, error)
	Update(id uint, req UpdateReq) error
	Delete(id uint) error
}

type service struct {
	repo Repository
	logger *log.Logger
	usuarioClient *httpclient.UsuarioClient
}

func NewService(repo Repository, logger *log.Logger, usuarioClient *httpclient.UsuarioClient) Service {
	return &service{
		repo:   repo,
		logger: logger,
		usuarioClient: usuarioClient,
	}
}

func (s *service) Create(req CreateReq) (*Solicitud, error) {
	// Validar que el usuario existe
	existe, err := s.usuarioClient.ValidarUsuario(req.UsuarioID)
	if err != nil {
		s.logger.Printf("Error al validar usuario ID=%d: %v", req.UsuarioID, err)
		return nil, fmt.Errorf("error al validar usuario: %v", err)
	}
	if !existe {
		s.logger.Printf("Usuario ID=%d no encontrado", req.UsuarioID)
		return nil, fmt.Errorf("el usuario con ID %d no existe", req.UsuarioID)
	}

	// Parsear la fecha de string a time.Time
	fechaInicio, err := time.Parse("2006-01-02", req.FechaInicioProyecto)
	if err != nil {
		s.logger.Printf("Error al parsear la fecha: %v", err)
		return nil, fmt.Errorf("formato de fecha inválido, use YYYY-MM-DD: %v", err)
	}

	solicitud := &Solicitud{
		Titulo:                   req.Titulo,
		Estado:                   req.Estado,
		Area:                     req.Area,
		Pais:                     req.Pais,
		Localizacion:             req.Localizacion,
		NumeroVacantes:           req.NumeroVacantes,
		Descripcion:              req.Descripcion,
		BaseEducacional:          req.BaseEducacional,
		ConocimientosExcluyentes: req.ConocimientosExcluyentes,
		RentaDesde:               req.RentaDesde,
		RentaHasta:               req.RentaHasta,
		ModalidadTrabajo:         req.ModalidadTrabajo,
		TipoServicio:             req.TipoServicio,
		NivelExperiencia:         req.NivelExperiencia,
		FechaInicioProyecto:      fechaInicio,
		UsuarioID:                req.UsuarioID,
	}

	if err := s.repo.Create(solicitud); err != nil {
		s.logger.Printf("Error al crear la solicitud: %v", err)
		return nil, err
	}
	
	s.logger.Printf("Solicitud creada exitosamente: ID=%d para Usuario ID=%d", solicitud.ID, solicitud.UsuarioID)
	return solicitud, nil
}

func (s *service) GetAll(filter GetAllReq) ([]Solicitud, error) {
	solicitudes, err := s.repo.GetAll(filter)
	if err != nil {
		s.logger.Printf("Error al obtener las solicitudes: %v", err)
		return nil, err
	}
	
	s.logger.Printf("Se obtuvieron %d solicitudes", len(solicitudes))
	return solicitudes, nil
}

func (s *service) GetByID(id uint) (*Solicitud, error) {
	solicitud, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Printf("Error al obtener la solicitud ID=%d: %v", id, err)
		return nil, err
	}
	
	return solicitud, nil
}

func (s *service) Update(id uint, req UpdateReq) error {
	if err := s.repo.Update(id, req); err != nil {
		s.logger.Printf("Error al actualizar la solicitud ID=%d: %v", id, err)
		return err
	}
	
	s.logger.Printf("Solicitud actualizada exitosamente: ID=%d", id)
	return nil
}

func (s *service) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		s.logger.Printf("Error al eliminar la solicitud ID=%d: %v", id, err)
		return err
	}
	
	s.logger.Printf("Solicitud eliminada exitosamente: ID=%d", id)
	return nil
}
