package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/entity"
	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/utils"
	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/utils/queries"
)

type EventDB struct {
	db *sql.DB
}

func NewEventDB(db *sql.DB) *EventDB {
	return &EventDB{
		db: db,
	}
}

func (edb *EventDB) CreateEvent(ctx context.Context, name, address, startDate, endDate, description string, institutionID int64) (*entity.Event, error) {
	tx, err := edb.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	event := entity.NewEvent(name, address, startDate, endDate, description, institutionID)
	result, err := tx.ExecContext(ctx, queries.CREATE_EVENT, event.Name, event.Address, event.StartDate, event.EndDate, event.Description)
	if err != nil {
		return nil, err
	}
	event.ID, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx, queries.SET_OWNER_EVENT, event.ID, event.InstitutionID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return event, nil
}

func (edb *EventDB) DeleteEvent(ctx context.Context, eventID, institutionID int64) error {
	isOwner, err := edb.IsEventOwner(ctx, eventID, institutionID)
	if err != nil {
		return err
	}
	if !isOwner {
		return utils.ErrInstitutionNotOwner
	}

	tx, err := edb.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, queries.DELETE_EVENT_ROOMS, eventID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, queries.DELETE_EVENT_OWNER, institutionID, eventID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, queries.DELETE_EVENT, eventID); err != nil {
		return err
	}

	return tx.Commit()
}

func (edb *EventDB) EditRoom(ctx context.Context, roomID, eventID, ownerID int64, roomName string) (*entity.RoomEvent, error) {
	isOwner, err := edb.IsEventOwner(ctx, eventID, ownerID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, utils.ErrInstitutionNotOwner
	}

	if roomName != "" {
		_, err := edb.db.ExecContext(ctx, queries.UPDATE_ROOM_NAME, roomName, roomID, eventID)
		if err != nil {
			return nil, err
		}
	}
	return &entity.RoomEvent{ID: roomID, Name: roomName}, nil
}

func (edb *EventDB) DeleteRoom(ctx context.Context, roomID, eventID, ownerID int64) error {
	isOwner, err := edb.IsEventOwner(ctx, eventID, ownerID)
	if err != nil {
		return err
	}
	if !isOwner {
		return utils.ErrInstitutionNotOwner
	}

	_, err = edb.db.ExecContext(ctx, queries.DELETE_ROOM, roomID, eventID)
	return err
}

func (edb *EventDB) AddRoomInEvent(ctx context.Context, eventID, ownerID int64, roomName string) ([]*entity.RoomEvent, error) {
	isOwner, err := edb.IsEventOwner(ctx, eventID, ownerID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, utils.ErrInstitutionNotOwner
	}

	_, err = edb.db.ExecContext(ctx, queries.ADD_ROOM_IN_EVENT, eventID, roomName)
	if err != nil {
		return nil, err
	}
	return edb.getRoomsByEventID(ctx, eventID)
}

func (edb *EventDB) GetEventByName(ctx context.Context, name string) ([]*entity.Event, error) {
	rows, err := edb.db.QueryContext(ctx, queries.GET_EVENT_BY_NAME, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*entity.Event
	for rows.Next() {
		var event entity.Event
		if err := rows.Scan(&event.ID, &event.Name, &event.Address, &event.StartDate, &event.EndDate, &event.Description, &event.InstitutionID, &event.InstitutionName); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (edb *EventDB) GetEvents(ctx context.Context) ([]*entity.Event, error) {
	rows, err := edb.db.QueryContext(ctx, queries.GET_EVENTS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*entity.Event
	for rows.Next() {
		var event entity.Event
		if err := rows.Scan(&event.ID, &event.Name, &event.Address, &event.StartDate, &event.EndDate, &event.Description, &event.InstitutionID, &event.InstitutionName); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (edb *EventDB) GetEventByOwnerID(ctx context.Context, ownerID int64) ([]*entity.Event, error) {
	rows, err := edb.db.QueryContext(ctx, queries.GET_EVENTS_BY_OWNER_ID, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*entity.Event
	for rows.Next() {
		var event entity.Event
		if err := rows.Scan(&event.ID, &event.Name, &event.Address, &event.StartDate, &event.EndDate, &event.Description, &event.InstitutionID, &event.InstitutionName); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (edb *EventDB) GetRoomsByEventID(ctx context.Context, eventID, ownerID int64) ([]*entity.RoomEvent, error) {
	isOwner, err := edb.IsEventOwner(ctx, eventID, ownerID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, utils.ErrInstitutionNotOwner
	}
	return edb.getRoomsByEventID(ctx, eventID)
}

func (edb *EventDB) getRoomsByEventID(ctx context.Context, eventID int64) ([]*entity.RoomEvent, error) {
	rows, err := edb.db.QueryContext(ctx, queries.GET_ROOMS_BY_EVENT_ID, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*entity.RoomEvent
	for rows.Next() {
		var room entity.RoomEvent
		if err := rows.Scan(&room.ID, &room.Name); err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (edb *EventDB) GetEventByID(ctx context.Context, eventID int64) (*entity.Event, error) {
	event := &entity.Event{}
	err := edb.db.QueryRowContext(ctx, queries.GET_EVENT_BY_ID, eventID).Scan(
		&event.ID, &event.Name, &event.Address, &event.StartDate, &event.EndDate, &event.Description, &event.InstitutionID, &event.InstitutionName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}

	rooms, err := edb.getRoomsByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	event.Rooms = rooms

	return event, nil
}

func (edb *EventDB) EditEvent(ctx context.Context, eventID, ownerID int64, name, address, startDate, endDate, description string) (*entity.Event, error) {
	isOwner, err := edb.IsEventOwner(ctx, eventID, ownerID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, utils.ErrInstitutionNotOwner
	}

	_, err = edb.db.ExecContext(ctx, queries.UPDATE_EVENT, name, address, startDate, endDate, description, eventID)
	if err != nil {
		return nil, err
	}

	return edb.GetEventByID(ctx, eventID)
}

func (edb *EventDB) IsEventOwner(ctx context.Context, eventID, ownerID int64) (bool, error) {
	var count int
	err := edb.db.QueryRowContext(ctx, queries.CHECK_EVENT_OWNER, ownerID, eventID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
