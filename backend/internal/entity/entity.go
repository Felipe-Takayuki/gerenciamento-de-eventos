package entity

type Institution struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email,omitempty"`
	Password string   `json:"-"`
	UserType string   `json:"user_type,omitempty"`
	CNPJ     string   `json:"cnpj,omitempty"`
	IsAdmin  bool     `json:"is_admin"`
	Events   []*Event `json:"events,omitempty"`
}

func NewInstitution(name, email, hashedPassword, cnpj string, isAdmin bool) *Institution {
	userType := "institution_user"
	if isAdmin {
		userType = "admin"
	}
	return &Institution{
		Name:     name,
		Email:    email,
		UserType: userType,
		Password: hashedPassword,
		CNPJ:     cnpj,
		IsAdmin:  isAdmin,
	}
}
