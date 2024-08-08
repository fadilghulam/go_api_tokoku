package controllers

import (
	"fmt"

	"go_api_tokoku/helpers"

	"github.com/gofiber/fiber/v2"
)

func GetCountTransactions(c *fiber.Ctx) error {

	customerId := c.Query("id")

	result, err := helpers.ExecuteQuery(fmt.Sprintf(`SELECT JSON_AGG(data.count_data) as count_data FROM (
														SELECT --ts.name as transaction_state, COUNT(tr.id) as count_data,
														JSONB_BUILD_OBJECT(ts.name, COUNT(tr.id)) as count_data
														FROM tk.transaction_state ts
														LEFT JOIN tk.transaction tr
															ON ts.id = tr.transaction_state_id
															AND tr.customer_id = %s
														GROUP BY ts.id
														ORDER BY ts.id
													) data`, customerId))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	if len(result) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.Response{
			Message: "Data not found",
			Success: true,
			Data:    nil,
		})
	}

	flattenedResult := make(map[string]int)
	for _, item := range result {
		for _, value := range item {
			for _, val := range value.([]interface{}) {
				for k, v := range val.(map[string]interface{}) {
					flattenedResult[k] = int(v.(float64))
				}
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Success",
		Success: true,
		Data:    flattenedResult,
	})

}

func GetTransactions(c *fiber.Ctx) error {

	customerId := c.Query("id")
	types := c.Query("type")

	results, err := helpers.ExecuteQuery(fmt.Sprintf(
		`WITH data_details as (
			SELECT 
				trd.transaction_id,
				JSONB_AGG(
					JSONB_BUILD_OBJECT(
						'id', trd.id,
						'produk_id', trd.produk_id,
						'point', trd.point,
						'code', p.code,
						'name', p.name,
						'image', p.foto,
						'qty', trd.qty,
						'unit', iu.name,
						'harga', trd.harga,
						'note', trd.note,
						'produk_satuan', JSONB_BUILD_OBJECT(
											'carton', ps.carton,
											'ball', ps.ball,
											'slof', ps.slof,
											'pack', ps.pack
										)
					) ORDER BY p.order
				) as details,
				JSONB_BUILD_OBJECT(
					'id', tr.reference_id,
					'name', s.name,
					'phone', s.phone,
					'type', tr.reference_name
				) as courier,
				JSONB_BUILD_OBJECT(
					'id', c.id,
					'name', c.name,
					'outlet_name', c.outlet_name,
					'phone', c.phone,
					'address', c.alamat,
					'kelurahan', c.kelurahan,
					'kabupaten', c.kabupaten,
					'kecamatan', c.kecamatan,
					'provinsi', c.provinsi
				) as customer,
				JSONB_BUILD_OBJECT(
					'no_transaksi', tr.id,
					'jumlah_pesanan', SUM(trd.qty),
					'jumlah_point', SUM(trd.point),
					'voucher_tokoku', 0,
					'total_pembayaran', tr.total_transaction
				) as invoice,
				rev.rating,
				JSONB_BUILD_OBJECT(
					'rating', rev.rating,
					'description', rev.description,
					'photo', rev.photo
				) as review,
				comp.id as complaint_id,
				JSONB_BUILD_OBJECT(
					 'description', comp.description,
					 'photo', comp.image
				) as complaint
			FROM tk.transaction_detail trd
			JOIN tk.transaction tr
				ON trd.transaction_id = tr.id
			JOIN tk.transaction_state ts
				ON tr.transaction_state_id = ts.id
			JOIN customer c
				ON tr.customer_id = c.id
			JOIN produk p
				ON trd.produk_id = p.id
			JOIN produk_satuan ps
				ON p.satuan_id = ps.id
			LEFT JOIN tk.unit_mapping um
				ON p.id = um.produk_id
			LEFT JOIN tk.item_unit_tk iu
				ON um.item_unit_id = iu.id
			LEFT JOIN salesman s
				ON (tr.reference_id = s.id
				AND tr.reference_name = 'SALESMAN')
			LEFT JOIN tk.review rev
				ON tr.id = rev.order_id
				AND rev.order_item_id IS NULL
			LEFT JOIN tk.complaints comp
				ON tr.id = comp.transaction_id
				AND comp.transaction_item_id IS NULL
			WHERE tr.customer_id = %s AND ts.name = UPPER('%s')
			GROUP BY trd.transaction_id, tr.id, c.id, s.id, rev.id, comp.id
		)
		
		SELECT LOWER(name) as name,
				JSONB_AGG(
					JSONB_BUILD_OBJECT(
						'id', tr.id,
						'type', ts.name,
						'transaction_date', tr.transaction_date,
						'total_order', tr.total_transaction,
						'products', trd.details,
						'courier', trd.courier,
						'customer', trd.customer,
						'invoice', trd.invoice,
						'note', tr.note,
						'review', CASE WHEN trd.rating IS NULL THEN NULL ELSE trd.review END,
						'complaints', CASE WHEN trd.complaint_id IS NULL THEN NULL ELSE trd.complaint END
					) ORDER BY tr.transaction_date DESC
				) FILTER (WHERE tr.id IS NOT NULL) as datas
		FROM tk.transaction_state ts
		LEFT JOIN tk.transaction tr
			ON ts.id = tr.transaction_state_id
			AND tr.customer_id = %s
		LEFT JOIN data_details trd
			ON tr.id = trd.transaction_id
		WHERE ts.name = UPPER('%s')
		GROUP BY ts.name`, customerId, types, customerId, types))

	if results == nil {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
			Message: "Data not found",
			Success: true,
			Data:    nil,
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseDataMultiple{
			Message: err.Error(),
			Success: false,
			Data:    nil,
		})
	}

	// newResults := make(map[string]interface{})
	// for _, val := range results {
	// 	// fmt.Println(i, val)
	// 	newResults[val["name"].(string)] = val["datas"]
	// }

	var newResults interface{}
	for _, val := range results {
		// fmt.Println(i, val)
		newResults = val["datas"]
	}

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
		Message: "Data has been loaded successfully",
		Success: true,
		Data:    newResults,
	})
}

