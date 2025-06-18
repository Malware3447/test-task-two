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
	"test-task-two/internal/api"
	"test-task-two/internal/api/crut"
	"test-task-two/internal/app"
	"test-task-two/internal/config"
	PgRepository "test-task-two/internal/db/pg"
	ServApi "test-task-two/internal/service/api"
	"test-task-two/internal/service/db/pg"
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

	PgParams := PgRepository.Params{Db: poolPg}
	PgRepo := PgRepository.NewRepository(&PgParams)

	PgService := pg.NewService(PgRepo)

	CrutParams := crut.Params{RepoPg: PgService}
	Crut := crut.NewCrut(&CrutParams)

	ApiServiceParams := ServApi.Params{Api: Crut}
	ApiService := ServApi.NewService(&ApiServiceParams)

	router := api.NewRouter(ApiService)

	App := app.NewApp(router)

	App.Init(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-quit:
		log.Println("Завершение работы сервиса")
	}

	log.Println("Сервис успешно завершил работу")
}
