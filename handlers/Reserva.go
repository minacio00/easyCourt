package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/minacio00/easyCourt/database"
	"github.com/minacio00/easyCourt/types"
)

func CreateReserva(c *fiber.Ctx) error {
	reserva := &types.Reserva{}
	err := json.Unmarshal(c.Body(), reserva)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	err = database.Db.Create(reserva).Error
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}

	return c.Status(201).JSON(&reserva)
}

func ReadReserva(c *fiber.Ctx) error {
	reserva := &types.Reserva{}
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}
	err = database.Db.First(&reserva, id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(struct{ Message string }{Message: err.Error()})
	}
	return c.Status(200).JSON(&reserva)
}

func GetReservas(c *fiber.Ctx) error {
	reservas := []types.Reserva{}

	err := database.Db.Find(&reservas).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(struct{ Message string }{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(&reservas)
}

func UpdateReserva(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id == 0 {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}

	body := types.Reserva{}
	err = c.BodyParser(&body)
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}
	reserva := types.Reserva{}
	err = database.Db.First(&reserva, id).Error
	if err != nil {
		return c.Status(404).SendString("reserva não encontrada")
	}

	err = json.Unmarshal(c.Body(), &reserva)
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}

	database.Db.Model(&reserva).Updates(&reserva)
	return c.SendStatus(fiber.StatusNoContent)
}

func DeleteReserva(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}
	reserva := &types.Reserva{}
	err = database.Db.First(&reserva, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(struct{ Message string }{Message: err.Error()})
	}
	database.Db.Unscoped().Delete(&reserva)
	return c.SendStatus(204)
}
