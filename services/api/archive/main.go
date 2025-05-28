package main

import (
	"os"

	"github.com/TheRangiCrew/WITS/services/api/archive/db"
	"github.com/TheRangiCrew/WITS/services/api/archive/routes"
	"github.com/joho/godotenv"
)

func main() {

	if len(os.Args) > 1 {
		godotenv.Load(os.Args[1])
	}

	pool, err := db.New()
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	db.DB = pool

	engine := routes.Setup()

	engine.Run(":8081")

}
