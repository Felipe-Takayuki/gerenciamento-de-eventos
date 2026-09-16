package webserver

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/entity/reqs"
	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/service"
	"github.com/go-chi/chi"
)

type WebEventHandler struct {
	eventService *service.EventService
}

func NewWebEventHandler(eventService *service.EventService) *WebEventHandler {
	return &WebEventHandler{
		eventService: eventService,
	}
}

type CreateEventRequest struct {
	Name        string `json:"name" example:"Semana da Tecnologia 2026"`
	Address     string `json:"address" example:"Av. Paulista, 1000, São Paulo - SP"`
	StartDate   string `json:"start_date" example:"2026-10-10 09:00:00"`
	EndDate     string `json:"end_date" example:"2026-10-12 18:00:00"`
	Description string `json:"description" example:"Evento anual de tecnologia e inovação"`
}

type EditEventRequest struct {
	Name        string `json:"name" example:"Semana da Tecnologia 2026 - Edição Especial"`
	Address     string `json:"address" example:"Av. Paulista, 1500, São Paulo - SP"`
	StartDate   string `json:"start_date" example:"2026-10-10 09:00:00"`
	EndDate     string `json:"end_date" example:"2026-10-12 18:00:00"`
	Description string `json:"description" example:"Programação atualizada com palestras internacionais"`
}

type EditRoomRequest struct {
	ID   int64  `json:"room_id" example:"1"`
	Name string `json:"name" example:"Auditório Alpha"`
}

type DeleteRoomRequest struct {
	ID int64 `json:"room_id" example:"1"`
}

// GetEventByName godoc
// @Summary      Buscar eventos por nome
// @Description  Retorna lista de eventos públicos que correspondam ao termo de busca.
// @Tags         events
// @Produce      json
// @Param        event path string true "Nome ou termo de busca do evento"
// @Success      200 {array} entity.Event
// @Failure      400 {object} utils.ErrorMessage
// @Router       /event/search/{event} [get]
func (h *WebEventHandler) GetEventByName(w http.ResponseWriter, r *http.Request) {
	eventName := chi.URLParam(r, "event")
	if eventName == "" {
		errorResponse(w, http.StatusBadRequest, "o nome do evento é obrigatório")
		return
	}

	events, err := h.eventService.GetEventByName(r.Context(), eventName)
	if err != nil {
		handleError(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, events)
}

// GetEventByID godoc
// @Summary      Obter detalhes de um evento
// @Description  Retorna as informações completas de um evento e suas salas pelo seu ID.
// @Tags         events
// @Produce      json
// @Param        event_id path int true "ID do evento"
// @Success      200 {object} entity.Event
// @Failure      400 {object} utils.ErrorMessage
// @Failure      404 {object} utils.ErrorMessage
// @Router       /event/{event_id} [get]
func (h *WebEventHandler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	eventID, err := strconv.ParseInt(chi.URLParam(r, "event_id"), 10, 64)
	if err != nil || eventID <= 0 {
		errorResponse(w, http.StatusBadRequest, "event_id deve ser um número inteiro válido")
		return
	}

	event, err := h.eventService.GetEventByID(r.Context(), eventID)
	if err != nil {
		handleError(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, event)
}

// GetEvents godoc
// @Summary      Listar todos os eventos
// @Description  Retorna a lista pública de todos os eventos cadastrados no sistema.
// @Tags         events
// @Produce      json
// @Success      200 {array} entity.Event
// @Router       /events [get]
func (h *WebEventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.eventService.GetEvents(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, events)
}

// GetEventByOwnerID godoc
// @Summary      Listar eventos por instituição
// @Description  Retorna todos os eventos públicos de uma instituição específica.
// @Tags         events
// @Produce      json
// @Param        institution_id path int true "ID da instituição proprietária"
// @Success      200 {array} entity.Event
// @Failure      400 {object} utils.ErrorMessage
// @Router       /event/institution/{institution_id} [get]
func (h *WebEventHandler) GetEventByOwnerID(w http.ResponseWriter, r *http.Request) {
	institutionID, err := strconv.ParseInt(chi.URLParam(r, "institution_id"), 10, 64)
	if err != nil || institutionID <= 0 {
		errorResponse(w, http.StatusBadRequest, "institution_id deve ser um número inteiro válido")
		return
	}

	events, err := h.eventService.GetEventByOwnerID(r.Context(), institutionID)
	if err != nil {
		handleError(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, events)
}

// CreateEvent godoc
// @Summary      Criar um novo evento
// @Description  Cadastra um novo evento vinculado à instituição autenticada.
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateEventRequest true "Dados do evento"
// @Success      201 {object} entity.Event
// @Failure      400 {object} utils.ErrorMessage
// @Failure      401 {object} utils.ErrorMessage
// @Failure      403 {object} utils.ErrorMessage
// @Router       /event [post]
func (h *WebEventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	authUser, err := requireInstitution(r)
	if err != nil {
		handleError(w, err)
		return
	}

	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "formato de requisição inválido")
		return
	}

	event, err := h.eventService.CreateEvent(r.Context(), req.Name, req.Address, req.StartDate, req.EndDate, req.Description, authUser.ID)
	if err != nil {
		handleError(w, err)
		return
	}
	jsonResponse(w, http.StatusCreated, event)
}

// AddRoomInEvent godoc
// @Summary      Adicionar sala ao evento
// @Description  Adiciona um espaço/sala ao evento especificado. Requer ser instituição proprietária ou admin.
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        event_id path int true "ID do evento"
// @Param        request body reqs.AddRoomRequest true "Nome da sala"
// @Success      201 {array} entity.RoomEvent
// @Failure      400 {object} utils.ErrorMessage
// @Failure      401 {object} utils.ErrorMessage
// @Failure      403 {object} utils.ErrorMessage
// @Router       /event/{event_id}/room [post]
func (h *WebEventHandler) AddRoomInEvent(w http.ResponseWriter, r *http.Request) {
	authUser, err := requireInstitution(r)
	if err != nil {
		handleError(w, err)
		return
	}

	eventID, err := strconv.ParseInt(chi.URLParam(r, "event_id"), 10, 64)
	if err != nil || eventID <= 0 {
		errorResponse(w, http.StatusBadRequest, "event_id deve ser um número inteiro válido")
		return
	}

	var req reqs.AddRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "formato de requisição inválido")
		return
	}

	rooms, err := h.eventService.AddRoomInEvent(r.Context(), eventID, authUser.ID, req.Name)
	if err != nil {
		handleError(w, err)
		return
	}
	jsonResponse(w, http.StatusCreated, rooms)
}

