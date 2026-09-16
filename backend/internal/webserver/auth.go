package webserver

import (
	"net/http"
	"time"

	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/utils"
	"github.com/go-chi/jwtauth"
)

type AuthUser struct {
	ID       int64
	Name     string
	Email    string
	UserType string
	IsAdmin  bool
}

func getAuthUser(r *http.Request) (*AuthUser, error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || claims == nil {
		return nil, utils.ErrUnauthorized
	}

	idFloat, ok := claims["id"].(float64)
	if !ok {
		return nil, utils.ErrUnauthorized
	}

	userType, _ := claims["user_type"].(string)
	name, _ := claims["name"].(string)
	email, _ := claims["email"].(string)
	isAdmin, _ := claims["is_admin"].(bool)
	if userType == "admin" {
		isAdmin = true
	}

	return &AuthUser{
		ID:       int64(idFloat),
		Name:     name,
		Email:    email,
		UserType: userType,
		IsAdmin:  isAdmin,
	}, nil
}

func requireAdmin(r *http.Request) (*AuthUser, error) {
	user, err := getAuthUser(r)
	if err != nil {
		return nil, err
	}
	if !user.IsAdmin && user.UserType != "admin" {
		return nil, utils.ErrForbidden
	}
	return user, nil
}

func requireInstitution(r *http.Request) (*AuthUser, error) {
	user, err := getAuthUser(r)
	if err != nil {
		return nil, err
	}
	if user.UserType != "institution_user" && user.UserType != "admin" && !user.IsAdmin {
		return nil, utils.ErrForbidden
	}
	return user, nil
}

func generateToken(tokenAuth *jwtauth.JWTAuth, id int64, name, email, userType string, isAdmin bool) (string, error) {
	claims := map[string]interface{}{
		"id":        id,
		"name":      name,
		"email":     email,
		"user_type": userType,
		"is_admin":  isAdmin,
		"exp":       jwtauth.ExpireIn(24 * time.Hour),
	}
	_, tokenString, err := tokenAuth.Encode(claims)
	return tokenString, err
}
