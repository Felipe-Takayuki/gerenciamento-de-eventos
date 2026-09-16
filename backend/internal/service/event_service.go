package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/entity"
	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/utils"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, name, address, startDate, endDate, description string, institutionID int64) (*entity.Event, error)
	DeleteEvent(ctx context.Context, eventID, institutionID int64) error
	EditEvent(ctx context.Context, eventID, ownerID int64, name, address, startDate, endDate, description string) (*entity.Event, error)
	GetEventByID(ctx context.Context, eventID int64) (*entity.Event, error)
	GetEvents(ctx context.Context) ([]*entity.Event, error)
	GetEventByName(ctx context.Context, name string) ([]*entity.Event, error)
	GetEventByOwnerID(ctx context.Context, ownerID int64) ([]*entity.Event, error)
	AddRoomInEvent(ctx context.Context, eventID, ownerID int64, roomName string) ([]*entity.RoomEvent, error)
	GetRoomsByEventID(ctx context.Context, eventID, ownerID int64) ([]*entity.RoomEvent, error)
	EditRoom(ctx context.Context, roomID, eventID, ownerID int64, roomName string) (*entity.RoomEvent, error)
	DeleteRoom(ctx context.Context, roomID, eventID, ownerID int64) error
}

type EventService struct {
	eventRepo EventRepository
}

func NewEventService(eventRepo EventRepository) *EventService {
	return &EventService{
		eventRepo: eventRepo,
	}
}

func (es *EventService) CreateEvent(ctx context.Context, name, address, startDate, endDate, description string, institutionID int64) (*entity.Event, error) {
	name = strings.TrimSpace(name)
	address = strings.TrimSpace(address)
	startDate = strings.TrimSpace(startDate)
	endDate = strings.TrimSpace(endDate)
	description = strings.TrimSpace(description)

	if name == "" || address == "" || startDate == "" || endDate == "" {
		return nil, errors.New("nome, endereço, data de início e data de término são obrigatórios")
	}
	if institutionID <= 0 {
		return nil, utils.ErrBadRequest
	}

	return es.eventRepo.CreateEvent(ctx, name, address, startDate, endDate, description, institutionID)
}

func (es *EventService) GetEventByName(ctx context.Context, name string) ([]*entity.Event, error) {
	return es.eventRepo.GetEventByName(ctx, strings.TrimSpace(name))
}

func (es *EventService) GetEventByID(ctx context.Context, eventID int64) (*entity.Event, error) {
	if eventID <= 0 {
		return nil, utils.ErrBadRequest
	}
	return es.eventRepo.GetEventByID(ctx, eventID)
}

func (es *EventService) GetEvents(ctx context.Context) ([]*entity.Event, error) {
	return es.eventRepo.GetEvents(ctx)
}

func (es *EventService) GetEventByOwnerID(ctx context.Context, ownerID int64) ([]*entity.Event, error) {
	if ownerID <= 0 {
		return nil, utils.ErrBadRequest
	}
	return es.eventRepo.GetEventByOwnerID(ctx, ownerID)
}

func (es *EventService) GetRoomsByEventID(ctx context.Context, eventID, ownerID int64) ([]*entity.RoomEvent, error) {
	if eventID <= 0 || ownerID <= 0 {
		return nil, utils.ErrBadRequest
	}
	return es.eventRepo.GetRoomsByEventID(ctx, eventID, ownerID)
}

func (es *EventService) DeleteEvent(ctx context.Context, eventID, institutionID int64) error {
	if eventID <= 0 || institutionID <= 0 {
		return utils.ErrBadRequest
	}
	return es.eventRepo.DeleteEvent(ctx, eventID, institutionID)
}

func (es *EventService) DeleteRoom(ctx context.Context, roomID, eventID, ownerID int64) error {
	if roomID <= 0 || eventID <= 0 || ownerID <= 0 {
		return utils.ErrBadRequest
	}
	return es.eventRepo.DeleteRoom(ctx, roomID, eventID, ownerID)
}

func (es *EventService) AddRoomInEvent(ctx context.Context, eventID, ownerID int64, roomName string) ([]*entity.RoomEvent, error) {
	roomName = strings.TrimSpace(roomName)
	if roomName == "" {
		return nil, errors.New("o nome da sala é obrigatório")
	}
	if eventID <= 0 || ownerID <= 0 {
		return nil, utils.ErrBadRequest
	}
	return es.eventRepo.AddRoomInEvent(ctx, eventID, ownerID, roomName)
}

func (es *EventService) EditEvent(ctx context.Context, eventID, ownerID int64, name, address, startDate, endDate, description string) (*entity.Event, error) {
	if eventID <= 0 || ownerID <= 0 {
		return nil, utils.ErrBadRequest
	}
	return es.eventRepo.EditEvent(ctx, eventID, ownerID, strings.TrimSpace(name), strings.TrimSpace(address), strings.TrimSpace(startDate), strings.TrimSpace(endDate), strings.TrimSpace(description))
}

func (es *EventService) EditRoom(ctx context.Context, roomID, eventID, ownerID int64, roomName string) (*entity.RoomEvent, error) {
	roomName = strings.TrimSpace(roomName)
	if roomName == "" {
		return nil, errors.New("o nome da sala é obrigatório")
	}
	if roomID <= 0 || eventID <= 0 || ownerID <= 0 {
		return nil, utils.ErrBadRequest
	}
	return es.eventRepo.EditRoom(ctx, roomID, eventID, ownerID, roomName)
}