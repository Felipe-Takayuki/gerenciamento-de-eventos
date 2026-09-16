package utils

import (
	"testing"
)

func TestHashPasswordAndCheck(t *testing.T) {
	password := "minhasenhasupersegura123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("esperava erro nil ao gerar hash, obteve: %v", err)
	}

	if hash == password {
		t.Fatal("o hash gerado não deve ser igual à senha pura")
	}

	if !CheckPasswordHash(password, hash) {
		t.Fatal("esperava que a senha correta validasse contra o hash bcrypt")
	}

	if CheckPasswordHash("senha_incorreta", hash) {
		t.Fatal("senha incorreta não deveria ser validada")
	}
}

func TestLegacySHA256Fallback(t *testing.T) {
	password := "senha_legada_123"
	shaHash := EncriptKey(password)

	if !CheckPasswordHash(password, shaHash) {
		t.Fatal("esperava que o hash legado SHA-256 fosse validado com sucesso pelo fallback")
	}

	if CheckPasswordHash("senha_incorreta", shaHash) {
		t.Fatal("senha incorreta não deveria ser validada com hash SHA-256")
	}
}
