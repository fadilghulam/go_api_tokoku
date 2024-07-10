package controllers

import (
	"fmt"

	"go_api_tokoku/helpers"

	"github.com/gofiber/fiber/v2"
)

func GetProdukTerkini(c *fiber.Ctx) error {

	customerId := c.Query("id")

	var query string

	if customerId == "" {
		query = `SELECT p.id,
					p.code,
					p.name,
					p.foto as image,
					0 as ranking,
					null as promo,
					0 as harga,
					0 as stock,
					-- CASE WHEN MOD(ROW_NUMBER() OVER(ORDER BY p.order), 2) = 0 THEN 100 ELSE 0 END as point,
					-- rh.diskon as discount
					-- CASE WHEN MOD(ROW_NUMBER() OVER(ORDER BY p.order), 2) > 0 THEN 10 ELSE 0 END as discount
					0 as point,
					0 as discount,
					null as produk_satuan
				FROM produk p
				WHERE p.is_aktif = 1
				ORDER BY p.order`
	} else {
		query = fmt.Sprintf(
			`SELECT p.id,
				p.code,
				p.name,
				p.foto as image,
				--CASE WHEN ROW_NUMBER() OVER(ORDER BY p.order) <= 3 THEN ROW_NUMBER() OVER(ORDER BY p.order) ELSE NULL END as ranking,
				0 as ranking,
				--JSONB_BUILD_OBJECT(
				--	'id', 1,
				--	'name', 'Promo 123'
				--) as promo,
				NULL as promo,
				rh.harga,
				0 as stock,
				-- CASE WHEN MOD(ROW_NUMBER() OVER(ORDER BY p.order), 2) = 0 THEN 100 ELSE 0 END as point,
				-- rh.diskon as discount
				-- CASE WHEN MOD(ROW_NUMBER() OVER(ORDER BY p.order), 2) > 0 THEN 10 ELSE 0 END as discount
				COALESCE(pt.value, 0) as point,
				COALESCE(dis.nominal,0) as discount,
				JSONB_BUILD_OBJECT(
					'carton', ps.carton,
					'ball', ps.ball,
					'slof', ps.slof,
					'pack', ps.pack
				) as produk_satuan
			FROM customer c
			JOIN salesman s
				ON c.salesman_id = s.id
			JOIN produk_branch pb
				ON COALESCE(c.branch_id,s.branch_id) = pb.branch_id
			JOIN produk p
				ON pb.produk_id = p.id
			JOIN produk_satuan ps
				ON p.satuan_id = ps.id
			LEFT JOIN ref_harga_master rhm
				ON COALESCE(c.branch_id,s.branch_id) = rhm.branch_id
				AND CURRENT_DATE BETWEEN rhm.date_start AND COALESCE(rhm.date_end, CURRENT_DATE)
			LEFT JOIN ref_harga rh
				ON rhm.id = rh.ref_harga_master_id
				AND c.tipe = rh.customer_tipe
				AND s.tipe_salesman = rh.salesman_tipe
				AND p.id = rh.produk_id
			LEFT JOIN tk.discount dis
				ON p.id = dis.produk_id
				AND CURRENT_DATE BETWEEN dis.date_start AND COALESCE(dis.date_end, CURRENT_DATE)
				AND COALESCE(c.branch_id,s.branch_id) = dis.branch_id
				AND c.tipe = dis.customer_type_id
			LEFT JOIN tk.points pt
				ON p.id = pt.produk_id
				AND CURRENT_DATE BETWEEN pt.date_start AND COALESCE(pt.date_end, CURRENT_DATE)
				AND COALESCE(c.branch_id,s.branch_id) = pt.branch_id
			WHERE c.id = %s
			ORDER BY p.order`, customerId)
	}
	results, err := helpers.ExecuteQuery(query)

	if err != nil {
		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"message": err.Error(),
		// 	"success": false,
		// 	"data":    nil,
		// })
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseDataMultiple{
			Message: err.Error(),
			Success: false,
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
		Message: "Data has been loaded successfully",
		Success: true,
		Data:    results,
	})
}

func GetProdukDetail(c *fiber.Ctx) error {
	produkId := c.Query("produkId")
	customerId := c.Query("customerId")

	var query string

	if customerId == "" {
		query = fmt.Sprintf(`SELECT p.id,
									p.name,
									p.code,
									p.foto as image,
									'-' as description,
									pk.name as type,
									0 as harga,
									0 as discount,
									0 as point,
									0 as stock,
									JSONB_BUILD_OBJECT(
										'carton', ps.carton,
										'ball', ps.ball,
										'slof', ps.slof,
										'pack', ps.pack
									) as produk_satuan
								FROM produk p
								JOIN produk_kategori pk
									ON p.kategori_id = pk.id
								JOIN produk_satuan ps
									ON p.satuan_id = ps.id
								WHERE p.id = %s`, produkId)
	} else {
		query = fmt.Sprintf(
			`SELECT p.id,
							p.name,
							p.code,
							p.foto as image,
							'-' as description,
							pk.name as type,
							rh.harga,
							COALESCE(dis.nominal,0) as discount,
							COALESCE(pt.value,0) as point,
							0 as stock,
							JSONB_BUILD_OBJECT(
								'carton', ps.carton,
								'ball', ps.ball,
								'slof', ps.slof,
								'pack', ps.pack
							) as produk_satuan
					FROM produk p
					JOIN produk_kategori pk
						ON p.kategori_id = pk.id
					JOIN produk_satuan ps
						ON p.satuan_id = ps.id
					JOIN customer c
						ON c.id = %s
					JOIN salesman s
						ON c.salesman_id = s.id
					LEFT JOIN ref_harga_master rhm
						ON COALESCE(c.branch_id,s.branch_id) = rhm.branch_id
						AND CURRENT_DATE BETWEEN rhm.date_start AND COALESCE(rhm.date_end, CURRENT_DATE)
					LEFT JOIN ref_harga rh
						ON rhm.id = rh.ref_harga_master_id
						AND c.tipe = rh.customer_tipe
						AND s.tipe_salesman = rh.salesman_tipe
						AND p.id = rh.produk_id
					LEFT JOIN tk.discount dis
						ON p.id = dis.produk_id
						AND CURRENT_DATE BETWEEN dis.date_start AND COALESCE(dis.date_end, CURRENT_DATE)
						AND COALESCE(c.branch_id,s.branch_id) = dis.branch_id
						AND c.tipe = dis.customer_type_id
					LEFT JOIN tk.points pt
						ON p.id = pt.produk_id
						AND CURRENT_DATE BETWEEN pt.date_start AND COALESCE(pt.date_end, CURRENT_DATE)
						AND COALESCE(c.branch_id,s.branch_id) = pt.branch_id
					WHERE p.id = %s`, customerId, produkId)
	}

	results, err := helpers.ExecuteQuery(query)

	if err != nil {
		// return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 	"message": err.Error(),
		// 	"success": false,
		// 	"data":    nil,
		// })
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.Response{
			Message: err.Error(),
			Success: false,
			Data:    nil,
		})
	}

	if len(results) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(helpers.Response{
			Message: "Data not found",
			Success: false,
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data has been loaded successfully",
		Success: true,
		Data:    results[0],
	})
}

func TestRoute(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}
