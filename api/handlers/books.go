package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/moonshoot17/bookdex/api/config"
	"github.com/moonshoot17/bookdex/api/models"
	"github.com/moonshoot17/bookdex/api/storage"
)

type Books struct {
	booksStorage *storage.Books
	customErrors *config.CustomErrors
	validator    *config.Validator
}

func newBooks(
	booksStorage *storage.Books,
	customErrors *config.CustomErrors,
	validator *config.Validator,
) *Books {
	return &Books{
		booksStorage: booksStorage,
		customErrors: customErrors,
		validator:    validator,
	}
}

func (b *Books) GetAll(ctx *fiber.Ctx) error {
	limit := ctx.QueryInt("limit", 9)
	offset := ctx.QueryInt("offset", 0)

	books, err := b.booksStorage.GetAll(limit, offset)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(books)
}

func (b *Books) GetAllBookmarked(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userId := claims["id"].(string)

	limit := ctx.QueryInt("limit", 9)
	offset := ctx.QueryInt("offset", 0)

	books, err := b.booksStorage.GetAllBookmarked(userId, limit, offset)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(books)
}

func (b *Books) GetById(ctx *fiber.Ctx) error {
	book, err := b.booksStorage.GetById(ctx.Params("id"))
	if err != nil {
		switch err {
		case b.customErrors.ErrNotFound:
			return fiber.ErrNotFound
		case b.customErrors.ErrInvalidId:
			return fiber.ErrBadRequest
		default:
			return fiber.ErrInternalServerError
		}
	}

	return ctx.JSON(book)
}

func (b *Books) GetAllBySubmitterId(ctx *fiber.Ctx) error {
	submitterId := ctx.Params("id")

	limit := ctx.QueryInt("limit", 9)
	offset := ctx.QueryInt("offset", 0)

	book, err := b.booksStorage.GetAllBySubmitterId(submitterId, limit, offset)
	if err != nil {
		switch err {
		case b.customErrors.ErrNotFound:
			return fiber.ErrNotFound
		case b.customErrors.ErrInvalidId:
			return fiber.ErrBadRequest
		default:
			return fiber.ErrInternalServerError
		}
	}

	return ctx.JSON(book)
}

func (b *Books) Create(ctx *fiber.Ctx) error {
	requestBody := new(models.InsertBookReqInput)
	err := ctx.BodyParser(requestBody)
	if err != nil {
		return fiber.ErrBadRequest
	}

	claims := ctx.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	submitterId := claims["id"].(string)

	err = b.validator.Validate(requestBody)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: err.Error(),
		}
	}

	id, err := b.booksStorage.Create(submitterId, requestBody)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(fiber.Map{"id": id})
}

func (b *Books) Update(ctx *fiber.Ctx) error {
	requestBody := new(models.UpdateBookInput)
	err := ctx.BodyParser(requestBody)
	if err != nil {
		return fiber.ErrBadRequest
	}

	err = b.booksStorage.Update(ctx.Params("id"), requestBody)
	if err != nil {
		switch err {
		case b.customErrors.ErrNotFound:
			return fiber.ErrNotFound
		case b.customErrors.ErrInvalidId:
			return fiber.ErrBadRequest
		default:
			return fiber.ErrInternalServerError
		}
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (b *Books) Delete(ctx *fiber.Ctx) error {
	err := b.booksStorage.Delete(ctx.Params("id"))
	if err != nil {
		switch err {
		case b.customErrors.ErrNotFound:
			return fiber.ErrNotFound
		case b.customErrors.ErrInvalidId:
			return fiber.ErrBadRequest
		default:
			return fiber.ErrInternalServerError
		}
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (b *Books) Bookmark(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userId := claims["id"].(string)

	err := b.booksStorage.Bookmark(ctx.Params("id"), userId)
	if err != nil {
		switch err {
		case b.customErrors.ErrNotFound:
			return fiber.ErrNotFound
		case b.customErrors.ErrInvalidId:
			return fiber.ErrBadRequest
		default:
			return fiber.ErrInternalServerError
		}
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (b *Books) Unbookmark(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userId := claims["id"].(string)

	err := b.booksStorage.Unbookmark(ctx.Params("id"), userId)
	if err != nil {
		switch err {
		case b.customErrors.ErrNotFound:
			return fiber.ErrNotFound
		case b.customErrors.ErrInvalidId:
			return fiber.ErrBadRequest
		default:
			return fiber.ErrInternalServerError
		}
	}

	return ctx.SendStatus(fiber.StatusOK)
}
