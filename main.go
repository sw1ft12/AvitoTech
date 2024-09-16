package main

import (
    "avitoTech/internal/config"
    "avitoTech/internal/server"
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "log"
)

func main() {
    cfg := config.GetConfig()
    pgConn, err := pgxpool.New(context.TODO(), cfg.PostgresConn)
    if err != nil {
        log.Fatal(err.Error())
    }
    if err = pgConn.Ping(context.TODO()); err != nil {
        log.Fatal(err)
    }

    s := server.NewServer(pgConn)
    err = s.Run(cfg.Address)
    if err != nil {
        log.Fatal(err)
    }
}
