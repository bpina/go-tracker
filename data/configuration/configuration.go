package configuration

import (
    "os"
    "path/filepath"
    "bufio"
    "io"
    "encoding/json"
    "log"
)

type DatabaseConfiguration struct {
    Host        string `json:"host"`
    Port        string `json:"port"`
    Database    string `json:"database"`
    User        string `json:"user"`
    Password    string `json:"password"`
}

func NewDatabaseConfiguration() (config DatabaseConfiguration, err error) {
    err = os.Chdir("../config")
    if err != nil {
        return config, err
    }

    log.Printf("found config dir and changed to it")

    configDir, configErr := os.Getwd()
    if configErr != nil {
        return config, configErr
    }

    log.Printf("got cwd")

    filePath := filepath.Join(configDir, "database.json")
    file, fileErr := os.Open(filePath)
    defer file.Close()
    if fileErr != nil {
        return config, fileErr
    }

    log.Printf("opened config file")

    reader := bufio.NewReader(file)
    var eofErr error
    buffer := []string{}
    for {
        line, readerErr := reader.ReadString('\n')
        eofErr = readerErr
        if readerErr != nil {
            break
        }

        buffer = append(buffer, line)
    }
    if eofErr == io.EOF {
        eofErr = nil
    } else {
        return config, eofErr
    }

    var jsonString string
    for i := range buffer {
        jsonString += buffer[i]
    }

    log.Printf("read lines into buffer")

    log.Printf(jsonString)

    jsonErr := json.Unmarshal([]byte(jsonString), &config)
    if jsonErr != nil {
        return config, jsonErr
    }

    log.Printf("loaded databasconfiguration from json string")

    return config, err
}
