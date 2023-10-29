package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/minacio00/easyCourt/database"
	"github.com/minacio00/easyCourt/types"
)

func CreateCliente(c *fiber.Ctx) error {
	cliente := &types.Cliente{}
	err := json.Unmarshal(c.Body(), cliente)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	err = database.Db.Create(cliente).Error
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}

	return c.Status(201).JSON(&cliente)
}

func ReadCliente(c *fiber.Ctx) error {
	cliente := &types.Cliente{}
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}
	err = database.Db.First(&cliente, id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(struct{ Message string }{Message: err.Error()})
	}
	return c.Status(200).JSON(&cliente)
}

func GetClientes(c *fiber.Ctx) error {
	clientes := []types.Cliente{}

	err := database.Db.Find(&clientes).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(struct{ Message string }{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(&clientes)
}

func UpdateCliente(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id == 0 {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}

	body := types.Cliente{}
	err = c.BodyParser(&body)
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}
	cliente := types.Cliente{}
	err = database.Db.First(&cliente, id).Error
	if err != nil {
		return c.Status(404).SendString("cliente não encontrado")
	}

	err = json.Unmarshal(c.Body(), &cliente)
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}

	database.Db.Model(&cliente).Updates(&cliente)
	return c.SendStatus(fiber.StatusNoContent)
}

func DeleteCliente(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}
	cliente := &types.Cliente{}
	err = database.Db.First(&cliente, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(struct{ Message string }{Message: err.Error()})
	}
	database.Db.Delete(&cliente)
	return c.SendStatus(204)
}