func GetPointsCustomer(c *fiber.Ctx) error {

	customerId := c.Query("customerId")

	results, err := helpers.ExecuteQuery(fmt.Sprintf(
		`SELECT cph.customer_id,
				SUM(CASE WHEN cph.exchange_id IS NOT NULL THEN cph.point ELSE 0 END) as total_point_digunakan,
				SUM(CASE WHEN cph.transaction_id IS NOT NULL THEN cph.point ELSE 0 END) as total_point_terkumpul
		FROM tk.customer_point_history cph
		LEFT JOIN tk.transaction t
			ON cph.transaction_id = t.id
		WHERE TRUE AND cph.customer_id = %s
		GROUP BY cph.customer_id`, customerId))

	if results == nil {
		return c.Status(fiber.StatusOK).JSON(helpers.Response{
			Message: "Data not found",
			Success: true,
			Data:    nil,
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.Response{
			Message: err.Error(),
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

func GetPointsHistory(c *fiber.Ctx) error {

	customerId := c.Query("customerId")
	dateStart := c.Query("dateStart")
	dateEnd := c.Query("dateEnd")
	types := c.Query("type")

	whereQuery := ""
	if dateStart != "" && dateEnd != "" {
		whereQuery = " AND DATE(cph.datetime) BETWEEN DATE('" + dateStart + "') AND DATE('" + dateEnd + "')"
	}

	whereType := ""
	if types != "" {
		whereType = " AND UPPER(cph.type) = UPPER('" + types + "')"
	}

	results, err := helpers.ExecuteQuery(fmt.Sprintf(
		`SELECT cph.id,
				cph.customer_id,
				cph.transaction_id,
				cph.exchange_id,
				cph.datetime,
				COALESCE(cph.point,0) as point,
				cph.expired_date,
				cph.created_at,
				cph.updated_at, 
				t.transaction_date,
				cph.type,
				CASE 
					WHEN CURRENT_DATE > cph.expired_date THEN 'Expired at '||cph.expired_date
					WHEN cph.transaction_id IS NOT NULL THEN 'Transaction at '||t.transaction_date
				END as note
		FROM tk.customer_point_history cph
		LEFT JOIN tk.transaction t
		ON cph.transaction_id = t.id
		WHERE TRUE AND cph.customer_id = %s %s %s
		ORDER BY cph.datetime DESC`, customerId, whereQuery, whereType))

	if results == nil {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
			Message: "Data not found",
			Success: true,
			Data:    nil,
		})
	}

	if err != nil {
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
