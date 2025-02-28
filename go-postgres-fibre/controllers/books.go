package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/makinde1034/go-postgres-fibre/models"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := models.Book{}

	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusExpectationFailed).JSON(&fiber.Map{"Messge": "Request failed"})
		return err
	}

	err = r.DB.Create(&book).Error

	if err != nil {
		context.Status(http.StatusExpectationFailed).JSON(&fiber.Map{"Messge": "Request failed"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"Messge": "Book created"})

	return nil
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	books := []models.Book{}

	err := r.DB.Find(&books).Error

	if err != nil {
		context.Status(http.StatusExpectationFailed).JSON(&fiber.Map{"Messge": "Failed Book created"})
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"Messge": "Books fetched successfully", "data": books})

	return nil

}

func (r *Repository) Ttest(context *fiber.Ctx) error {

	context.Status(http.StatusOK).JSON(&fiber.Map{"Messge": "OK"})

	return nil

}

func (r *Repository) UpdateBook(context *fiber.Ctx) error {

	book := models.Book{}

	err := r.DB.Model(&book).Where("id = ?", 1).Update("title", "new title")

	if err != nil {
		context.Status(http.StatusExpectationFailed).JSON(&fiber.Map{"Messge": "Failed update"})
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"Messge": "OK"})

	return nil

}

func SetupRoutes(app *fiber.App, r *Repository) {
	api := app.Group("/api")
	api.Post("/create-book", r.CreateBook)
	api.Get("/get-books", r.GetBooks)
	api.Post("/update-book", r.UpdateBook)
	api.Get("/test", r.Ttest)
}
