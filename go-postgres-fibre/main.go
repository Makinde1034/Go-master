package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/makinde1034/go-postgres-fibre/controllers"
	"github.com/makinde1034/go-postgres-fibre/models"
	"github.com/makinde1034/go-postgres-fibre/storage"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.NewConnection(os.Getenv("POSTGRES_URL"))

	if err != nil {
		log.Fatal(err)
	}

	err = models.AutoMigrate(db)

	if err != nil {
		log.Fatal(err)
	}

	r :=
		controllers.Repository{
			DB: db,
		}

	app := fiber.New()
	controllers.SetupRoutes(app, (*controllers.Repository)(&r))
	app.Listen(":8000")
	fmt.Println("Listening on port :8080")
}
