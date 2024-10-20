package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/moonshoot17/bookdex/api/config"
	"github.com/moonshoot17/bookdex/api/models"
	"github.com/moonshoot17/bookdex/api/storage"
)

type Comments struct {
	commentsStorage *storage.Comments
	customErrors    *config.CustomErrors
}

func newComments(
	commentsStorage *storage.Comments,
	customErrors *config.CustomErrors,
) *Comments {
	return &Comments{
		commentsStorage: commentsStorage,
		customErrors:    customErrors,
	}
}

func (r *Comments) GetAllByBookId(ctx *fiber.Ctx) error {
	comments, err := r.commentsStorage.GetAllByBookId(ctx.Params("id"))

	if err != nil {
		if err == r.customErrors.ErrInvalidId {
			return fiber.ErrBadRequest
		}

		return fiber.ErrInternalServerError
	}

	return ctx.JSON(comments)
}

func (r *Comments) Create(ctx *fiber.Ctx) error {
	requestBody := new(models.InsertCommentReqInput)
	err := ctx.BodyParser(requestBody)
	if err != nil {
		return fiber.ErrBadRequest
	}

	claims := ctx.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	authorId := claims["id"].(string)

	booksId := ctx.Params("id")

	id, err := r.commentsStorage.Create(authorId, booksId, requestBody)
	if err != nil {
		if err == r.customErrors.ErrInvalidId {
			return fiber.ErrBadRequest
		}
		return fiber.ErrInternalServerError
	}

	return ctx.JSON(fiber.Map{"id": id})
}

func (r *Comments) Update(ctx *fiber.Ctx) error {
	requestBody := new(models.UpdateCommentInput)
	err := ctx.BodyParser(requestBody)
	if err != nil {
		return fiber.ErrBadRequest
	}

	claims := ctx.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userId := claims["id"].(string)

	commentId := ctx.Params("id")

	comment, err := r.commentsStorage.GetById(commentId)
	if err != nil {
		return fiber.ErrInternalServerError
	}

	if comment.AuthorID.Hex() != userId {
		return fiber.ErrForbidden
	}

	err = r.commentsStorage.Update(commentId, requestBody)
	if err != nil {
		switch err {
		case r.customErrors.ErrNotFound:
			return fiber.ErrNotFound
		case r.customErrors.ErrInvalidId:
			return fiber.ErrBadRequest
		default:
			return fiber.ErrInternalServerError
		}
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (r *Comments) Delete(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
	userId := claims["id"].(string)

	commentId := ctx.Params("id")

	comment, err := r.commentsStorage.GetById(commentId)
	if err != nil {
		switch err {
		case r.customErrors.ErrNotFound:
			return fiber.ErrNotFound
		case r.customErrors.ErrInvalidId:
			return fiber.ErrBadRequest
		default:
			return fiber.ErrInternalServerError
		}
	}

	if comment.AuthorID.Hex() != userId {
		return fiber.ErrForbidden
	}

	err = r.commentsStorage.Delete(commentId)
	if err != nil {
		if err == r.customErrors.ErrInvalidId {
			return fiber.ErrBadRequest
		}
		return fiber.ErrInternalServerError
	}

	return ctx.SendStatus(fiber.StatusOK)
}
