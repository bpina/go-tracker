package data

import (
	"fmt"
	"github.com/bpina/go-tracker/data/configuration"
	"github.com/jackc/pgx"
)

var Database *pgx.ConnPool

func GetConnectionString(config configuration.DatabaseConfiguration) string {
	var sslMode string
	if config.SSLMode != "" {
		sslMode = config.SSLMode
	} else {
		sslMode = "disabled"
	}

	var port string
	if config.Port != "" {
		port = config.Port
	} else {
		port = "5432"
	}

	properties := map[string]string{
		"dbname":   config.Database,
		"host":     config.Host,
		"user":     config.User,
		"password": config.Password,
		"port":     port,
		"sslmode":  sslMode,
	}

	runes := []rune{}
	i := 1
	max := len(properties)
	for key, value := range properties {
		property := key + "=" + value
		if i != max {
			property = property + " "
		}
		runes = append(runes, []rune(property)...)
		i += 1
	}

	return string(runes)
}

func OpenDatabaseConnection(config configuration.DatabaseConfiguration) (pool *pgx.ConnPool, err error) {
	connectionUri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.SSLMode)

	connectionConfig, err := pgx.ParseURI(connectionUri)

	if err != nil {
		return pool, err
	}

	maxConnections := 50

	poolConfig := pgx.ConnPoolConfig{connectionConfig, maxConnections, nil}

	pool, err = pgx.NewConnPool(poolConfig)

	if err != nil {
		return pool, err
	}

	Database = pool
	return pool, err
}

func CloseDatabaseConnection() {
	Database.Close()
}
