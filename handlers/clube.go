package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/minacio00/easyCourt/database"
	"github.com/minacio00/easyCourt/types"
)

func CreateClube(c *fiber.Ctx) error {
	clube := &types.Clube{}
	err := json.Unmarshal(c.Body(), clube)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	err = clube.Validate()
	if err != nil {
		return c.Status(400).JSON(struct{ Message string }{Message: err.Error()})
	}
	err = database.Db.Create(clube).Error
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}

	return c.Status(201).JSON(&clube)
}

func ReadClube(c *fiber.Ctx) error {
	clube := &types.Clube{}
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}
	err = database.Db.Preload("quadras").Preload("clientes").First(&clube, id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(struct{ Message string }{Message: err.Error()})
	}
	return c.Status(200).JSON(&clube)
}

func GetClubes(c *fiber.Ctx) error {
	clubes := []types.Clube{}

	err := database.Db.Find(&clubes).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(struct{ Message string }{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(&clubes)
}

func UpdateClube(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id == 0 {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}

	body := types.Clube{}
	err = c.BodyParser(&body)
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}
	clube := types.Clube{}
	err = database.Db.First(&clube, id).Error
	if err != nil {
		return c.Status(404).SendString("clube não encontrado")
	}

	err = json.Unmarshal(c.Body(), &clube)
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}

	database.Db.Model(&clube).Updates(&clube)
	return c.SendStatus(fiber.StatusNoContent)
}

func DeleteClube(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}
	clube := &types.Clube{}
	err = database.Db.First(&clube, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(struct{ Message string }{Message: err.Error()})
	}
	database.Db.Delete(&clube)
	return c.SendStatus(204)
}
