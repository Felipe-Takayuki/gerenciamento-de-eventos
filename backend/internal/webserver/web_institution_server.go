package webserver

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	CNPJ     string `json:"cnpj"`
	IsAdmin  bool   `json:"is_admin,omitempty"`
}

type LoginInstitutionRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"token":       token,
		"institution": result,
	})
}

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
