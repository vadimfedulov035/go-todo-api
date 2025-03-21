package initializers

import "log"

func Init() {
	// Connect to database
	if err := connectToDB(); err != nil {
		log.Fatal(err)
	}

	// Run migrations
	if err := runMigrations(); err != nil {
		log.Fatal(err)
	}
}
