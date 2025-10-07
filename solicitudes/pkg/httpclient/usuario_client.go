package httpclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UsuarioClient struct {
	baseURL    string
	httpClient *http.Client
}

type UsuarioResponse struct {
	ID            uint   `json:"id"`
	NombreUsuario string `json:"nombre_usuario"`
	EmailUsuario  string `json:"email_usuario"`
	Nombre        string `json:"nombre"`
	Apellidos     string `json:"apellidos"`
}

func NewUsuarioClient(baseURL string) *UsuarioClient {
	return &UsuarioClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ValidarUsuario verifica si un usuario existe en el servicio de usuarios
func (c *UsuarioClient) ValidarUsuario(usuarioID uint) (bool, error) {
	url := fmt.Sprintf("%s/usuarios/%d", c.baseURL, usuarioID)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return false, fmt.Errorf("error al conectar con el servicio de usuarios: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("error del servicio de usuarios: status %d", resp.StatusCode)
	}

	return true, nil
}

// ObtenerUsuario obtiene la información completa de un usuario
func (c *UsuarioClient) ObtenerUsuario(usuarioID uint) (*UsuarioResponse, error) {
	url := fmt.Sprintf("%s/usuarios/%d", c.baseURL, usuarioID)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al conectar con el servicio de usuarios: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error del servicio de usuarios: status %d", resp.StatusCode)
	}

	var usuario UsuarioResponse
	if err := json.NewDecoder(resp.Body).Decode(&usuario); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	return &usuario, nil
}
