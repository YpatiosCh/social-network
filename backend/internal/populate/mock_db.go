package main

import (
	"log"

	"platform.zone01.gr/git/kvamvasa/real-time-forum/backend/internal/database"
)

func main() {
	dbPool, err := database.StartDB()
	if err != nil {
		log.Fatal(err)
	}
	db := database.Database{Pool: dbPool}

	db.PopulateUsers()
	db.PopulatePosts()
	db.PopulateConversations()
}
