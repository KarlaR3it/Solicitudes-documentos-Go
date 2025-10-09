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
	repo            Repository
	logger          *log.Logger
	documentoClient DocumentoClient
}

func NewService(repo Repository, logger *log.Logger, docClient DocumentoClient) Service {
	return &service{
		repo:            repo,
		logger:          logger,
		documentoClient: docClient,
	}
}

// validateCreateRequest valida los campos requeridos de la solicitud
func validateCreateRequest(req CreateReq) error {
	if req.Titulo == "" {
		return fmt.Errorf("el título es requerido")
	}
	if req.Estado == "" {
		req.Estado = "pendiente" // Valor por defecto
	}
	if req.Area == "" {
		return fmt.Errorf("el área es requerida")
	}
	if req.Pais == "" {
		return fmt.Errorf("el país es requerido")
	}
	if req.Localizacion == "" {
		return fmt.Errorf("la localización es requerida")
	}
	if req.FechaInicioProyecto == "" {
		return fmt.Errorf("la fecha de inicio del proyecto es requerida")
	}
	if req.UsuarioID == nil {
		return fmt.Errorf("el ID de usuario es requerido")
	}

	// Validar formato de fecha
	_, err := time.Parse("2006-01-02", req.FechaInicioProyecto)
	if err != nil {
		return fmt.Errorf("formato de fecha inválido, debe ser YYYY-MM-DD")
	}

	// Validar rango de renta
	if req.RentaDesde > 0 && req.RentaHasta > 0 && req.RentaDesde > req.RentaHasta {
		return fmt.Errorf("el rango de renta es inválido")
	}

	return nil
}

func (s *service) Create(ctx context.Context, req CreateReq) (*Solicitud, error) {
	// Validar campos requeridos
	if err := validateCreateRequest(req); err != nil {
		s.logger.Printf("Validación fallida: %v", err)
		return nil, err
	}

	// Parsear la fecha de string a time.Time
	fechaInicio, err := time.Parse("2006-01-02", req.FechaInicioProyecto)
	if err != nil {
		s.logger.Printf("Error al parsear la fecha: %v", err)
		return nil, fmt.Errorf("formato de fecha inválido, use YYYY-MM-DD: %v", err)
	}

	// Establecer valor por defecto para Estado si está vacío
	estado := req.Estado
	if estado == "" {
		estado = "pendiente"
	}

	solicitud := &Solicitud{
		Titulo:                   req.Titulo,
		Estado:                   estado,
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
	// Verificar que la solicitud exista antes de actualizar
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Printf("Error al buscar solicitud ID=%d: %v", id, err)
		return fmt.Errorf("solicitud no encontrada")
	}

	if err := s.repo.Update(ctx, id, req); err != nil {
		s.logger.Printf("Error al actualizar la solicitud ID=%d: %v", id, err)
		return err
	}

	s.logger.Printf("Solicitud actualizada exitosamente: ID=%d", id)
	return nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	// Verificar que la solicitud exista antes de eliminar
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Printf("Error al buscar solicitud ID=%d: %v", id, err)
		return fmt.Errorf("solicitud no encontrada")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Printf("Error al eliminar la solicitud ID=%d: %v", id, err)
		return err
	}

	s.logger.Printf("Solicitud eliminada exitosamente: ID=%d", id)
	return nil
}
