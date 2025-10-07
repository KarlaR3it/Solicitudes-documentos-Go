package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SolicitudClient struct {
	baseURL    string
	httpClient *http.Client
}

type SolicitudResponse struct {
	ID            uint   `json:"id"`
	Extension     string `json:"extension"`
	NombreArchivo string `json:"nombre_archivo"`
}

func NewSolicitudClient(baseURL string) *SolicitudClient {
	return &SolicitudClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ValidarUsuario verifica si un usuario existe en el servicio de usuarios
func (c *SolicitudClient) ValidarSolicitud(solicitudID uint) (bool, error) {
	url := fmt.Sprintf("%s/solicitudes/%d", c.baseURL, solicitudID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return false, fmt.Errorf("error al conectar con el servicio de solicitudes: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("error del servicio de solicitudes: status %d", resp.StatusCode)
	}

	return true, nil
}

// ObtenerSolicitud obtiene la información completa de una solicitud
func (c *SolicitudClient) ObtenerSolicitud(solicitudID uint) (*SolicitudResponse, error) {
	url := fmt.Sprintf("%s/solicitudes/%d", c.baseURL, solicitudID)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al conectar con el servicio de solicitudes: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("solicitud no encontrada")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error del servicio de solicitudes: status %d", resp.StatusCode)
	}

	var solicitud SolicitudResponse
	if err := json.NewDecoder(resp.Body).Decode(&solicitud); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	return &solicitud, nil
}
