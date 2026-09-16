package queries

const CREATE_EVENT = "INSERT INTO EVENT(name, address, start_date, end_date, description) VALUES (?, ?, ?, ?, ?)"

const SET_OWNER_EVENT = "INSERT INTO OWNER_EVENT(event_id, owner_id) VALUES (?, ?)"

const ADD_ROOM_IN_EVENT = "INSERT INTO ROOM_IN_EVENT(event_id, name) VALUES (?, ?)"

const GET_EVENT_BY_NAME = `
	SELECT e.id, e.name, e.address, e.start_date, e.end_date, e.description, o.owner_id, i.name 
	FROM EVENT e
	JOIN OWNER_EVENT o ON e.id = o.event_id
	JOIN INSTITUTION_USER i ON o.owner_id = i.id
	WHERE e.name LIKE ?
	ORDER BY e.created_at DESC`

const GET_EVENT_BY_ID = `
	SELECT e.id, e.name, e.address, e.start_date, e.end_date, e.description, o.owner_id, i.name 
	FROM EVENT e
	JOIN OWNER_EVENT o ON e.id = o.event_id
	JOIN INSTITUTION_USER i ON o.owner_id = i.id
	WHERE e.id = ?`

const GET_EVENTS_BY_OWNER_ID = `
	SELECT e.id, e.name, e.address, e.start_date, e.end_date, e.description, o.owner_id, i.name 
	FROM EVENT e
	JOIN OWNER_EVENT o ON e.id = o.event_id
	JOIN INSTITUTION_USER i ON o.owner_id = i.id
	WHERE i.id = ?
	ORDER BY e.created_at DESC`

const GET_ROOMS_BY_EVENT_ID = "SELECT rie.id, rie.name FROM ROOM_IN_EVENT rie WHERE rie.event_id = ?"

const CHECK_EVENT_OWNER = "SELECT COUNT(*) FROM OWNER_EVENT WHERE owner_id = ? AND event_id = ?"

const GET_EVENTS = `
	SELECT e.id, e.name, e.address, e.start_date, e.end_date, e.description, o.owner_id, i.name 
	FROM EVENT e
	JOIN OWNER_EVENT o ON e.id = o.event_id
	JOIN INSTITUTION_USER i ON o.owner_id = i.id
	ORDER BY e.created_at DESC`

const UPDATE_EVENT = `
UPDATE EVENT 
SET name = COALESCE(NULLIF(?, ''), name),
    address = COALESCE(NULLIF(?, ''), address),
    start_date = COALESCE(NULLIF(?, ''), start_date),
    end_date = COALESCE(NULLIF(?, ''), end_date),
    description = COALESCE(NULLIF(?, ''), description)
WHERE id = ?`

const UPDATE_ROOM_NAME = "UPDATE ROOM_IN_EVENT SET name = ? WHERE id = ? AND event_id = ?"

const DELETE_EVENT = "DELETE FROM EVENT WHERE id = ?"

const DELETE_EVENT_ROOMS = "DELETE FROM ROOM_IN_EVENT WHERE event_id = ?"

const DELETE_ROOM = "DELETE FROM ROOM_IN_EVENT WHERE id = ? AND event_id = ?"