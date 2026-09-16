package service

import (
	"context"
	"testing"

	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/entity"
	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/utils"
)

type mockInstitutionRepo struct {
	institutions map[string]*entity.Institution
}

func newMockInstitutionRepo() *mockInstitutionRepo {
	return &mockInstitutionRepo{
		institutions: make(map[string]*entity.Institution),
	}
}

func (m *mockInstitutionRepo) CreateInstitution(ctx context.Context, name, email, hashedPassword, cnpj string, isAdmin bool) (*entity.Institution, error) {
	inst := entity.NewInstitution(name, email, hashedPassword, cnpj, isAdmin)
	inst.ID = int64(len(m.institutions) + 1)
	m.institutions[email] = inst
	return inst, nil
}

func (m *mockInstitutionRepo) GetInstitutionByEmail(ctx context.Context, email string) (*entity.Institution, error) {
	inst, ok := m.institutions[email]
	if !ok {
		return nil, utils.ErrNotFound
	}
	return inst, nil
}

func (m *mockInstitutionRepo) GetInstitutionByID(ctx context.Context, institutionID int64) (*entity.Institution, error) {
	for _, inst := range m.institutions {
		if inst.ID == institutionID {
			return inst, nil
		}
	}
	return nil, utils.ErrNotFound
}

func TestInstitutionService_CreateAndLogin(t *testing.T) {
	repo := newMockInstitutionRepo()
	svc := NewInstitutionService(repo)
	ctx := context.Background()

	// Validation: empty name
	_, err := svc.CreateInstitution(ctx, "", "inst@example.com", "senha123", "12345678000199", false)
	if err == nil {
		t.Fatal("esperava erro para nome vazio")
	}

	// Validation: short password
	_, err = svc.CreateInstitution(ctx, "Univ Tech", "inst@example.com", "123", "12345678000199", false)
	if err == nil {
		t.Fatal("esperava erro para senha curta")
	}

	// Valid creation (regular institution)
	inst, err := svc.CreateInstitution(ctx, "Univ Tech", "inst@example.com", "senhaSegura123", "12345678000199", false)
	if err != nil {
		t.Fatalf("erro ao criar instituição: %v", err)
	}
	if inst.ID == 0 {
		t.Fatal("esperava ID atribuído")
	}
	if inst.IsAdmin {
		t.Fatal("esperava IsAdmin falso para instituição comum")
	}

	// Valid creation (admin institution)
	adminInst, err := svc.CreateInstitution(ctx, "Admin Org", "admin@example.com", "senhaAdmin123", "00000000000000", true)
	if err != nil {
		t.Fatalf("erro ao criar instituição admin: %v", err)
	}
	if !adminInst.IsAdmin || adminInst.UserType != "admin" {
		t.Fatal("esperava IsAdmin verdadeiro e UserType admin")
	}

	// Valid login
	logged, err := svc.LoginInstitution(ctx, "inst@example.com", "senhaSegura123")
	if err != nil {
		t.Fatalf("esperava sucesso no login: %v", err)
	}
	if logged.Name != "Univ Tech" {
		t.Fatalf("nome incorreto: %s", logged.Name)
	}

	// Wrong password
	_, err = svc.LoginInstitution(ctx, "inst@example.com", "senha_incorreta")
	if err != utils.ErrInvalidCredentials {
		t.Fatalf("esperava ErrInvalidCredentials, obteve: %v", err)
	}
}