// GetRoomsByEventID godoc
// @Summary      Listar salas do evento
// @Description  Retorna as salas de um evento. Requer ser instituição proprietária ou admin.
// @Tags         rooms
// @Produce      json
// @Security     BearerAuth
// @Param        event_id path int true "ID do evento"
// @Success      200 {array} entity.RoomEvent
// @Failure      400 {object} utils.ErrorMessage
// @Failure      401 {object} utils.ErrorMessage
// @Failure      403 {object} utils.ErrorMessage
// @Router       /event/{event_id}/room [get]
func (h *WebEventHandler) GetRoomsByEventID(w http.ResponseWriter, r *http.Request) {
	authUser, err := requireInstitution(r)
	if err != nil {
		handleError(w, err)
		return
	}

	eventID, err := strconv.ParseInt(chi.URLParam(r, "event_id"), 10, 64)
	if err != nil || eventID <= 0 {
		errorResponse(w, http.StatusBadRequest, "event_id deve ser um número inteiro válido")
		return
	}

	rooms, err := h.eventService.GetRoomsByEventID(r.Context(), eventID, authUser.ID)
	if err != nil {
		handleError(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, rooms)
}

// DeleteRoom godoc
// @Summary      Excluir sala do evento
// @Description  Exclui uma sala de um evento. Requer ser instituição proprietária ou admin.
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        event_id path int true "ID do evento"
// @Param        request body DeleteRoomRequest true "ID da sala"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} utils.ErrorMessage
// @Failure      401 {object} utils.ErrorMessage
// @Failure      403 {object} utils.ErrorMessage
// @Router       /event/{event_id}/room [delete]
func (h *WebEventHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	authUser, err := requireInstitution(r)
	if err != nil {
		handleError(w, err)
		return
	}

	eventID, err := strconv.ParseInt(chi.URLParam(r, "event_id"), 10, 64)
	if err != nil || eventID <= 0 {
		errorResponse(w, http.StatusBadRequest, "event_id deve ser um número inteiro válido")
		return
	}

	var req DeleteRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "formato de requisição inválido")
		return
	}

	if err := h.eventService.DeleteRoom(r.Context(), req.ID, eventID, authUser.ID); err != nil {
		handleError(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, MessageResponse{Message: "sala excluída com sucesso"})
}

