package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/config"
	"github.com/Felipe-Takayuki/Adamas/adamas-api/internal/router"
	_ "github.com/go-sql-driver/mysql"
)

// @title           Gerenciamento de Eventos API
// @version         1.0
// @description     API RESTful para gerenciamento de eventos, salas e instituições.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Suporte da API
// @contact.email  suporte@adamas.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:3000
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Insira o token JWT no formato: Bearer {seu_token}
func main() {
	cfg := config.Load()

	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		log.Fatalf("falha ao abrir conexão com o banco de dados: %v", err)
	}
	defer db.Close()

	// Configuração do pool de conexões
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()

	if err := db.PingContext(pingCtx); err != nil {
		log.Printf("aviso: não foi possível conectar ao banco de dados imediatamente: %v", err)
	} else {
		log.Println("conexão com o banco de dados estabelecida com sucesso")
	}

	appRouter := router.NewRouter(db, cfg)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      appRouter,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Encerramento gracioso (Graceful Shutdown)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("servidor iniciado com sucesso na porta %s", cfg.Port)
		log.Printf("documentação Swagger disponível em: http://localhost:%s/swagger/index.html", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("erro ao iniciar o servidor: %v", err)
		}
	}()

	<-stop
	log.Println("recebido sinal de encerramento, finalizando servidor...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("falha ao encerrar o servidor com segurança: %v", err)
	}

	log.Println("servidor finalizado com sucesso")
}
