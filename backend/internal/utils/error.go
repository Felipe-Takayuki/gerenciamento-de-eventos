package utils

import "errors"

var (
	ErrNotFound            = errors.New("recurso não encontrado")
	ErrInvalidCredentials  = errors.New("credenciais inválidas")
	ErrUnauthorized        = errors.New("não autorizado")
	ErrForbidden           = errors.New("este usuário não possui essa permissão")
	ErrBadRequest          = errors.New("dados da requisição inválidos")
	ErrInstitutionNotOwner = errors.New("a instituição não possui o evento")
)

type ErrorMessage struct {
	Message string `json:"message"`
}