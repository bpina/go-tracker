package data

import (
    _ "github.com/bmizerany/pq"
    "database/sql"
    "github.com/bpina/go-tracker/data/configuration"
    "strings"
    "log"
)

var Database *sql.DB

func GetConnectionString(config configuration.DatabaseConfiguration) string {
    //TODO: figure out how to handle port configuration and sslmode
    properties := map[string] string {
        "dbname": config.Database,
        "host": config.Host,
        "user": config.User,
        "password": config.Password,
        "sslmode": "disable",
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

func OpenDatabaseConnection(config configuration.DatabaseConfiguration) error {
    connection := GetConnectionString(config)
    log.Printf("'%v'", connection)

    db, err := sql.Open("postgres", connection)
    if err != nil {
        return err
    }

    Database = db
    return err
}

func CloseDatabaseConnection() {
    Database.Close()
}

func InsertRow(table string, fields map[string] string) error {
    var columns string
    var values string

    i := 1
    max := len(fields)
    for key, value := range fields {

        column := key + ", "
        columnValue := value + ", "
        if i == max {
            column = key
            columnValue = value
        }

        columns = columns + column
        values = values + columnValue
        i += 1
    }

    log.Printf("columns: %v", columns)
    log.Printf("values: %v", values)

    sql := "INSERT INTO " + table + " (" + columns + ") VALUES (" + values + ")"

    log.Printf(sql)

    _, err := Database.Exec(sql)
    if err != nil {
        return err
    }

    return nil
}

func UpdateRow(table string, fields map[string] string, where string) error {
    var updates []rune

    i := 1
    max := len(fields)
    for key, value := range fields {
        update := key + "=" + value + ", "
        if i == max {
            update = key + "=" + value
        }

        updates = append(updates, []rune(update)...)
    }

    sql := "UPDATE " + table + " SET " + string(updates) + " WHERE " + where
    log.Printf(sql)

    _, err := Database.Exec(sql)
    if err != nil {
        return err
    }

    return nil
}

func Sanitize(sql string) string {
    //TODO: figure out what the fuck actually
    return strings.Replace(sql, "'", "\\'", -1)
}
