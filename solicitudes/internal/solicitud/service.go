package solicitud

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Service interface {
	Create(ctx context.Context, req CreateReq) (*Solicitud, error)
	GetAll(ctx context.Context, filter GetAllReq) ([]SolicitudResponse, error)
	GetByID(ctx context.Context, id uint) (*SolicitudResponse, error)
	GetByIDWithDocuments(ctx context.Context, id uint) (*SolicitudResponse, error) // New method
	Update(ctx context.Context, id uint, req UpdateReq) error
	Delete(ctx context.Context, id uint) error
}

// DocumentoClient define la interfaz para el cliente de documentos
type DocumentoClient interface {
	GetBySolicitudID(solicitudID uint) ([]Documento, error)
}

type service struct {
	repo          Repository
	logger        *log.Logger
	documentoClient DocumentoClient
}

func NewService(repo Repository, logger *log.Logger, docClient DocumentoClient) Service {
	return &service{
		repo:          repo,
		logger:        logger,
		documentoClient: docClient,
	}
}

func (s *service) Create(ctx context.Context, req CreateReq) (*Solicitud, error) {
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

	if err := s.repo.Create(ctx, solicitud); err != nil {
		s.logger.Printf("Error al crear la solicitud: %v", err)
		return nil, err
	}

	s.logger.Printf("Solicitud creada exitosamente: ID=%d para Usuario ID=%d", solicitud.ID, solicitud.UsuarioID)
	return solicitud, nil
}

func (s *service) GetAll(ctx context.Context, filter GetAllReq) ([]SolicitudResponse, error) {
	solicitudes, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		s.logger.Printf("Error al obtener las solicitudes: %v", err)
		return nil, err
	}

	// Inicializar el slice de respuestas
	responses := make([]SolicitudResponse, len(solicitudes))

	for i, solicitud := range solicitudes {
		// Convertir a respuesta básica primero
		responses[i] = solicitud.ToResponse()
		
		// Obtener documentos del microservicio
		documentos, err := s.documentoClient.GetBySolicitudID(solicitud.ID)
		if err != nil {
			s.logger.Printf("Advertencia: No se pudieron obtener documentos para solicitud ID=%d: %v", solicitud.ID, err)
			responses[i].Documentos = []DocumentoResponse{}
		} else {
			// Mapear documentos a DocumentoResponse (solo información básica)
			docResponses := make([]DocumentoResponse, len(documentos))
			for j, doc := range documentos {
				docResponses[j] = DocumentoResponse{
					ID:            doc.ID,
					NombreArchivo: doc.NombreArchivo,
					Extension:     doc.Extension,
				}
			}
			responses[i].Documentos = docResponses
		}
	}

	s.logger.Printf("Se obtuvieron %d solicitudes", len(responses))
	return responses, nil

	s.logger.Printf("Se obtuvieron %d solicitudes", len(responses))
	return responses, nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*SolicitudResponse, error) {
	solicitud, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Printf("Error al obtener solicitud ID=%d: %v", id, err)
		return nil, fmt.Errorf("error al obtener la solicitud: %v", err)
	}

	// Devolver solo la información básica, sin documentos
	response := solicitud.ToResponse()
	response.Documentos = []DocumentoResponse{}

	return &response, nil
}

// GetByIDWithDocuments obtiene una solicitud con sus documentos asociados
func (s *service) GetByIDWithDocuments(ctx context.Context, id uint) (*SolicitudResponse, error) {
	solicitud, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Printf("Error al obtener solicitud ID=%d: %v", id, err)
		return nil, fmt.Errorf("error al obtener la solicitud: %v", err)
	}

	response := solicitud.ToResponse()
	
	// Obtener documentos del microservicio
	documentos, err := s.documentoClient.GetBySolicitudID(solicitud.ID)
	if err != nil {
		s.logger.Printf("Advertencia: No se pudieron obtener documentos para solicitud ID=%d: %v", solicitud.ID, err)
		response.Documentos = []DocumentoResponse{}
	} else {
		response.Documentos = make([]DocumentoResponse, len(documentos))
		for i, doc := range documentos {
			response.Documentos[i] = DocumentoResponse{
				ID:            doc.ID,
				NombreArchivo: doc.NombreArchivo,
				Extension:     doc.Extension,
			}
		}
	}

	return &response, nil
}

func (s *service) Update(ctx context.Context, id uint, req UpdateReq) error {
	if err := s.repo.Update(ctx, id, req); err != nil {
		s.logger.Printf("Error al actualizar la solicitud ID=%d: %v", id, err)
		return err
	}

	s.logger.Printf("Solicitud actualizada exitosamente: ID=%d", id)
	return nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Printf("Error al eliminar la solicitud ID=%d: %v", id, err)
		return err
	}

	s.logger.Printf("Solicitud eliminada exitosamente: ID=%d", id)
	return nil
}
