package main

import (
	"engine-db/database"
	"engine-db/database/seeder"
	"testing"
)

func TestEngineDB(t *testing.T) {
	db := database.New()

	seeder := seeder.NewSeeder(db)

	// seeder.SeedUsers()
	seeder.SeedUsersFromCSV()
	seeder.SeedEventCycle()
}
