package service

import (
	"context"
	"errors"
	"strings"

	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/entity"
	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/utils"
)

type InstitutionRepository interface {
	CreateInstitution(ctx context.Context, name, email, hashedPassword, cnpj string, isAdmin bool) (*entity.Institution, error)
	GetInstitutionByEmail(ctx context.Context, email string) (*entity.Institution, error)
	GetInstitutionByID(ctx context.Context, institutionID int64) (*entity.Institution, error)
}

type InstitutionService struct {
	institutionRepo InstitutionRepository
}

func NewInstitutionService(institutionRepo InstitutionRepository) *InstitutionService {
	return &InstitutionService{
		institutionRepo: institutionRepo,
	}
}

func (is *InstitutionService) CreateInstitution(ctx context.Context, name, email, password, cnpj string, isAdmin bool) (*entity.Institution, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	cnpj = strings.TrimSpace(cnpj)

	if name == "" || email == "" || password == "" || cnpj == "" {
		return nil, errors.New("todos os campos são obrigatórios")
	}
	if !strings.Contains(email, "@") {
		return nil, errors.New("formato de email inválido")
	}
	if len(password) < 6 {
		return nil, errors.New("a senha deve ter no mínimo 6 caracteres")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	return is.institutionRepo.CreateInstitution(ctx, name, email, hashedPassword, cnpj, isAdmin)
}

func (is *InstitutionService) LoginInstitution(ctx context.Context, email, password string) (*entity.Institution, error) {
	email = strings.TrimSpace(email)
	if email == "" || password == "" {
		return nil, utils.ErrInvalidCredentials
	}

	institution, err := is.institutionRepo.GetInstitutionByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return nil, utils.ErrInvalidCredentials
		}
		return nil, err
	}

	if !utils.CheckPasswordHash(password, institution.Password) {
		return nil, utils.ErrInvalidCredentials
	}

	institution.Password = ""
	return institution, nil
}

func (is *InstitutionService) GetInstitutionByID(ctx context.Context, institutionID int64) (*entity.Institution, error) {
	if institutionID <= 0 {
		return nil, utils.ErrBadRequest
	}
	return is.institutionRepo.GetInstitutionByID(ctx, institutionID)
}