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
	Name        string `json:"name"`
	Address     string `json:"address"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
}

type EditEventRequest struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
}

type EditRoomRequest struct {
	ID   int64  `json:"room_id"`
	Name string `json:"name"`
}

type DeleteRoomRequest struct {
	ID int64 `json:"room_id"`
}

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

func (h *WebEventHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	events, err := h.eventService.GetEvents(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}
	jsonResponse(w, http.StatusOK, events)
}

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
	jsonResponse(w, http.StatusOK, map[string]interface{}{"message": "sala excluída com sucesso"})
}

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
	jsonResponse(w, http.StatusOK, map[string]interface{}{"message": "evento excluído com sucesso"})
}