// EditRoom godoc
// @Summary      Editar sala do evento
// @Description  Atualiza o nome de uma sala do evento. Requer ser instituição proprietária ou admin.
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        event_id path int true "ID do evento"
// @Param        request body EditRoomRequest true "Dados da sala"
// @Success      200 {object} entity.RoomEvent
// @Failure      400 {object} utils.ErrorMessage
// @Failure      401 {object} utils.ErrorMessage
// @Failure      403 {object} utils.ErrorMessage
// @Router       /event/{event_id}/room [put]
func (h *WebEventHandler) EditRoom(w http.ResponseWriter, r *http.Request) {
	authUser, err := requireInstitution(r)
	if err != nil {
		handleError(w, err)
		return
	}

	eventID, err := strconv.ParseInt(chi.URLParam(r, "event_id"), 10, 64)
	if err != nil || eventID <= 0 {
		errorResponse(w, http.StatusBadRequest, "event_id deve ser um número inteiro válido")
		return
	}

	var req EditRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "formato de requisição inválido")
		return
	}

	room, err := h.eventService.EditRoom(r.Context(), req.ID, eventID, authUser.ID, req.Name)
	if err != nil {
		handleError(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, room)
}

// EditEvent godoc
// @Summary      Editar evento
// @Description  Atualiza as informações de um evento existente. Requer ser instituição proprietária ou admin.
// @Tags         events
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        event_id path int true "ID do evento"
// @Param        request body EditEventRequest true "Dados atualizados do evento"
// @Success      200 {object} entity.Event
// @Failure      400 {object} utils.ErrorMessage
// @Failure      401 {object} utils.ErrorMessage
// @Failure      403 {object} utils.ErrorMessage
// @Router       /event/{event_id} [put]
func (h *WebEventHandler) EditEvent(w http.ResponseWriter, r *http.Request) {
	authUser, err := requireInstitution(r)
	if err != nil {
		handleError(w, err)
		return
	}

	eventID, err := strconv.ParseInt(chi.URLParam(r, "event_id"), 10, 64)
	if err != nil || eventID <= 0 {
		errorResponse(w, http.StatusBadRequest, "event_id deve ser um número inteiro válido")
		return
	}

	var req EditEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "formato de requisição inválido")
		return
	}

	event, err := h.eventService.EditEvent(r.Context(), eventID, authUser.ID, req.Name, req.Address, req.StartDate, req.EndDate, req.Description)
	if err != nil {
		handleError(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, event)
}

// DeleteEvent godoc
// @Summary      Excluir evento
// @Description  Exclui um evento e suas salas associadas. Requer ser instituição proprietária ou admin.
// @Tags         events
// @Produce      json
// @Security     BearerAuth
// @Param        event_id path int true "ID do evento"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} utils.ErrorMessage
// @Failure      401 {object} utils.ErrorMessage
// @Failure      403 {object} utils.ErrorMessage
// @Router       /event/{event_id} [delete]
func (h *WebEventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	authUser, err := requireInstitution(r)
	if err != nil {
		handleError(w, err)
		return
	}

	eventID, err := strconv.ParseInt(chi.URLParam(r, "event_id"), 10, 64)
	if err != nil || eventID <= 0 {
		errorResponse(w, http.StatusBadRequest, "event_id deve ser um número inteiro válido")
		return
	}

	if err := h.eventService.DeleteEvent(r.Context(), eventID, authUser.ID); err != nil {
		handleError(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, MessageResponse{Message: "evento excluído com sucesso"})
}
