package db

import (
	"fmt"
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
)

func NewDB() (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "core",
		Addr:     "db:5432",
	})

	err := migrate(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func migrate(db *pg.DB) error {
	collection := migrations.NewCollection()
	err := collection.DiscoverSQLMigrations("migrations")
	if err != nil {
		return err
	}
	_, _, err = collection.Run(db, "init")
	if err != nil {
		return errw
	}
	oldVer, newVer, err := collection.Run(db, "up")
	if err != nil {
		return err
	}
	if oldVer != newVer {
		fmt.Printf("Migrated from version %d to version %d\n", oldVer, newVer)
	} else {
		fmt.Printf("Version %d\n", oldVer)
	}
	return nil
}
