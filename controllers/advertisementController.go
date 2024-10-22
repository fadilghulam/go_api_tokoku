package controllers

import (
	"fmt"
	db "go_api_tokoku/config"
	"go_api_tokoku/helpers"
	model "go_api_tokoku/models"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
)

func GetAdvertisement(c *fiber.Ctx) error {

	advertisements := []model.Advertisement{}
	err := db.DB.Find(&advertisements).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}
	return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
		Message: "Success",
		Success: true,
		Data:    advertisements,
	})
}

func executeWithoutResult(query string, wg *sync.WaitGroup) {
	defer wg.Done()

	db.DB.Exec(query)
}

func GenerateTransactionsUserId(c *fiber.Ctx) error {

	masterTable := []string{"kunjungan", "penjualan", "payment", "piutang",
		"pengembalian", "pembayaran_piutang", "stok_salesman_riwayat",
		"stok_salesman", "md.stok_merchandiser_riwayat", "md.stok_merchandiser",
		"kunjungan_log", "md.transaction"}

	masterQuery := `UPDATE :tableName SET user_id = data.user_id, user_id_subtitute = data.user_id_subtitute
					FROM (
					SELECT k.id, 
							COALESCE(s.user_id, mm.user_id, tl.user_id) as user_id,
							CASE WHEN k.teamleader_id IS NULL THEN NULL ELSE COALESCE(tl.user_id, k.teamleader_id) END as user_id_subtitute
					FROM :tableName k
					LEFT JOIN salesman s
						ON k.salesman_id = s.id
					LEFT JOIN md.merchandiser mm
						ON k.merchandiser_id = mm.id
					LEFT JOIN teamleader tl
						ON k.teamleader_id = tl.id
					WHERE k.user_id IS NULL
					) data
					WHERE :tableName.id = data.id`

	queries := []string{}

	for _, tableName := range masterTable {
		// Replace the placeholder with the actual table name
		query := strings.ReplaceAll(masterQuery, ":tableName", tableName)
		// Append the modified query to the queries slice
		queries = append(queries, query)
	}

	// var wg sync.WaitGroup

	// // Launch concurrent Goroutines
	// for _, query := range queries {
	// 	wg.Add(1)
	// 	go executeWithoutResult(query, &wg)
	// }

	// // Wait for all Goroutines to finish
	// wg.Wait()

	for _, query := range queries {
		fmt.Println(query) // Prints with proper formatting
		fmt.Println("------")
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data has been loaded successfully",
		Success: true,
		Data:    queries,
	})

}
