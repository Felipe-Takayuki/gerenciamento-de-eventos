package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/entity"
	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/utils"
	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/utils/queries"
)

type InstitutionDB struct {
	db *sql.DB
}

func NewInstitutionDB(db *sql.DB) *InstitutionDB {
	return &InstitutionDB{
		db: db,
	}
}

func (idb *InstitutionDB) CreateInstitution(ctx context.Context, name, email, hashedPassword, cnpj string, isAdmin bool) (*entity.Institution, error) {
	institution := entity.NewInstitution(name, email, hashedPassword, cnpj, isAdmin)
	result, err := idb.db.ExecContext(ctx, queries.CREATE_INSTITUTION, institution.Name, institution.Email, institution.Password, institution.CNPJ, institution.IsAdmin)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	institution.ID = id
	return institution, nil
}

func (idb *InstitutionDB) GetInstitutionByEmail(ctx context.Context, email string) (*entity.Institution, error) {
	var institution entity.Institution
	err := idb.db.QueryRowContext(ctx, queries.GET_INSTITUTION_BY_EMAIL, email).Scan(
		&institution.ID, &institution.Name, &institution.Email, &institution.Password, &institution.CNPJ, &institution.IsAdmin,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	institution.UserType = "institution_user"
	if institution.IsAdmin {
		institution.UserType = "admin"
	}
	return &institution, nil
}

func (idb *InstitutionDB) GetInstitutionByID(ctx context.Context, institutionID int64) (*entity.Institution, error) {
	var institution entity.Institution
	err := idb.db.QueryRowContext(ctx, queries.GET_INSTITUTION_BY_ID, institutionID).Scan(
		&institution.ID, &institution.Name, &institution.Email, &institution.CNPJ, &institution.IsAdmin,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	institution.UserType = "institution_user"
	if institution.IsAdmin {
		institution.UserType = "admin"
	}
	return &institution, nil
}