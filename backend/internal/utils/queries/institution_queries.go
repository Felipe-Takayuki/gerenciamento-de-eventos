package queries

const CREATE_INSTITUTION = "INSERT INTO INSTITUTION_USER(name, email, password, cnpj, is_admin) VALUES (?, ?, ?, ?, ?)"

const GET_INSTITUTION_BY_EMAIL = "SELECT id, name, email, password, cnpj, is_admin FROM INSTITUTION_USER WHERE email = ?"

const GET_INSTITUTION_BY_ID = "SELECT id, name, email, cnpj, is_admin FROM INSTITUTION_USER WHERE id = ?"

const DELETE_EVENT_OWNER = "DELETE FROM OWNER_EVENT WHERE owner_id = ? AND event_id = ?"
