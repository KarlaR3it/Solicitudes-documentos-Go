package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kramirez/solicitudes/internal/solicitud"
)

// DocumentoDTO represents a document from the documents microservice
type DocumentoDTO struct {
	ID            uint   `json:"id"`
	Extension     string `json:"extension"`
	NombreArchivo string `json:"nombre_archivo"`
	SolicitudID   uint   `json:"solicitud_id"`
}

// DocumentoClient implementa la interfaz solicitud.DocumentoClient
type DocumentoClient struct {
	baseURL string
	client  *http.Client
}

// toSolicitudDocumento convierte un DocumentoDTO a un Documento de solicitud
func (d *DocumentoDTO) toSolicitudDocumento() solicitud.Documento {
	return solicitud.Documento{
		ID:            d.ID,
		NombreArchivo: d.NombreArchivo,
		Extension:     d.Extension,
	}
}

func NewDocumentoClient(baseURL string) *DocumentoClient {
	return &DocumentoClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetBySolicitudID obtiene los documentos de una solicitud
func (c *DocumentoClient) GetBySolicitudID(solicitudID uint) ([]solicitud.Documento, error) {
	var documentosDTO []DocumentoDTO
	var documentos []solicitud.Documento
	
	// Construir la URL para obtener los documentos de la solicitud
	url := fmt.Sprintf("%s/documentos?solicitud_id=%d", c.baseURL, solicitudID)
	
	// Realizar la petición HTTP
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear la petición: %v", err)
	}

	// Configurar headers si es necesario
	req.Header.Set("Content-Type", "application/json")

	// Realizar la petición
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al conectar con el servicio de documentos: %v", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error al obtener documentos: status %d", resp.StatusCode)
	}

	// Decodificar la respuesta
	if err := json.NewDecoder(resp.Body).Decode(&documentosDTO); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	// Convertir DTOs a modelos de dominio
	for _, dto := range documentosDTO {
		documentos = append(documentos, dto.toSolicitudDocumento())
	}

	return documentos, nil
}

// DeleteBySolicitudID elimina (soft delete) todos los documentos asociados a una solicitud
func (c *DocumentoClient) DeleteBySolicitudID(solicitudID uint) error {
	// Construir la URL para eliminar los documentos de la solicitud
	url := fmt.Sprintf("%s/documentos/solicitud/%d", c.baseURL, solicitudID)
	
	// Realizar la petición HTTP DELETE
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error al crear la petición: %v", err)
	}

	// Configurar headers
	req.Header.Set("Content-Type", "application/json")

	// Realizar la petición
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("error al conectar con el servicio de documentos: %v", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error al eliminar documentos: status %d", resp.StatusCode)
	}

	return nil
}