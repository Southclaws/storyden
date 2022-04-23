package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	fmt.Println(os.Getwd())

	m, err := migrate.New("file://migrations/", "postgres://default:default@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return err
	}

	fmt.Println(m)
	fmt.Println(m.Version())

	return m.Up()
}
