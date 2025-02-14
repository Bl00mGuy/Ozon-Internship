package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/Bl00mGuy/url-shortener/internal/repository"
	"github.com/Bl00mGuy/url-shortener/internal/server"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

func main() {
	repoType := flag.String("repo", "inmemory", "Тип хранилища: inmemory или postgres")
	postgresConnStr := flag.String("pgconn", "", "Строка подключения к PostgreSQL (используется, если repo=postgres)")
	grpcPort := flag.String("grpc-port", "50051", "Порт для gRPC сервера")
	httpPort := flag.String("http-port", "8080", "Порт для HTTP сервера")
	flag.Parse()

	logger := logrus.New()

	var repo repository.URLStorage
	var err error

	switch *repoType {
	case "inmemory":
		repo = repository.NewInMemoryURLStore()
	case "postgres":
		if *postgresConnStr == "" {
			logger.Fatal("Необходимо указать строку подключения к PostgreSQL при выборе репозитория postgres")
		}
		repo, err = repository.NewPostgresRepository(*postgresConnStr)
		if err != nil {
			logger.Fatalf("Не удалось инициализировать PostgreSQL репозиторий: %v", err)
		}

		db := repo.(*repository.PostgresRepository).Db
		if err := applyMigration(db, logger); err != nil {
			logger.Fatalf("Ошибка при применении миграции: %v", err)
		}
	default:
		logger.Fatalf("Неподдерживаемый тип репозитория: %s", *repoType)
	}

	go server.InitializeGRPCServer(*grpcPort, repo, logger)

	go server.InitializeHTTPGateway(*grpcPort, *httpPort, logger)

	select {}
}

func applyMigration(db *sql.DB, logger *logrus.Logger) error {
	migrationFile, err := ioutil.ReadFile("/home/bloomguy/GolandProjects/Ozon-Internship/urlshortener/migrations/migration.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	_, err = db.Exec(string(migrationFile))
	if err != nil {
		return fmt.Errorf("failed to apply migration: %w", err)
	}

	logger.Info("Migration applied successfully")
	return nil
}
