package initializer

import (
    "log"
    "os"
)

func LoadEnvVariables() {
    dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
    if dbConnectionString == "" {
        log.Fatal("DB_CONNECTION_STRING environment variable not set")
    }

    secret := os.Getenv("SECRET")
    if secret == "" {
        log.Fatal("SECRET environment variable not set")
    }

    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        log.Fatal("JWT_SECRET environment variable not set")
    }

    issuer := os.Getenv("ISSUER")
    if issuer == "" {
        log.Fatal("ISSUER environment variable not set")
    }

    audience := os.Getenv("AUDIENCE")
    if audience == "" {
        log.Fatal("AUDIENCE environment variable not set")
    }
}