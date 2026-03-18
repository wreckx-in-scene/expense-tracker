package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func Connect(connString string) error {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return fmt.Errorf("unable to connect to database : %w", err)
	}

	DB = conn
	fmt.Println("Connected to database!")
	return nil
}
