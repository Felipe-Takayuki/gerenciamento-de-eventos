package webserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/entity"
	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

type WebInstitutionHandler struct {
	institutionService *service.InstitutionService
}

func NewWebInstitutionHandler(institutionService *service.InstitutionService) *WebInstitutionHandler {
	return &WebInstitutionHandler{
		institutionService: institutionService,
	}
}

type CreateInstitutionRequest struct {
	Name     string `json:"name" example:"Universidade Federal"`
	Email    string `json:"email" example:"contato@universidade.edu.br"`
	Password string `json:"password" example:"senhaSegura123"`
	CNPJ     string `json:"cnpj" example:"12345678000199"`
	IsAdmin  bool   `json:"is_admin,omitempty" example:"false"`
}

type LoginInstitutionRequest struct {
	Email    string `json:"email" example:"admin@adamas.com"`
	Password string `json:"password" example:"admin123"`
}

type LoginResponse struct {
	Token       string              `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Institution *entity.Institution `json:"institution"`
}

// CreateInstitution godoc
// @Summary      Criar uma nova instituição
// @Description  Cadastra uma nova instituição no sistema. Apenas instituições administradoras têm permissão para esta ação.
// @Tags         institutions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateInstitutionRequest true "Dados da nova instituição"
// @Success      201 {object} entity.Institution
// @Failure      400 {object} utils.ErrorMessage
// @Failure      401 {object} utils.ErrorMessage
// @Failure      403 {object} utils.ErrorMessage
// @Router       /institution [post]
func (h *WebInstitutionHandler) CreateInstitution(w http.ResponseWriter, r *http.Request) {
	// Apenas instituição do tipo admin pode criar novas instituições
	_, err := requireAdmin(r)
	if err != nil {
		handleError(w, err)
		return
	}

	var req CreateInstitutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "formato de requisição inválido")
		return
	}

	result, err := h.institutionService.CreateInstitution(r.Context(), req.Name, req.Email, req.Password, req.CNPJ, req.IsAdmin)
	if err != nil {
		handleError(w, err)
		return
	}

	jsonResponse(w, http.StatusCreated, result)
}

// LoginInstitution godoc
// @Summary      Autenticar instituição
// @Description  Realiza login com email e senha retornando o token JWT.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginInstitutionRequest true "Credenciais de acesso"
// @Success      200 {object} LoginResponse
// @Failure      400 {object} utils.ErrorMessage
// @Failure      401 {object} utils.ErrorMessage
// @Router       /login [post]
func (h *WebInstitutionHandler) LoginInstitution(w http.ResponseWriter, r *http.Request, tokenAuth *jwtauth.JWTAuth) {
	var req LoginInstitutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "formato de requisição inválido")
		return
	}

	result, err := h.institutionService.LoginInstitution(r.Context(), req.Email, req.Password)
	if err != nil {
		handleError(w, err)
		return
	}

	token, err := generateToken(tokenAuth, result.ID, result.Name, result.Email, result.UserType, result.IsAdmin)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "erro ao gerar token de autenticação")
		return
	}

	jsonResponse(w, http.StatusOK, LoginResponse{
		Token:       token,
		Institution: result,
	})
}

// GetInstitutionByID godoc
// @Summary      Obter instituição por ID
// @Description  Retorna os dados públicos de uma instituição cadastrada.
// @Tags         institutions
// @Produce      json
// @Param        institution_id path int true "ID da instituição"
// @Success      200 {object} entity.Institution
// @Failure      400 {object} utils.ErrorMessage
// @Failure      404 {object} utils.ErrorMessage
// @Router       /institution/{institution_id} [get]
func (h *WebInstitutionHandler) GetInstitutionByID(w http.ResponseWriter, r *http.Request) {
	institutionID, err := strconv.ParseInt(chi.URLParam(r, "institution_id"), 10, 64)
	if err != nil || institutionID <= 0 {
		errorResponse(w, http.StatusBadRequest, "institution_id deve ser um número inteiro válido")
		return
	}

	institution, err := h.institutionService.GetInstitutionByID(r.Context(), institutionID)
	if err != nil {
		handleError(w, err)
		return
	}

	jsonResponse(w, http.StatusOK, institution)
}
