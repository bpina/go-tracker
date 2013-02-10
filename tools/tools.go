package tools

import (
    "strconv"
)

func IntOrDefault(s string) int {
    i, err := strconv.Atoi(s)

    if err != nil {
        i = 0
    }

    return i
}

func FormatErrors(errors []string) string {
    message := ""

    for i := range errors {
        message += errors[i] + "\n"
    }

    return message
}


