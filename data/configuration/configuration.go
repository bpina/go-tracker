package configuration

import (
    "os"
    "path/filepath"
    "bufio"
    "io"
    "encoding/json"
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

    configDir, err := os.Getwd()
    if err != nil {
        return config, err
    }

    filePath := filepath.Join(configDir, "database.json")
    file, err := os.Open(filePath)
    if err != nil {
        return config, err
    }
    defer file.Close()

    reader := bufio.NewReader(file)
    buffer := []string{}
    for {
      line, err := reader.ReadString('\n')
      if err == io.EOF {
        break
      } else {
        if err != nil {
          return config, err
        }
      }

      buffer = append(buffer, line)
    }

    var jsonString string
    for i := range buffer {
        jsonString += buffer[i]
    }

    err = json.Unmarshal([]byte(jsonString), &config)
    if err != nil {
        return config, err
    }

    return config, err
}
