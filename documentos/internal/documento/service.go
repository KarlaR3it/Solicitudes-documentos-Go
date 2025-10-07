package documento

import (
	"context"
	"fmt"
	"log"

	"github.com/kramirez/documentos/pkg/httpclient"
)

type Service interface {
	Create(ctx context.Context, req CreateReq) (*Documento, error)
	GetAll(ctx context.Context, filter GetAllReq) ([]Documento, error)
	GetByID(ctx context.Context, id uint) (*Documento, error)
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

func (s *service) Create(ctx context.Context, req CreateReq) (*Documento, error) {
	// Validar que la solicitud existe
	existe, err := s.solicitudClient.ValidarSolicitud(req.SolicitudID)
	if err != nil {
		s.logger.Printf("Error al validar solicitud ID=%d: %v", req.SolicitudID, err)
		return nil, fmt.Errorf("error al validar solicitud: %v", err)
	}
	if !existe {
		s.logger.Printf("Solicitud ID=%d no encontrado", req.SolicitudID)
		return nil, fmt.Errorf("la solicitud con ID %d no existe", req.SolicitudID)
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

	s.logger.Printf("Documento creado exitosamente: ID=%d para Solicitud ID=%d", documento.ID, documento.SolicitudID)
	return documento, nil
}

func (s *service) GetAll(ctx context.Context, filter GetAllReq) ([]Documento, error) {
	documentos, err := s.repo.GetAll(ctx, filter)
	if err != nil {
		s.logger.Printf("Error al obtener los documentos: %v", err)
		return nil, err
	}

	s.logger.Printf("Se obtuvieron %d documentos", len(documentos))
	return documentos, nil
}

func (s *service) GetByID(ctx context.Context, id uint) (*Documento, error) {
	documento, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Printf("Error al obtener el documento ID=%d: %v", id, err)
		return nil, err
	}

	return documento, nil
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
