package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() {
	var err error
	Pool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}

	err = Pool.Ping(context.Background())
	if err != nil {
		log.Fatal("Database ping failed: ", err)
	}

	log.Println("✅ Database connected successfully!")
}
