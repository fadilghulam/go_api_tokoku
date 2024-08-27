package controllers

import (
	"go_api_tokoku/helpers"

	"github.com/gofiber/fiber/v2"
)

func RefreshUser(c *fiber.Ctx) error {

	userId := c.FormValue("userId")

	// datas, err := helpers.ExecuteQuery(fmt.Sprintf(`SELECT NULL as employee,
	// 													u.id,
	// 													u.full_name as name,
	// 													u.username,
	// 													u.profile_photo,
	// 													p.email,
	// 													p.phone,
	// 													p.ktp,
	// 													ARRAY[]::varchar[] as permission,
	// 													CASE WHEN MAX(c.id) IS NULL THEN NULL ELSE JSONB_AGG(c.*) END as "userInfo"
	// 												FROM public.user u
	// 												LEFT JOIN customer c
	// 													ON u.id = c.user_id
	// 												LEFT JOIN hr.person p
	// 													ON u.id = p.user_id
	// 												WHERE u.id = %v
	// 												GROUP BY u.id, p.id`, userId))

	datas, err := helpers.RefreshUser(userId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	// fmt.Println("flag datas")
	// fmt.Println(datas)

	if len(datas) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(helpers.ResponseWithoutData{
			Message: "User not found",
			Success: false,
		})
	}

	returnMap := make(map[string]interface{})

	tokenString, expired, err := CreateJWT(userId)
	jwtMap := map[string]interface{}{
		"expired": expired,
		"token":   tokenString,
	}
	returnMap["auth"] = true
	returnMap["data"] = datas
	returnMap["jwt"] = jwtMap

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Success",
		Success: true,
		Data:    returnMap,
	})
}
