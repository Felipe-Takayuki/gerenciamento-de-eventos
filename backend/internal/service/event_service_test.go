package service

import (
	"context"
	"testing"

	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/entity"
	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/utils"
)

type mockEventRepo struct {
	events map[int64]*entity.Event
	rooms  map[int64][]*entity.RoomEvent
}

func newMockEventRepo() *mockEventRepo {
	return &mockEventRepo{
		events: make(map[int64]*entity.Event),
		rooms:  make(map[int64][]*entity.RoomEvent),
	}
}

func (m *mockEventRepo) CreateEvent(ctx context.Context, name, address, startDate, endDate, description string, institutionID int64) (*entity.Event, error) {
	ev := &entity.Event{
		ID:            int64(len(m.events) + 1),
		Name:          name,
		Address:       address,
		StartDate:     startDate,
		EndDate:       endDate,
		Description:   description,
		InstitutionID: institutionID,
	}
	m.events[ev.ID] = ev
	return ev, nil
}

func (m *mockEventRepo) DeleteEvent(ctx context.Context, eventID, institutionID int64) error {
	delete(m.events, eventID)
	delete(m.rooms, eventID)
	return nil
}

func (m *mockEventRepo) EditEvent(ctx context.Context, eventID, ownerID int64, name, address, startDate, endDate, description string) (*entity.Event, error) {
	ev, ok := m.events[eventID]
	if !ok {
		return nil, utils.ErrNotFound
	}
	if name != "" {
		ev.Name = name
	}
	return ev, nil
}

func (m *mockEventRepo) GetEventByID(ctx context.Context, eventID int64) (*entity.Event, error) {
	ev, ok := m.events[eventID]
	if !ok {
		return nil, utils.ErrNotFound
	}
	return ev, nil
}

func (m *mockEventRepo) GetEvents(ctx context.Context) ([]*entity.Event, error) {
	var list []*entity.Event
	for _, ev := range m.events {
		list = append(list, ev)
	}
	return list, nil
}

func (m *mockEventRepo) GetEventByName(ctx context.Context, name string) ([]*entity.Event, error) {
	var list []*entity.Event
	for _, ev := range m.events {
		list = append(list, ev)
	}
	return list, nil
}

func (m *mockEventRepo) GetEventByOwnerID(ctx context.Context, ownerID int64) ([]*entity.Event, error) {
	var list []*entity.Event
	for _, ev := range m.events {
		if ev.InstitutionID == ownerID {
			list = append(list, ev)
		}
	}
	return list, nil
}

func (m *mockEventRepo) AddRoomInEvent(ctx context.Context, eventID, ownerID int64, roomName string) ([]*entity.RoomEvent, error) {
	r := &entity.RoomEvent{
		ID:   int64(len(m.rooms[eventID]) + 1),
		Name: roomName,
	}
	m.rooms[eventID] = append(m.rooms[eventID], r)
	return m.rooms[eventID], nil
}

func (m *mockEventRepo) GetRoomsByEventID(ctx context.Context, eventID, ownerID int64) ([]*entity.RoomEvent, error) {
	return m.rooms[eventID], nil
}

func (m *mockEventRepo) EditRoom(ctx context.Context, roomID, eventID, ownerID int64, roomName string) (*entity.RoomEvent, error) {
	return &entity.RoomEvent{ID: roomID, Name: roomName}, nil
}

func (m *mockEventRepo) DeleteRoom(ctx context.Context, roomID, eventID, ownerID int64) error {
	return nil
}

func TestEventService_CreateAndManage(t *testing.T) {
	repo := newMockEventRepo()
	svc := NewEventService(repo)
	ctx := context.Background()

	// Missing fields
	_, err := svc.CreateEvent(ctx, "", "Rua A", "2026-10-01", "2026-10-02", "Desc", 1)
	if err == nil {
		t.Fatal("esperava erro ao omitir nome")
	}

	// Invalid institution ID
	_, err = svc.CreateEvent(ctx, "Tech Fest", "Rua A", "2026-10-01", "2026-10-02", "Desc", 0)
	if err == nil {
		t.Fatal("esperava erro para institution ID inválido")
	}

	// Valid creation
	ev, err := svc.CreateEvent(ctx, "Tech Fest", "Rua A", "2026-10-01", "2026-10-02", "Desc", 1)
	if err != nil {
		t.Fatalf("erro ao criar evento: %v", err)
	}
	if ev.ID != 1 {
		t.Fatalf("esperava ID 1, obteve: %d", ev.ID)
	}

	// Add room: empty name
	_, err = svc.AddRoomInEvent(ctx, ev.ID, 1, "")
	if err == nil {
		t.Fatal("esperava erro ao adicionar sala com nome vazio")
	}

	// Add room: valid
	rooms, err := svc.AddRoomInEvent(ctx, ev.ID, 1, "Auditório Principal")
	if err != nil {
		t.Fatalf("erro ao adicionar sala: %v", err)
	}
	if len(rooms) != 1 {
		t.Fatalf("esperava 1 sala, obteve %d", len(rooms))
	}
	if rooms[0].Name != "Auditório Principal" {
		t.Fatalf("nome da sala incorreto: %s", rooms[0].Name)
	}
}
