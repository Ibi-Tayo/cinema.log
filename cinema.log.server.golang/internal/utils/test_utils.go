package utils

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    "cinema.log.server.golang/internal/migration"
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
    "github.com/testcontainers/testcontainers-go/wait"
    _ "github.com/jackc/pgx/v5/stdlib"
)

type TestDatabase struct {
    DB           *sql.DB
    TeardownFunc func(context.Context, ...testcontainers.TerminateOption) error
    ConnStr      string
}

func StartTestPostgres() (*TestDatabase, error) {
    var (
        dbName = "test_database"
        dbPwd  = "test_password"
        dbUser = "test_user"
    )

    dbContainer, err := postgres.Run(
        context.Background(),
        "postgres:latest",
        postgres.WithDatabase(dbName),
        postgres.WithUsername(dbUser),
        postgres.WithPassword(dbPwd),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(5*time.Second)),
    )
    if err != nil {
        return nil, err
    }

    dbHost, err := dbContainer.Host(context.Background())
    if err != nil {
        return nil, err
    }

    dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
    if err != nil {
        return nil, err
    }

    connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
        dbUser, dbPwd, dbHost, dbPort.Port(), dbName)

    db, err := sql.Open("pgx", connStr)
    if err != nil {
        return nil, err
    }

    // Run migrations
    if err := migration.RunMigrations(db); err != nil {
        return nil, fmt.Errorf("failed to run migrations: %w", err)
    }

    return &TestDatabase{
        DB:           db,
        TeardownFunc: dbContainer.Terminate,
        ConnStr:      connStr,
    }, nil
}

func (td *TestDatabase) Close() error {
    if td.DB != nil {
        td.DB.Close()
    }
    if td.TeardownFunc != nil {
        return td.TeardownFunc(context.Background())
    }
    return nil
}