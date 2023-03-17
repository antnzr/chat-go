package db

import (
	"context"
	"fmt"
	"log"
	"path"
	"runtime"

	"net/url"
	"path/filepath"
	"strconv"
	"time"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/antnzr/chat-go/config"
	"github.com/docker/go-connections/nat"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabase struct {
	DbInstance *pgxpool.Pool
	DbAddress  string
	container  testcontainers.Container
}

func SetupTestDatabase(conf *config.Config) *TestDatabase {
	// setup db container
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*60)

	container, dbInstance, dbURL, err := createContainer(ctx, conf)
	if err != nil {
		log.Fatal("failed to setup test", err)
	}

	err = migrate(dbURL)
	if err != nil {
		log.Fatal("failed to perform db migration", err)
	}
	cancel()

	return &TestDatabase{
		container:  container,
		DbInstance: dbInstance,
		DbAddress:  dbURL,
	}
}

func (tdb *TestDatabase) TearDown() {
	tdb.DbInstance.Close()
	// remove test container
	_ = tdb.container.Terminate(context.Background())
}

func createContainer(ctx context.Context, conf *config.Config) (testcontainers.Container, *pgxpool.Pool, string, error) {
	// setup db container
	pgPort := fmt.Sprintf("%d/tcp", conf.PgPort)
	containerReq := testcontainers.ContainerRequest{
		ExposedPorts: []string{pgPort},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(pgPort)),
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(1).
				WithStartupTimeout(5*time.Second)).
			WithDeadline(1 * time.Minute),
		FromDockerfile: testcontainers.FromDockerfile{
			Dockerfile: "Dockerfile",
			Context:    getDbDirAbsolutePath(),
		},
		Env: map[string]string{
			"POSTGRES_DB":       conf.PgDbName,
			"POSTGRES_PASSWORD": conf.PgPassword,
			"POSTGRES_USER":     conf.PgUser,
		},
	}

	// start db container
	container, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		return nil, nil, "", err
	}

	dbURL, err := getDbURL(ctx, container, conf)
	if err != nil {
		return nil, nil, "", err
	}

	conf.DatabaseURL = dbURL
	dbPool, err := DBPool(ctx, *conf)
	if err != nil {
		return nil, nil, "", err
	}

	return container, dbPool, dbURL, nil
}

func getDbDirAbsolutePath() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Dir(filename)
}

func migrate(dbURL string) error {
	url, err := url.Parse(dbURL)
	if err != nil {
		return err
	}
	migr := dbmate.New(url)

	migr.AutoDumpSchema = false
	migr.WaitInterval = time.Millisecond
	migr.WaitTimeout = 5 * time.Millisecond
	migr.MigrationsDir = getMigrationDir()

	fmt.Println("Migrations:")
	migrations, err := migr.FindMigrations()
	if err != nil {
		fmt.Println("No migrations")
		return err
	}
	for _, m := range migrations {
		fmt.Println(m.Version, m.FilePath)
	}

	fmt.Println("\nApplying...")
	err = migr.CreateAndMigrate()
	if err != nil {
		return err
	}
	return nil
}

func getMigrationDir() string {
	dir := getDbDirAbsolutePath()
	return filepath.Join(dir, "migrations")
}

func getDbURL(ctx context.Context, container testcontainers.Container, config *config.Config) (string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", err
	}

	port, err := container.MappedPort(ctx, nat.Port(strconv.Itoa(config.PgPort)))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
			"postgres://%s:%s@%v:%v/%s?sslmode=disable",
			config.PgUser,
			config.PgPassword,
			host,
			port.Port(),
			config.PgDbName),
		nil
}
