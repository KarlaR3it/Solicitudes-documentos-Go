package documento

import (
	"context"
	"fmt"
	"log"

	"github.com/kramirez/documentos/pkg/httpclient"
)

type Service interface {
	Create(ctx context.Context, req CreateReq) (*DocumentoResponse, error)
	GetAll(ctx context.Context, filter GetAllReq) ([]DocumentoResponse, error)
	GetByID(ctx context.Context, id uint) (*DocumentoResponse, error)
	Update(ctx context.Context, id uint, req UpdateReq) error
	Delete(ctx context.Context, id uint) error
}

type service struct {
	repo            Repository
	logger          *log.Logger
	solicitudClient *httpclient.SolicitudClient
}

func NewService(repo Repository, logger *log.Logger, solicitudClient *httpclient.SolicitudClient) Service {
	return &service{
		repo:            repo,
		logger:          logger,
		solicitudClient: solicitudClient,
	}
}

func (s *service) Create(ctx context.Context, req CreateReq) (*DocumentoResponse, error) {
	// Validar que la solicitud existe
	solicitud, err := s.solicitudClient.GetSolicitud(req.SolicitudID)
	if err != nil {
		s.logger.Printf("Error al validar solicitud ID=%d: %v", req.SolicitudID, err)
		return nil, fmt.Errorf("error al validar solicitud: %v", err)
	}

	documento := &Documento{
		Extension:     req.Extension,
		NombreArchivo: req.NombreArchivo,
		SolicitudID:   req.SolicitudID,
	}

	if err := s.repo.Create(ctx, documento); err != nil {
		s.logger.Printf("Error al crear el documento: %v", err)
		return nil, err
	}

	// Crear la respuesta con la información de la solicitud
	response := s.toDocumentoResponse(documento, solicitud)

	s.logger.Printf("Documento creado exitosamente: ID=%d para Solicitud ID=%d", documento.ID, documento.SolicitudID)
	return &response, nil
}

func (s *service) toDocumentoResponse(doc *Documento, solicitud *httpclient.SolicitudResponse) DocumentoResponse {
	var response DocumentoResponse
	response.ID = doc.ID
	response.Extension = doc.Extension
	response.NombreArchivo = doc.NombreArchivo
	response.CreatedAt = doc.CreatedAt
	response.UpdatedAt = doc.UpdatedAt
	response.Solicitud.ID = solicitud.ID
	response.Solicitud.Titulo = solicitud.Titulo
	response.Solicitud.Area = solicitud.Area
	return response
}

func (s *service) GetAll(ctx context.Context, filter GetAllReq) ([]DocumentoResponse, error) {
	documentos, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		s.logger.Printf("Error al obtener los documentos: %v", err)
		return nil, err
	}

	// Crear un slice para las respuestas
	responses := make([]DocumentoResponse, 0, len(documentos))

	// Obtener los IDs de las solicitudes únicas
	solicitudIDs := make(map[uint]bool)
	for _, doc := range documentos {
		solicitudIDs[doc.SolicitudID] = true
	}

	// Obtener los detalles de las solicitudes
	solicitudes := make(map[uint]*httpclient.SolicitudResponse)
	for id := range solicitudIDs {
		solicitud, err := s.solicitudClient.GetSolicitud(id)
		if err != nil {
			s.logger.Printf("Advertencia: No se pudo obtener la solicitud ID=%d: %v", id, err)
			continue
		}
		solicitudes[id] = solicitud
	}

	// Construir las respuestas
	for _, doc := range documentos {
		if solicitud, ok := solicitudes[doc.SolicitudID]; ok {
			response := s.toDocumentoResponse(&doc, solicitud)
			responses = append(responses, response)
		} else {
			s.logger.Printf("Advertencia: No se encontró la solicitud ID=%d para el documento ID=%d", doc.SolicitudID, doc.ID)
		}
	}

	s.logger.Printf("Se obtuvieron %d documentos", len(responses))
	return responses, nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*DocumentoResponse, error) {
	documento, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Printf("Error al obtener el documento ID=%d: %v", id, err)
		return nil, err
	}

	// Obtener los detalles de la solicitud
	solicitud, err := s.solicitudClient.GetSolicitud(documento.SolicitudID)
	if err != nil {
		s.logger.Printf("Advertencia: No se pudo obtener la solicitud ID=%d: %v", documento.SolicitudID, err)
		return nil, fmt.Errorf("no se pudo obtener la información de la solicitud: %v", err)
	}

	// Crear la respuesta
	response := s.toDocumentoResponse(documento, solicitud)
	return &response, nil
}

func (s *service) Update(ctx context.Context, id uint, req UpdateReq) error {
	if err := s.repo.Update(ctx, id, req); err != nil {
		s.logger.Printf("Error al actualizar el documento ID=%d: %v", id, err)
		return err
	}

	s.logger.Printf("Documento actualizado exitosamente: ID=%d", id)
	return nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Printf("Error al eliminar el documento ID=%d: %v", id, err)
		return err
	}

	s.logger.Printf("Documento eliminado exitosamente: ID=%d", id)
	return nil
}
