package controllers

import (
	db "go_api_tokoku/config"
	"go_api_tokoku/helpers"
	model "go_api_tokoku/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetMembership(c *fiber.Ctx) error {

	type returnData struct {
		Message      string      `json:"message"`
		Success      bool        `json:"success"`
		Data         interface{} `json:"datas"`
		Requirements interface{} `json:"requirements"`
	}

	datas, err := helpers.ExecuteQuery(strings.ReplaceAll(`WITH split_parts AS (
											SELECT id, 
												regexp_split_to_table(mb.formula_requirement, @\sOR\s|\sAND\s|,\s*@) AS part
											FROM tk.membership mb
												GROUP BY mb.id
										), extracted_ids AS (
											SELECT id, 
												regexp_replace(part, @[^0-9]@, @@, @g@) AS ids
											FROM split_parts
										), requirement_ids as (
												SELECT id, array_agg(ids) AS extracted_ids
												FROM extracted_ids
												GROUP BY id
										)

											SELECT mb.id,
															mb.name,
															mb.description,
															mb.formula_requirement,
															mb.formula_benefit,
										-- 					ri.extracted_ids as requirement_ids,
															JSONB_AGG(
																JSONB_BUILD_OBJECT(
																	'id', mbd.id,
																	'type', mbd.type,
																	'value', mbd.value,
																	'quota', mbd.quota,
																	'is_percentage_value', mbd.is_percentage_value
																) ORDER BY mbd.id
															) as benefit
											FROM tk.membership mb
											JOIN requirement_ids ri
												ON mb.id = ri.id
											LEFT JOIN tk.membership_benefit mbd
												ON mb.id = mbd.membership_id
											GROUP BY mb.id, ri.id, ri.extracted_ids`, "@", "'"))

	MembershipRequirements := []model.MembershipRequirement{}
	err2 := db.DB.Find(&MembershipRequirements).Error

	if err != nil || err2 != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(
		returnData{
			Message:      "Success",
			Success:      true,
			Data:         datas,
			Requirements: MembershipRequirements,
		},
	)
}
