package handlers

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/minacio00/easyCourt/database"
	"github.com/minacio00/easyCourt/types"
)

func CreateTenant(c *fiber.Ctx) error {
	tenant := &types.Tenant{}
	err := json.Unmarshal(c.Body(), tenant)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	err = tenant.Validate()
	if err != nil {
		return c.Status(400).JSON(struct{ Message string }{Message: err.Error()}) // todo: debugar isso daqui
	}
	err = database.Db.Create(tenant).Error
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}

	return c.Status(201).JSON(&tenant)
}

// todo: fazer dto
func ReadTenant(c *fiber.Ctx) error {
	tenant := &types.Tenant{}
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}
	err = database.Db.Preload("clubes").First(&tenant, id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(struct{ Message string }{Message: err.Error()})
	}
	return c.Status(200).JSON(&tenant)
}

func GetTenants(c *fiber.Ctx) error {
	tenants := []types.Tenant{}

	err := database.Db.Find(&tenants).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(struct{ Message string }{Message: err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(&tenants)
}
func UpdateTenant(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id == 0 {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}

	body := types.Tenant{}
	err = c.BodyParser(&body)
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}
	tenant := types.Tenant{}
	err = database.Db.First(&tenant, id).Error
	if err != nil {
		return c.Status(404).SendString("tenant não encontrado")
	}

	err = json.Unmarshal(c.Body(), &tenant)
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}

	database.Db.Model(&tenant).Updates(&tenant)
	return c.SendStatus(fiber.StatusNoContent)
}

func DeleteTenant(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(struct{ Message string }{Message: err.Error()})
	}
	tenant := &types.Tenant{}
	err = database.Db.First(&tenant, "id = ?", id).Error
	if err != nil {
		return c.Status(404).JSON(struct{ Message string }{Message: err.Error()})
	}
	return c.SendStatus(204)

}
