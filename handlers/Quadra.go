package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/minacio00/easyCourt/database"
	"github.com/minacio00/easyCourt/types"
)

func CreateQuadra(c *fiber.Ctx) error {
	quadra := &types.Quadra{}
	err := json.Unmarshal(c.Body(), quadra)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	err = database.Db.Create(quadra).Error
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}

	return c.Status(201).JSON(&quadra)
}

func ReadQuadra(c *fiber.Ctx) error {
	quadra := &types.Quadra{}
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}
	err = database.Db.First(&quadra, id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(struct{ Message string }{Message: err.Error()})
	}
	return c.Status(200).JSON(&quadra)
}

func GetQuadras(c *fiber.Ctx) error {
	quadras := []types.Quadra{}

	err := database.Db.Find(&quadras).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(struct{ Message string }{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(&quadras)
}

func UpdateQuadra(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id == 0 {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}

	body := types.Quadra{}
	err = c.BodyParser(&body)
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}
	quadra := types.Quadra{}
	err = database.Db.First(&quadra, id).Error
	if err != nil {
		return c.Status(404).SendString("quadra não encontrada")
	}

	err = json.Unmarshal(c.Body(), &quadra)
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}

	database.Db.Model(&quadra).Updates(&quadra)
	return c.SendStatus(fiber.StatusNoContent)
}

func DeleteQuadra(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}
	quadra := &types.Quadra{}
	err = database.Db.First(&quadra, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(struct{ Message string }{Message: err.Error()})
	}
	database.Db.Delete(&quadra)
	return c.SendStatus(204)
}
