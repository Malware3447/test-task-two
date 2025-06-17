package main

import (
	"context"
	"fmt"
	"github.com/Malware3447/configo"
	"github.com/Malware3447/spg"
	"log"
	"os"
	"os/signal"
	"syscall"
	"test-task-two/internal/config"
)

func main() {
	const op = "cmd.runner.main"
	cfg, _ := configo.MustLoad[config.Config]()

	ctx := context.Background()
	ctx = context.WithValue(ctx, "main", op)

	poolPg, err := spg.NewClient(ctx, &cfg.DatabasePg)
	if err != nil {
		log.Println(fmt.Errorf("ошибка при запуске Postgres: %s", err))
		panic(err)
	}
	log.Println("Postgres успешно запущен")
	log.Println("Сервис успешно запущен")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-quit:
		log.Println("Завершение работы сервиса")
	}

	log.Println("Сервис успешно завершил работу")
}
